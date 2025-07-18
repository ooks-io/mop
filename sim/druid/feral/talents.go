package feral

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/druid"
)

func (cat *FeralDruid) applySpecTalents() {
	cat.registerSoulOfTheForest()
	cat.registerIncarnation()
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

func (cat *FeralDruid) registerIncarnation() {
	if !cat.Talents.Incarnation {
		return
	}

	actionID := core.ActionID{SpellID: 102543}

	var oldExtraCastCondition core.CanCastCondition

	cat.IncarnationAura = cat.RegisterAura(core.Aura{
		Label:    "Incarnation: King of the Jungle",
		ActionID: actionID,
		Duration: time.Second * 30,

		OnGain: func(_ *core.Aura, _ *core.Simulation) {
			oldExtraCastCondition = cat.Ravage.ExtraCastCondition
			cat.Ravage.ExtraCastCondition = nil
		},

		OnExpire: func(_ *core.Aura, _ *core.Simulation) {
			cat.Ravage.ExtraCastCondition = oldExtraCastCondition
		},
	})

	cat.Incarnation = cat.RegisterSpell(druid.Any, core.SpellConfig{
		ActionID:        actionID,
		Flags:           core.SpellFlagAPL,
		RelatedSelfBuff: cat.IncarnationAura,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},

			CD: core.Cooldown{
				Timer:    cat.NewTimer(),
				Duration: time.Minute * 3,
			},

			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if !cat.InForm(druid.Cat) {
				cat.CatFormAura.Activate(sim)
			}

			cat.IncarnationAura.Activate(sim)
		},
	})

	cat.AddMajorCooldown(core.MajorCooldown{
		Spell: cat.Incarnation.Spell,
		Type:  core.CooldownTypeDPS,

		ShouldActivate: func(sim *core.Simulation, _ *core.Character) bool {
			return cat.BerserkCatAura.IsActive() && !cat.ClearcastingAura.IsActive() && (cat.CurrentEnergy() + cat.EnergyRegenPerSecond() < 100)
		},
	})
}
