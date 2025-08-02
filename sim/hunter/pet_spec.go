package hunter

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (hp *HunterPet) ApplySpikedCollar() {
	if hp.hunterOwner.Options.PetSpec != proto.PetSpec_Ferocity {
		return
	}

	basicAttackDamageMod := hp.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  HunterPetFocusDump,
		FloatValue: 0.1,
	})

	critMod := hp.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: 10,
	})

	core.MakePermanent(hp.RegisterAura(core.Aura{
		Label:    "Spiked Collar",
		ActionID: core.ActionID{SpellID: 53184},
		Duration: core.NeverExpires,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			critMod.Activate()
			basicAttackDamageMod.Activate()
			hp.MultiplyMeleeSpeed(sim, 1.1)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			critMod.Deactivate()
			basicAttackDamageMod.Deactivate()
			hp.MultiplyMeleeSpeed(sim, 1/1.1)
		},
	}))
}

func (hp *HunterPet) ApplyCombatExperience() {
	core.MakePermanent(hp.RegisterAura(core.Aura{
		Label:    "Combat Experience",
		ActionID: core.ActionID{SpellID: 20782},
		Duration: core.NeverExpires,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			hp.PseudoStats.DamageDealtMultiplier *= 1.5
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hp.PseudoStats.DamageDealtMultiplier /= 1.5
		},
	}))
}
