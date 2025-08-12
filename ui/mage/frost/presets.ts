import { Encounter } from '../../core/encounter';
import * as PresetUtils from '../../core/preset_utils';
import { ConsumesSpec, Glyphs, Profession, Race, Stat } from '../../core/proto/common';
import { FrostMage_Options as MageOptions, MageMajorGlyph, MageMinorGlyph, MageArmor } from '../../core/proto/mage';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import FrostApl from './apls/frost.apl.json';
import FrostAoeApl from './apls/frost_aoe.apl.json';
import FrostCleaveApl from './apls/frost_cleave.apl.json';
import P1PreBISRealisticGear from './gear_sets/p1_prebis_realistic.gear.json';
import P1PreBISGear from './gear_sets/p1_prebis.gear.json';
import P1PostMSVGear from './gear_sets/p1_post_msv.gear.json';
import P1PostHOFGear from './gear_sets/p1_post_hof.gear.json';
import P1BISGear from './gear_sets/p1_bis.gear.json';
// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const P1_PREBIS = PresetUtils.makePresetGear('P1 - Pre-BIS', P1PreBISGear);
export const P1_PREBIS_REALISTIC = PresetUtils.makePresetGear('P1 - Pre-BIS (Realistic)', P1PreBISRealisticGear);
export const P1_POST_MSV = PresetUtils.makePresetGear('P1 - Post-MSV', P1PostMSVGear);
export const P1_POST_HOF = PresetUtils.makePresetGear('P1 - Post-HoF', P1PostHOFGear);
export const P1_BIS = PresetUtils.makePresetGear('P1 - BIS', P1BISGear);

export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Default', FrostApl);
export const ROTATION_PRESET_AOE = PresetUtils.makePresetAPLRotation('AOE', FrostAoeApl);
// export const ROTATION_PRESET_CLEAVE = PresetUtils.makePresetAPLRotation('Cleave', FrostCleaveApl);

// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1',
	Stats.fromMap({
		[Stat.StatIntellect]: 1.00,
		[Stat.StatSpellPower]: 0.80,
		[Stat.StatHitRating]: 1.15,
		[Stat.StatCritRating]: 0.47,
		[Stat.StatHasteRating]: 0.45,
		[Stat.StatMasteryRating]: 0.39,
	}),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const FrostDefaultTalents = {
	name: 'Default',
	data: SavedTalents.create({
		talentsString: '311122',
		glyphs: Glyphs.create({
			major1: MageMajorGlyph.GlyphOfSplittingIce,
			major2: MageMajorGlyph.GlyphOfIcyVeins,
			major3: MageMajorGlyph.GlyphOfWaterElemental,
			minor1: MageMinorGlyph.GlyphOfMomentum,
			minor2: MageMinorGlyph.GlyphOfMirrorImage,
			minor3: MageMinorGlyph.GlyphOfTheUnboundElemental,
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
			minor3: MageMinorGlyph.GlyphOfTheUnboundElemental,
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
			minor3: MageMinorGlyph.GlyphOfTheUnboundElemental,
		}),
	}),
};

export const DefaultFrostOptions = MageOptions.create({
	classOptions: {
		defaultMageArmor: MageArmor.MageArmorFrostArmor,
	},
});

export const ENCOUNTER_SINGLE_TARGET = PresetUtils.makePresetEncounter('Single Target', Encounter.defaultEncounterProto());
export const ENCOUNTER_CLEAVE = PresetUtils.makePresetEncounter('Cleave', Encounter.defaultEncounterProto(2));
export const ENCOUNTER_AOE = PresetUtils.makePresetEncounter('AoE (5+)', Encounter.defaultEncounterProto(5));

export const P1_PRESET_BUILD_DEFAULT = PresetUtils.makePresetBuild('Single Target', {
	talents: FrostDefaultTalents,
	rotation: ROTATION_PRESET_DEFAULT,
	encounter: ENCOUNTER_SINGLE_TARGET,
});

export const P1_PRESET_BUILD_CLEAVE = PresetUtils.makePresetBuild('Cleave', {
	talents: FrostTalentsCleave,
	rotation: ROTATION_PRESET_DEFAULT,
	encounter: ENCOUNTER_CLEAVE,
});

export const P1_PRESET_BUILD_AOE = PresetUtils.makePresetBuild('AoE (5+)', {
	talents: FrostTalentsAoE,
	rotation: ROTATION_PRESET_AOE,
	encounter: ENCOUNTER_AOE,
});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
	race: Race.RaceTroll,
};
