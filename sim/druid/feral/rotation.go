package feral

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (cat *FeralDruid) NewAPLAction(_ *core.APLRotation, config *proto.APLAction) core.APLActionImpl {
	switch config.Action.(type) {
	case *proto.APLAction_CatOptimalRotationAction:
		return cat.newActionCatOptimalRotationAction(config.GetCatOptimalRotationAction())
	default:
		return nil
	}
}

func (cat *FeralDruid) newActionCatOptimalRotationAction(config *proto.APLActionCatOptimalRotationAction) core.APLActionImpl {
	rotation := &FeralDruidRotation{
		APLActionCatOptimalRotationAction: config,
		agent:                             cat,
	}

	// Process rotation parameters
	rotation.ForceMangleFiller = cat.PseudoStats.InFrontOfTarget || cat.CannotShredTarget
	rotation.UseBerserk = (rotation.RotationType == proto.FeralDruid_Rotation_SingleTarget) || rotation.AllowAoeBerserk

	if rotation.ManualParams {
		rotation.BiteTime = core.DurationFromSeconds(config.BiteTime)
		rotation.BerserkBiteTime = core.DurationFromSeconds(config.BerserkBiteTime)
		rotation.MinRoarOffset = core.DurationFromSeconds(config.MinRoarOffset)
		rotation.RipLeeway = core.DurationFromSeconds(config.RipLeeway)
	} else {
		rotation.UseBite = true
		rotation.BiteTime = time.Second * 13
		rotation.BerserkBiteTime = time.Second * 8
		rotation.MinRoarOffset = time.Second * 37
		rotation.RipLeeway = time.Second * 8
	}

	// Pre-allocate PoolingActions
	rotation.pendingPool = &PoolingActions{}
	rotation.pendingPool.create(3)
	rotation.pendingPoolWeaves = &PoolingActions{}
	rotation.pendingPoolWeaves.create(2)

	return rotation
}

type FeralDruidRotation struct {
	*proto.APLActionCatOptimalRotationAction

	// Overwritten parameters
	BiteTime          time.Duration
	BerserkBiteTime   time.Duration
	MinRoarOffset     time.Duration
	RipLeeway         time.Duration
	ForceMangleFiller bool
	UseBerserk        bool

	// Bookkeeping fields
	agent             *FeralDruid
	lastActionAt      time.Duration
	nextActionAt      time.Duration
	pendingPool       *PoolingActions
	pendingPoolWeaves *PoolingActions
	readyToShift      bool
	lastShiftAt       time.Duration
}

func (rotation *FeralDruidRotation) Finalize(_ *core.APLRotation)                     {}
func (rotation *FeralDruidRotation) GetAPLValues() []core.APLValue                    { return nil }
func (rotation *FeralDruidRotation) GetInnerActions() []*core.APLAction               { return nil }
func (rotation *FeralDruidRotation) GetNextAction(_ *core.Simulation) *core.APLAction { return nil }
func (rotation *FeralDruidRotation) PostFinalize(_ *core.APLRotation)                 {}

func (rotation *FeralDruidRotation) IsReady(sim *core.Simulation) bool {
	return sim.CurrentTime > rotation.lastActionAt
}

func (rotation *FeralDruidRotation) Reset(_ *core.Simulation) {
	rotation.lastActionAt = -core.NeverExpires
	rotation.nextActionAt = -core.NeverExpires
	rotation.readyToShift = false
	rotation.lastShiftAt = -core.NeverExpires
}

func (rotation *FeralDruidRotation) Execute(sim *core.Simulation) {
	rotation.lastActionAt = sim.CurrentTime
	cat := rotation.agent

	// If a melee swing resulted in an Omen proc, then schedule the next
	// player decision based on latency.
	ccRefreshTime := cat.ClearcastingAura.ExpiresAt() - cat.ClearcastingAura.Duration

	if ccRefreshTime >= sim.CurrentTime - cat.ReactionTime {
		rotation.WaitUntil(sim, max(cat.NextGCDAt(), ccRefreshTime + cat.ReactionTime))
	}

	// Keep up Sunder debuff if not provided externally. Do this here since
	// FF can be cast while moving.
	for _, aoeTarget := range sim.Encounter.ActiveTargetUnits {
		if cat.ShouldFaerieFire(sim, aoeTarget) {
			cat.FaerieFire.CastOrQueue(sim, aoeTarget)
		}
	}

	// Off-GCD Maul check
	if cat.BearFormAura.IsActive() && !cat.ClearcastingAura.IsActive() && cat.Maul.CanCast(sim, cat.CurrentTarget) {
		cat.Maul.Cast(sim, cat.CurrentTarget)
	}

	// Handle movement before any rotation logic
	if cat.Moving || (cat.Hardcast.Expires > sim.CurrentTime) {
		return
	}

	if cat.DistanceFromTarget > core.MaxMeleeRange {
		// TODO: Wild Charge or Displacer Beast usage here
		if sim.Log != nil {
			cat.Log(sim, "Out of melee range (%.6fy) and cannot charge or teleport, initiating manual run-in...", cat.DistanceFromTarget)
		}

		cat.MoveTo(core.MaxMeleeRange - 1, sim) // movement aura is discretized in 1 yard intervals, so need to overshoot to guarantee melee range
		return
	}

	if !cat.GCD.IsReady(sim) {
		cat.WaitUntil(sim, cat.NextGCDAt())
		return
	}

	rotation.TryTigersFury(sim)
	rotation.TryBerserk(sim)

	if sim.CurrentTime < rotation.nextActionAt {
		cat.WaitUntil(sim, rotation.nextActionAt)
	} else if rotation.readyToShift {
		rotation.ShiftBearCat(sim)
	} else if rotation.RotationType == proto.FeralDruid_Rotation_SingleTarget {
		rotation.PickSingleTargetGCDAction(sim)
	} else {
		panic("AoE rotation not yet supported!")
	}
}

func (rotation *FeralDruidRotation) PickSingleTargetGCDAction(sim *core.Simulation) {
	// Store state variables for re-use
	cat := rotation.agent
	curEnergy := cat.CurrentEnergy()
	curCp := cat.ComboPoints()
	regenRate := cat.EnergyRegenPerSecond()
	isExecutePhase := sim.IsExecutePhase25()
	isClearcast := cat.ClearcastingAura.IsActive()
	anyBleedActive := cat.AssumeBleedActive || (cat.BleedsActive[cat.CurrentTarget] > 0)
	fightDur := sim.GetRemainingDuration()

	// Rip logic
	ripDot := cat.Rip.CurDot()
	ripDur := ripDot.RemainingDuration(sim)
	roarBuff := cat.SavageRoarBuff
	roarDur := roarBuff.RemainingDuration(sim)
	ripRefreshTime := cat.calcBleedRefreshTime(sim, cat.Rip, ripDot, isExecutePhase, true)
	ripNow := (curCp >= 5) && (!ripDot.IsActive() || ((sim.CurrentTime > ripRefreshTime) && (!isExecutePhase || (cat.Rip.NewSnapshotPower > cat.Rip.CurrentSnapshotPower + 0.001))) || (!isExecutePhase && (roarDur < rotation.RipLeeway) && (ripDot.ExpiresAt() < roarBuff.ExpiresAt() + rotation.RipLeeway))) && (fightDur > ripDot.BaseTickLength) && (!isClearcast || !anyBleedActive) && !cat.shouldDelayBleedRefreshForTf(sim, ripDot, true)

	// Roar logic
	newRoarDur := cat.SavageRoarDurationTable[curCp]
	roarRefreshTime := cat.calcRoarRefreshTime(sim, ripRefreshTime, rotation.RipLeeway, rotation.MinRoarOffset)
	roarNow := (newRoarDur > 0) && (!roarBuff.IsActive() || (sim.CurrentTime > roarRefreshTime))

	// Bite logic
	biteTime := core.TernaryDuration(cat.BerserkCatAura.IsActive(), rotation.BerserkBiteTime, rotation.BiteTime)
	shouldBite := (curCp >= 5) && ripDot.IsActive() && roarBuff.IsActive() && ((rotation.UseBite && (min(ripRefreshTime, roarRefreshTime) - sim.CurrentTime >= biteTime)) || isExecutePhase) && !isClearcast
	shouldEmergencyBite := isExecutePhase && ripDot.IsActive() && (ripDur < ripDot.BaseTickLength) && (curCp >= 1)
	biteNow := shouldBite || shouldEmergencyBite

	// Rake logic
	rakeDot := cat.Rake.CurDot()
	rakeDur := rakeDot.RemainingDuration(sim)
	rakeRefreshTime := cat.calcBleedRefreshTime(sim, cat.Rake, rakeDot, isExecutePhase, false)
	rakeNow := (!rakeDot.IsActive() || (sim.CurrentTime > rakeRefreshTime)) && (fightDur > rakeDot.BaseTickLength) && (!isClearcast || !rakeDot.IsActive() || (rakeDur < time.Second)) && !cat.shouldDelayBleedRefreshForTf(sim, rakeDot, false) && roarBuff.IsActive()

	// Pooling calcs
	ripRefreshPending := ripDot.IsActive() && (ripDur < fightDur - ripDot.BaseTickLength) && (curCp >= core.TernaryInt32(isExecutePhase, 1, 5))
	rakeRefreshPending := rakeDot.IsActive() && (rakeDur < fightDur - rakeDot.BaseTickLength)
	roarRefreshPending := roarBuff.IsActive() && (roarDur < fightDur - cat.ReactionTime) && (newRoarDur > 0)
	rotation.pendingPool.reset()
	rotation.pendingPoolWeaves.reset()

	if ripRefreshPending && (sim.CurrentTime < ripRefreshTime) {
		ripRefreshCost := core.Ternary(isExecutePhase, cat.FerociousBite.DefaultCast.Cost, cat.Rip.DefaultCast.Cost)
		rotation.pendingPool.addAction(ripRefreshTime, ripRefreshCost)
		rotation.pendingPoolWeaves.addAction(ripRefreshTime, ripRefreshCost)
	}

	if rakeRefreshPending && (sim.CurrentTime < rakeRefreshTime) {
		rotation.pendingPool.addAction(rakeRefreshTime, cat.Rake.DefaultCast.Cost)
		rotation.pendingPoolWeaves.addAction(rakeRefreshTime, cat.Rake.DefaultCast.Cost)
	}

	if roarRefreshPending && (sim.CurrentTime < roarRefreshTime) {
		rotation.pendingPool.addAction(roarRefreshTime, cat.SavageRoar.DefaultCast.Cost)
	}

	rotation.pendingPool.sort()
	rotation.pendingPoolWeaves.sort()
	floatingEnergy := rotation.pendingPool.calcFloatingEnergy(cat, sim)
	excessE := curEnergy - floatingEnergy

	// Check bear-weaving conditions.
	furorCap := 100.0 - 1.5 * regenRate
	bearWeaveNow := rotation.BearWeave && cat.canBearWeave(sim, furorCap, regenRate, curEnergy, excessE, rotation.pendingPoolWeaves)

	// Main decision tree starts here.
	var timeToNextAction time.Duration

	if cat.BearFormAura.IsActive() {
		if rotation.shouldTerminateBearWeave(sim, isClearcast, curEnergy, furorCap, regenRate, rotation.pendingPoolWeaves) {
			rotation.readyToShift = true
		} else if cat.ThrashBear.CanCast(sim, cat.CurrentTarget) {
			cat.ThrashBear.Cast(sim, cat.CurrentTarget)
		} else if cat.MangleBear.CanCast(sim, cat.CurrentTarget) {
			cat.MangleBear.Cast(sim, cat.CurrentTarget)
		} else if cat.Lacerate.CanCast(sim, cat.CurrentTarget) {
			cat.Lacerate.Cast(sim, cat.CurrentTarget)
		} else {
			rotation.readyToShift = true
		}

		// Last second Maul check if we are about to shift back.
		if rotation.readyToShift && !isClearcast && cat.Maul.CanCast(sim, cat.CurrentTarget) {
			cat.Maul.Cast(sim, cat.CurrentTarget)
		}

		if !rotation.readyToShift {
			timeToNextAction = cat.ReactionTime
		}
	} else if roarNow {
		if cat.SavageRoar.CanCast(sim, cat.CurrentTarget) {
			cat.SavageRoar.Cast(sim, nil)
			return
		}

		timeToNextAction = core.DurationFromSeconds((cat.CurrentSavageRoarCost() - curEnergy) / regenRate)
	} else if ripNow {
		if cat.Rip.CanCast(sim, cat.CurrentTarget) {
			cat.Rip.Cast(sim, cat.CurrentTarget)
			return
		}

		timeToNextAction = core.DurationFromSeconds((cat.CurrentRipCost() - curEnergy) / regenRate)
	} else if biteNow && ((curEnergy >= cat.CurrentFerociousBiteCost()) || !bearWeaveNow) {
		if cat.FerociousBite.CanCast(sim, cat.CurrentTarget) {
			cat.FerociousBite.Cast(sim, cat.CurrentTarget)
			return
		}

		timeToNextAction = core.DurationFromSeconds((cat.CurrentFerociousBiteCost() - curEnergy) / regenRate)
	} else if rakeNow {
		if cat.Rake.CanCast(sim, cat.CurrentTarget) {
			cat.Rake.Cast(sim, cat.CurrentTarget)
			return
		}

		timeToNextAction = core.DurationFromSeconds((cat.CurrentRakeCost() - curEnergy) / regenRate)
	} else if bearWeaveNow {
		rotation.readyToShift = true
	} else if isClearcast && !cat.ThrashCat.CurDot().IsActive() {
		cat.ThrashCat.Cast(sim, cat.CurrentTarget)
		return
	} else if isClearcast || !ripRefreshPending || !cat.tempSnapshotAura.IsActive() || (ripRefreshTime + cat.ReactionTime - sim.CurrentTime > core.GCDMin) {
		fillerSpell := core.Ternary(rotation.ForceMangleFiller || (!ripDot.IsActive() && (curCp < 5) && !isClearcast), cat.MangleCat, cat.Shred)
		fillerDpc := fillerSpell.ExpectedInitialDamage(sim, cat.CurrentTarget)
		rakeDpc := cat.Rake.ExpectedInitialDamage(sim, cat.CurrentTarget)

		if (fillerDpc < rakeDpc) || (!cat.BerserkCatAura.IsActive() && !isClearcast && (fillerDpc / fillerSpell.DefaultCast.Cost < rakeDpc / cat.Rake.DefaultCast.Cost)) {
			fillerSpell = cat.Rake
		}

		// Force filler on Clearcasts or when about to Energy cap.
		if isClearcast || (curEnergy > cat.MaximumEnergy() - regenRate * cat.ReactionTime.Seconds()) {
			fillerSpell.Cast(sim, cat.CurrentTarget)
			return
		}

		fillerCost := fillerSpell.Cost.GetCurrentCost()
		energyForCalc := core.TernaryFloat64(cat.BerserkCatAura.IsActive(), curEnergy, excessE)

		if energyForCalc >= fillerCost {
			fillerSpell.Cast(sim, cat.CurrentTarget)
			return
		}

		timeToNextAction = core.DurationFromSeconds((fillerCost - energyForCalc) / regenRate)
	}

	nextActionAt := sim.CurrentTime + timeToNextAction
	isPooling, nextRefresh := rotation.pendingPool.nextRefreshTime()

	if isPooling {
		nextActionAt = min(nextActionAt, nextRefresh)
	}

	rotation.ProcessNextPlannedAction(sim, nextActionAt)
}
