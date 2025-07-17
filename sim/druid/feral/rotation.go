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
		rotation.BiteTime = time.Second * 11
		rotation.BerserkBiteTime = time.Second * 6
		rotation.MinRoarOffset = time.Second * 31
		rotation.RipLeeway = time.Second * 1
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
}
