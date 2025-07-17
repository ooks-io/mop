package feral

import (
	"time"

	"github.com/wowsims/mop/sim/core"
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
