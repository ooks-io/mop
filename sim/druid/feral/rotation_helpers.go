package feral

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/druid"
)

func (cat *FeralDruid) tfExpectedBefore(sim *core.Simulation, futureTime time.Duration) bool {
	if !cat.TigersFury.IsReady(sim) {
		return cat.TigersFury.ReadyAt() < futureTime
	}
	if cat.BerserkCatAura.IsActive() {
		return cat.BerserkCatAura.ExpiresAt() < futureTime
	}
	return true
}

func (rotation *FeralDruidRotation) WaitUntil(sim *core.Simulation, nextEvaluation time.Duration) {
	rotation.nextActionAt = nextEvaluation
	rotation.agent.WaitUntil(sim, nextEvaluation)
}

func (cat *FeralDruid) calcTfEnergyThresh() float64 {
	delayTime := cat.ReactionTime + core.TernaryDuration(cat.ClearcastingAura.IsActive(), time.Second, 0)
	return 40.0 - delayTime.Seconds() * cat.EnergyRegenPerSecond()
}

func (rotation *FeralDruidRotation) TryTigersFury(sim *core.Simulation) {
	cat := rotation.agent

	if !cat.TigersFury.IsReady(sim) {
		return
	}

	tfEnergyThresh := cat.calcTfEnergyThresh()
	tfNow := (cat.CurrentEnergy() < tfEnergyThresh) && !cat.BerserkCatAura.IsActive()

	if tfNow {
		cat.TigersFury.Cast(sim, nil)
		rotation.WaitUntil(sim, sim.CurrentTime + cat.ReactionTime)
	}
}

func (rotation *FeralDruidRotation) TryBerserk(sim *core.Simulation) {
	// Berserk algorithm: time Berserk for just after a Tiger's Fury
	// *unless* we'll lose Berserk uptime by waiting for Tiger's Fury to
	// come off cooldown.
	cat := rotation.agent
	simTimeRemain := sim.GetRemainingDuration()
	tfCdRemain := cat.TigersFury.TimeToReady(sim)
	waitForTf := (tfCdRemain <= cat.BerserkCatAura.Duration) && (tfCdRemain + cat.ReactionTime < simTimeRemain - cat.BerserkCatAura.Duration)
	berserkNow := rotation.UseBerserk && cat.Berserk.IsReady(sim) && !waitForTf && !cat.ClearcastingAura.IsActive()

	if berserkNow {
		cat.Berserk.Cast(sim, nil)
		cat.UpdateMajorCooldowns()
		rotation.WaitUntil(sim, sim.CurrentTime + cat.ReactionTime)
	}
}

func (rotation *FeralDruidRotation) ShiftBearCat(sim *core.Simulation) {
	rotation.readyToShift = false
	rotation.lastShiftAt = sim.CurrentTime
	cat := rotation.agent

	if cat.InForm(druid.Cat) {
		cat.BearForm.Cast(sim, nil)
	} else {
		cat.CatForm.Cast(sim, nil)

		// Reset swing timer with Albino Snake when advantageous
		if cat.AutoAttacks.NextAttackAt() - sim.CurrentTime > cat.AutoAttacks.MainhandSwingSpeed() {
			cat.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime)
		}
	}
}

func (cat *FeralDruid) calcBleedRefreshTime(sim *core.Simulation, bleedSpell *druid.DruidSpell, bleedDot *core.Dot, isExecutePhase bool, isRip bool) time.Duration {
	if !bleedDot.IsActive() {
		return sim.CurrentTime - cat.ReactionTime
	}

	// If we're not gaining a stronger snapshot, then use the standard 1
	// tick refresh window.
	bleedEnd := bleedDot.ExpiresAt()
	standardRefreshTime := bleedEnd - bleedDot.BaseTickLength

	if !cat.tempSnapshotAura.IsActive() {
		return standardRefreshTime
	}

	// For Rip specifically, also bypass clipping calculations if CP count
	// is too low for the calculation to be relevant.
	if isRip && (cat.ComboPoints() < 5) {
		return standardRefreshTime
	}

	// Likewise, if the existing buff will still be up at the start of the normal
	// window, then don't clip unnecessarily. For long buffs that cover a full bleed
	// duration, project "buffEnd" forward in time such that we block clips if we are
	// already maxing out the number of full durations we can snapshot.
	buffRemains := cat.tempSnapshotAura.RemainingDuration(sim)
	maxTickCount := core.TernaryInt32(isRip, druid.RipMaxNumTicks, bleedDot.BaseTickCount)
	maxBleedDur := bleedDot.BaseTickLength * time.Duration(maxTickCount)
	numCastsCovered := buffRemains / maxBleedDur
	buffEnd := cat.tempSnapshotAura.ExpiresAt() - numCastsCovered*maxBleedDur

	if buffEnd > standardRefreshTime+cat.ReactionTime {
		return standardRefreshTime
	}

	// Potential clips for a buff snapshot should be done as late as possible
	latestPossibleSnapshot := buffEnd - cat.ReactionTime*time.Duration(2)
	numClippedTicks := (bleedEnd - latestPossibleSnapshot) / bleedDot.BaseTickLength
	targetClipTime := standardRefreshTime - numClippedTicks*bleedDot.BaseTickLength

	// Since the clip can cost us 30-35 Energy, we need to determine whether the damage gain is worth the
	// spend. First calculate the maximum number of buffed bleed ticks we can get out before the fight
	// ends.
	buffedTickCount := min(maxTickCount, int32((sim.Duration-targetClipTime)/bleedDot.BaseTickLength))

	// Perform a DPE comparison vs. Shred
	expectedDamageGain := (bleedSpell.NewSnapshotPower - bleedSpell.CurrentSnapshotPower) * float64(buffedTickCount)

	// For Rake specifically, we get 1 free "tick" immediately upon cast.
	if !isRip {
		expectedDamageGain += bleedSpell.NewSnapshotPower
	}

	energyEquivalent := expectedDamageGain / cat.Shred.ExpectedInitialDamage(sim, cat.CurrentTarget) * cat.Shred.DefaultCast.Cost

	// Finally, discount the effective Energy cost of the clip based on the number of clipped ticks.
	discountedRefreshCost := core.TernaryFloat64(isRip, float64(numClippedTicks) / float64(maxTickCount), 1.0) * bleedSpell.DefaultCast.Cost

	if sim.Log != nil {
		cat.Log(sim, "%s buff snapshot is worth %.1f Energy, discounted refresh cost is %.1f Energy.", bleedSpell.ShortName, energyEquivalent, discountedRefreshCost)
	}

	return core.TernaryDuration(energyEquivalent > discountedRefreshCost, targetClipTime, standardRefreshTime)
}

// Determine whether Tiger's Fury will be usable soon enough for the snapshot to
// outweigh the lost Rip/Rake ticks from delaying a refresh.
func (cat *FeralDruid) shouldDelayBleedRefreshForTf(sim *core.Simulation, bleedDot *core.Dot, isRip bool) bool {
	if cat.TigersFuryAura.IsActive() || cat.BerserkCatAura.IsActive() {
		return false
	}

	finalTickLeeway := core.TernaryDuration(bleedDot.IsActive(), bleedDot.TimeUntilNextTick(sim), 0)
	maxTickCount := core.TernaryInt32(isRip, druid.RipMaxNumTicks, bleedDot.BaseTickCount)
	buffedTickCount := min(maxTickCount, int32((sim.GetRemainingDuration() - finalTickLeeway) / bleedDot.BaseTickLength))
	delayBreakpoint := finalTickLeeway + core.DurationFromSeconds(0.15 * float64(buffedTickCount) * bleedDot.BaseTickLength.Seconds())

	if !cat.tfExpectedBefore(sim, sim.CurrentTime + delayBreakpoint) {
		return false
	}

	if isRip && cat.tempSnapshotAura.IsActive() && (cat.tempSnapshotAura.RemainingDuration(sim) <= delayBreakpoint) {
		return false
	}

	delaySeconds := delayBreakpoint.Seconds()
	energyToDump := cat.CurrentEnergy() + delaySeconds * cat.EnergyRegenPerSecond() - cat.calcTfEnergyThresh()
	secondsToDump := math.Ceil(energyToDump / cat.Shred.DefaultCast.Cost)
	return secondsToDump < delaySeconds
}

func (cat *FeralDruid) calcRoarRefreshTime(sim *core.Simulation, ripLeeway time.Duration, minRoarOffset time.Duration) time.Duration {
	roarBuff := cat.SavageRoarBuff
	ripDot := cat.Rip.CurDot()

	if !roarBuff.IsActive() {
		return sim.CurrentTime - cat.ReactionTime
	}

	// If we're not proactively offsetting the Roar, then use the standard 1
	// tick refresh window.
	roarEnd := roarBuff.ExpiresAt()
	standardRefreshTime := roarEnd - roarBuff.BaseTickLength

	if !ripDot.IsActive() {
		return standardRefreshTime
	}

	// Project Rip end time assuming full Bloodletting extensions
	remainingExtensions := druid.RipMaxNumTicks - ripDot.BaseTickCount
	ripEnd := ripDot.ExpiresAt() + time.Duration(remainingExtensions) * ripDot.BaseTickLength
	fightEnd := sim.Duration

	if roarEnd > (ripEnd + ripLeeway) {
		return standardRefreshTime
	}

	if roarEnd >= fightEnd {
		return standardRefreshTime
	}

	// Potential clips for offsetting timers should be done just after a
	// Roar "tick" in order to exploit the Pandemic behavior in MoP.
	targetClipTime := roarBuff.NextTickAt()

	// Calculate when Roar would end if refreshed at the optimal clip time.
	newRoarDur := cat.SavageRoarDurationTable[cat.ComboPoints()]
	newRoarEnd := targetClipTime + newRoarDur + roarBuff.BaseTickLength

	// If a fresh Roar cast would cover us to the end of the fight, then
	// clip at the next tick for maximum CP efficiency.
	if newRoarEnd >= fightEnd {
		return targetClipTime
	}

	// Outside of Execute, use offset rule to determine whether to clip.
	if !sim.IsExecutePhase25() {
		return core.TernaryDuration(newRoarEnd >= ripEnd + minRoarOffset, targetClipTime, standardRefreshTime)
	}

	// Under Execute conditions, ignore the offset rule and instead optimize
	// for as few Roar casts as possible.
	if cat.ComboPoints() < 5 {
		return standardRefreshTime
	}

	minRoarsPossible := (fightEnd - roarEnd) / newRoarDur
	projectedRoarCasts := (fightEnd - newRoarEnd) / newRoarDur + 1
	return core.TernaryDuration(projectedRoarCasts == minRoarsPossible, targetClipTime, standardRefreshTime)
}

func (cat *FeralDruid) canBearWeave(sim *core.Simulation, furorCap float64, regenRate float64, currentEnergy float64, excessEnergy float64, upcomingTimers *PoolingActions) bool {
	if cat.ClearcastingAura.IsActive() || cat.BerserkCatAura.IsActive() {
		return false
	}

	// If we can Shred now and then weave on the next GCD, prefer that.
	if excessEnergy > cat.Shred.DefaultCast.Cost {
		return false
	}

	// Calculate effective Energy cap for out-of-form pooling.
	targetWeaveDuration := core.GCDDefault*3 + cat.ReactionTime*2
	maxStartingEnergy := furorCap - targetWeaveDuration.Seconds() * regenRate

	if currentEnergy > maxStartingEnergy {
		return false
	}

	// Prioritize all timers over weaving.
	earliestWeaveEnd := sim.CurrentTime + core.GCDDefault*3 + cat.ReactionTime*2
	isPooling, nextRefresh := upcomingTimers.nextRefreshTime()

	if isPooling && (nextRefresh < earliestWeaveEnd) {
		return false
	}

	// Mana check
	if cat.CurrentMana() < cat.CatForm.DefaultCast.Cost * 2 {
		cat.Metrics.MarkOOM(sim)
		return false
	}

	// Also add a condition to make sure we can spend down our Energy
	// post-weave before the encounter ends or TF is ready.
	energyToDump := currentEnergy + (earliestWeaveEnd - sim.CurrentTime).Seconds() * regenRate
	timeToDump := earliestWeaveEnd + core.DurationFromSeconds(math.Floor(energyToDump / cat.Shred.DefaultCast.Cost))
	return (timeToDump < sim.Duration) && !cat.tfExpectedBefore(sim, timeToDump)
}

func (rotation *FeralDruidRotation) shouldTerminateBearWeave(sim *core.Simulation, isClearcast bool, currentEnergy float64, furorCap float64, regenRate float64, upcomingTimers *PoolingActions) bool {
	// Shift back early if a bear auto resulted in an Omen proc.
	if isClearcast && (sim.CurrentTime - rotation.lastShiftAt > core.GCDDefault) {
		return true
	}

	// Check Energy pooling leeway.
	cat := rotation.agent
	smallestWeaveExtension := core.GCDDefault + cat.ReactionTime
	finalEnergy := currentEnergy + smallestWeaveExtension.Seconds() * regenRate

	if finalEnergy > furorCap {
		return true
	}

	// Check timer leeway.
	earliestWeaveEnd := sim.CurrentTime + smallestWeaveExtension + core.GCDDefault
	isPooling, nextRefresh := upcomingTimers.nextRefreshTime()

	if isPooling && (nextRefresh < earliestWeaveEnd) {
		return true
	}

	// Also add a condition to prevent extending a weave if we don't have
	// enough time to spend the pooled Energy thus far.
	energyToDump := finalEnergy + 1.5 * regenRate // need to include Cat Form GCD here
	timeToDump := earliestWeaveEnd + core.DurationFromSeconds(math.Floor(energyToDump / cat.Shred.DefaultCast.Cost))
	return (timeToDump >= sim.Duration) || cat.tfExpectedBefore(sim, timeToDump)
}

func (rotation *FeralDruidRotation) ProcessNextPlannedAction(sim *core.Simulation, nextActionAt time.Duration) {
	// Also schedule an action right at Energy cap to make sure we never
	// accidentally over-cap while waiting on other timers.
	cat := rotation.agent
	timeToCap := core.DurationFromSeconds((cat.MaximumEnergy() - cat.CurrentEnergy()) / cat.EnergyRegenPerSecond())
	nextActionAt = min(nextActionAt, sim.CurrentTime + timeToCap)

	// Offset the ideal evaluation time by player latency.
	nextActionAt += cat.ReactionTime

	if nextActionAt <= sim.CurrentTime {
		panic("nextActionAt in the past!")
	} else {
		rotation.WaitUntil(sim, nextActionAt)
	}
}
