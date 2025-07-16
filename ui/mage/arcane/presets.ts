import { Encounter } from '../../core/encounter';
import * as PresetUtils from '../../core/preset_utils';
import { Class, ConsumesSpec, Debuffs, Glyphs, Profession, Race, RaidBuffs, Stat } from '../../core/proto/common';
import { ArcaneMage_Options as MageOptions, MageMajorGlyph as MajorGlyph, MageMinorGlyph } from '../../core/proto/mage';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import { defaultRaidBuffMajorDamageCooldowns } from '../../core/proto_utils/utils';
import ArcaneCleaveApl from './apls/arcane_cleave.apl.json';
import ArcaneApl from './apls/default.apl.json';
import P1ArcaneBisGear from './gear_sets/p1_bis.gear.json';
import ArcanePreBisGear from './gear_sets/prebis.gear.json';
import ArcaneRichPreBisGear from './gear_sets/rich_prebis.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
export const P1_BIS_PRESET = PresetUtils.makePresetGear('P1 Heroic Pre-BIS', ArcanePreBisGear);
export const PREBIS_PRESET = PresetUtils.makePresetGear('P1 BIS', P1ArcaneBisGear);
export const RICH_PREBIS_PRESET = PresetUtils.makePresetGear('P1 Rich Pre-BIS', ArcaneRichPreBisGear);

export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Arcane ST', ArcaneApl);
export const ROTATION_PRESET_CLEAVE = PresetUtils.makePresetAPLRotation('Arcane Cleave', ArcaneCleaveApl);

// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Default',
	Stats.fromMap({
		[Stat.StatIntellect]: 1.24,
		[Stat.StatSpellPower]: 1,
		[Stat.StatHitRating]: 1.31,
		[Stat.StatCritRating]: 0.53,
		[Stat.StatHasteRating]: 0.7,
		[Stat.StatMasteryRating]: 0.68,
	}),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/mop-classic/talent-calc and copy the numbers in the url.
export const ArcaneTalents = {
	name: 'Arcane',
	data: SavedTalents.create({
		talentsString: '311122',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfArcanePower,
			major2: MajorGlyph.GlyphOfRapidDisplacement,
			major3: MajorGlyph.GlyphOfConeOfCold,
			minor1: MageMinorGlyph.GlyphOfMomentum,
			minor2: MageMinorGlyph.GlyphOfRapidTeleportation,
			minor3: MageMinorGlyph.GlyphOfLooseMana,
		}),
	}),
};

export const ArcaneTalentsCleave = {
	name: 'Cleave',
	data: SavedTalents.create({
		talentsString: '311112',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfArcanePower,
			major2: MajorGlyph.GlyphOfRapidDisplacement,
			major3: MajorGlyph.GlyphOfConeOfCold,
			minor1: MageMinorGlyph.GlyphOfMomentum,
			minor2: MageMinorGlyph.GlyphOfRapidTeleportation,
			minor3: MageMinorGlyph.GlyphOfLooseMana,
		}),
	}),
};

// Default buffs and debuffs
export const DefaultRaidBuffs = RaidBuffs.create({
	...defaultRaidBuffMajorDamageCooldowns(Class.ClassMage),
	blessingOfKings: true,
	leaderOfThePack: true,
	serpentsSwiftness: true,
	bloodlust: true,
});

export const DefaultDebuffs = Debuffs.create({
	curseOfElements: true,
});

// Encounter presets
export const ENCOUNTER_SINGLE_TARGET = PresetUtils.makePresetEncounter('Arcane ST', Encounter.defaultEncounterProto());
export const ENCOUNTER_CLEAVE = PresetUtils.makePresetEncounter('Arcane Cleave (2 targets)', Encounter.defaultEncounterProto(2));

export const P1_PRESET_BUILD_DEFAULT = PresetUtils.makePresetBuild('Arcane ST', {
	talents: ArcaneTalents,
	rotation: ROTATION_PRESET_DEFAULT,
	encounter: ENCOUNTER_SINGLE_TARGET,
	epWeights: P1_EP_PRESET,
});

export const P1_PRESET_BUILD_CLEAVE = PresetUtils.makePresetBuild('Arcane Cleave (2 targets)', {
	talents: ArcaneTalentsCleave,
	rotation: ROTATION_PRESET_CLEAVE,
	encounter: ENCOUNTER_CLEAVE,
	epWeights: P1_EP_PRESET,
});

export const DefaultArcaneOptions = MageOptions.create({
	classOptions: {},
});
export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 76085, // Flask of the Warm Sun
	foodId: 74650, // Mogu Fish Stew
	potId: 76093, // Potion of the Jade Serpent
	prepotId: 76093, // Potion of the Jade Serpent
});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
	race: Race.RaceTroll,
};
