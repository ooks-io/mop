package feral

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/druid"
)

func (cat *FeralDruid) applySpecTalents() {
	cat.registerSoulOfTheForest()
}

func (cat *FeralDruid) registerSoulOfTheForest() {
	if !cat.Talents.SoulOfTheForest {
		return
	}

	energyMetrics := cat.NewEnergyMetrics(core.ActionID{SpellID: 114113})

	var cpSnapshot int32

	procSotf := func(sim *core.Simulation) {
		if cpSnapshot > 0 {
			cat.AddEnergy(sim, 4.0 * float64(cpSnapshot), energyMetrics)
			cpSnapshot = 0
		}
	}

	cat.RegisterAura(core.Aura{
		Label:    "Soul of the Forest Trigger",
		Duration: core.NeverExpires,

		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},

		OnApplyEffects: func(aura *core.Aura, _ *core.Simulation, _ *core.Unit, spell *core.Spell) {
			if spell.Matches(druid.DruidSpellFinisher) {
				cpSnapshot = aura.Unit.ComboPoints()
			}
		},

		OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(druid.DruidSpellFinisher) && result.Landed() {
				procSotf(sim)
			}
		},

		OnCastComplete: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(druid.DruidSpellSavageRoar) {
				procSotf(sim)
			}
		},
	})
}
