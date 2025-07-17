import { Encounter } from '../../core/encounter';
import * as PresetUtils from '../../core/preset_utils';
import { Class, ConsumesSpec, Debuffs, Glyphs, Profession, Race, RaidBuffs, Stat } from '../../core/proto/common';
import { FrostMage_Options as MageOptions, MageMajorGlyph, MageMinorGlyph } from '../../core/proto/mage';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import { defaultRaidBuffMajorDamageCooldowns } from '../../core/proto_utils/utils';
import FrostApl from './apls/frost.apl.json';
import FrostAoeApl from './apls/frost_aoe.apl.json';
import FrostCleaveApl from './apls/frost_cleave.apl.json';
import P1BISGear from './gear_sets/p1_bis.gear.json';
import P1PreBISPoorGear from './gear_sets/p1_prebis_poor.gear.json';
import P1PreBISRichGear from './gear_sets/p1_prebis_rich.gear.json';
// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const P1_PREBIS_RICH = PresetUtils.makePresetGear('P1 - Pre-BIS (Rich)', P1PreBISRichGear);
export const P1_PREBIS_POOR = PresetUtils.makePresetGear('P1 - Pre-BIS (Budget)', P1PreBISPoorGear);

export const P1_BIS = PresetUtils.makePresetGear('P1 - BIS', P1BISGear);

export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Frost ST', FrostApl);
export const ROTATION_PRESET_AOE = PresetUtils.makePresetAPLRotation('Frost AOE', FrostAoeApl);
export const ROTATION_PRESET_CLEAVE = PresetUtils.makePresetAPLRotation('Frost Cleave', FrostCleaveApl);

// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Frost P1',
	Stats.fromMap({
		[Stat.StatIntellect]: 1.23,
		[Stat.StatSpellPower]: 1,
		[Stat.StatHitRating]: 1.15,
		[Stat.StatCritRating]: 0.49,
		[Stat.StatHasteRating]: 0.60,
		[Stat.StatMasteryRating]: 0.47,
	}),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const FrostDefaultTalents = {
	name: 'Frost',
	data: SavedTalents.create({
		talentsString: '311122',
		glyphs: Glyphs.create({
			major1: MageMajorGlyph.GlyphOfSplittingIce,
			major2: MageMajorGlyph.GlyphOfIcyVeins,
			major3: MageMajorGlyph.GlyphOfWaterElemental,
			minor1: MageMinorGlyph.GlyphOfMomentum,
			minor2: MageMinorGlyph.GlyphOfMirrorImage,
			minor3: MageMinorGlyph.GlyphOfTheUnboundElemental
		}),
	}),
};

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 76085, // Flask of the Warm Sun
	foodId: 74650, // Mogu Fish Stew
	potId: 76093, // Potion of the Jade Serpent
	prepotId: 76093, // Potion of the Jade Serpent
});

export const FrostTalentsCleave = {
	name: 'Cleave',
	data: SavedTalents.create({
		talentsString: '311112',
		glyphs: Glyphs.create({
			major1: MageMajorGlyph.GlyphOfSplittingIce,
			major2: MageMajorGlyph.GlyphOfIcyVeins,
			major3: MageMajorGlyph.GlyphOfWaterElemental,
			minor1: MageMinorGlyph.GlyphOfMomentum,
			minor2: MageMinorGlyph.GlyphOfMirrorImage,
			minor3: MageMinorGlyph.GlyphOfTheUnboundElemental
		}),
	}),
};

export const FrostTalentsAoE = {
	name: 'AoE (5+)',
	data: SavedTalents.create({
		talentsString: '311112',
		glyphs: Glyphs.create({
			major1: MageMajorGlyph.GlyphOfSplittingIce,
			major2: MageMajorGlyph.GlyphOfIcyVeins,
			major3: MageMajorGlyph.GlyphOfWaterElemental,
			minor1: MageMinorGlyph.GlyphOfMomentum,
			minor2: MageMinorGlyph.GlyphOfMirrorImage,
			minor3: MageMinorGlyph.GlyphOfTheUnboundElemental
		}),
	}),
};

export const DefaultFrostOptions = MageOptions.create({
	classOptions: {},
});

export const ENCOUNTER_SINGLE_TARGET = PresetUtils.makePresetEncounter('Frost ST', Encounter.defaultEncounterProto());
export const ENCOUNTER_CLEAVE = PresetUtils.makePresetEncounter('Frost Cleave', Encounter.defaultEncounterProto(2));
export const ENCOUNTER_AOE = PresetUtils.makePresetEncounter('Frost AoE (5+)', Encounter.defaultEncounterProto(5));

export const P1_PRESET_BUILD_DEFAULT = PresetUtils.makePresetBuild('Frost ST', {
	talents: FrostDefaultTalents,
	rotation: ROTATION_PRESET_DEFAULT,
	encounter: ENCOUNTER_SINGLE_TARGET,
	epWeights: P1_EP_PRESET,
});

export const P1_PRESET_BUILD_CLEAVE = PresetUtils.makePresetBuild('Frost Cleave', {
	talents: FrostTalentsCleave,
	rotation: ROTATION_PRESET_CLEAVE,
	encounter: ENCOUNTER_CLEAVE,
	epWeights: P1_EP_PRESET,
});

export const P1_PRESET_BUILD_AOE = PresetUtils.makePresetBuild('Frost AoE (5+)', {
	talents: FrostTalentsAoE,
	rotation: ROTATION_PRESET_AOE,
	encounter: ENCOUNTER_AOE,
	epWeights: P1_EP_PRESET,
});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
	race: Race.RaceTroll,
};
