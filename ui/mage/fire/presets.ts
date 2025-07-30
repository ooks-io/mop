import { Encounter } from '../../core/encounter';
import * as PresetUtils from '../../core/preset_utils';
import { ConsumesSpec, Glyphs, Profession, PseudoStat, Race, Stat } from '../../core/proto/common';
import { MageArmor, FireMage_Options as MageOptions, MageMajorGlyph as MajorGlyph, MageMinorGlyph as MinorGlyph } from '../../core/proto/mage';
import { SavedTalents } from '../../core/proto/ui';
import { Stats, UnitStat, UnitStatPresets } from '../../core/proto_utils/stats';
import FireApl from './apls/fire.apl.json';
import FireCleaveApl from './apls/fire_cleave.apl.json';
import P1PreBISGear from './gear_sets/p1_prebis.gear.json';
import P1PostMSVGear from './gear_sets/p1_post_msv.gear.json';
import P1PostHOFGear from './gear_sets/p1_post_hof.gear.json';
import P1BISGear from './gear_sets/p1_bis.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
export const P1_PREBIS = PresetUtils.makePresetGear('P1 - Pre-BIS', P1PreBISGear);
export const P1_POST_MSV = PresetUtils.makePresetGear('P1 - Post-MSV', P1PostMSVGear);
export const P1_POST_HOF = PresetUtils.makePresetGear('P1 - Post-HoF', P1PostHOFGear);
export const P1_BIS = PresetUtils.makePresetGear('P1 - BIS', P1BISGear);


// export const P1TrollDefaultSimpleRotation = FireMage_Rotation.create({
// 	combustThreshold: 515000,
// 	combustLastMomentLustPercentage: 140000,
// 	combustNoLustPercentage: 260000,
// });

// export const P1NoTrollDefaultSimpleRotation = FireMage_Rotation.create({
// 	combustThreshold: 470000,
// 	combustLastMomentLustPercentage: 115000,
// 	combustNoLustPercentage: 225000,
// });

// export const P1_SIMPLE_ROTATION_DEFAULT = PresetUtils.makePresetSimpleRotation('P1 - Default', Spec.SpecFireMage, P1TrollDefaultSimpleRotation);
// export const P1_SIMPLE_ROTATION_NO_TROLL = PresetUtils.makePresetSimpleRotation('P1 - Not Troll', Spec.SpecFireMage, P1NoTrollDefaultSimpleRotation);

//export const ROTATION_PRESET_SIMPLE = PresetUtils.makePresetSimpleRotation('Simple Default', Spec.SpecFireMage, DefaultSimpleRotation);
export const FIRE_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Default', FireApl);

// export const FIRE_ROTATION_PRESET_CLEAVE = PresetUtils.makePresetAPLRotation('Cleave', FireCleaveApl);

// Preset options for EP weights
export const DEFAULT_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Default',
	Stats.fromMap({
		[Stat.StatIntellect]: 1.37,
		[Stat.StatSpellPower]: 1.0,
		[Stat.StatHitRating]: 1.21,
		[Stat.StatCritRating]: 0.88,
		[Stat.StatHasteRating]: 0.73,
		[Stat.StatMasteryRating]: 0.73,
	}),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const FireTalents = {
	name: 'Default',
	data: SavedTalents.create({
		talentsString: '111122',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfCombustion,
			major2: MajorGlyph.GlyphOfInfernoBlast,
			major3: MajorGlyph.GlyphOfRapidDisplacement,
			minor1: MinorGlyph.GlyphOfMomentum,
			minor2: MinorGlyph.GlyphOfLooseMana,
		}),
	}),
};

export const FireTalentsCleave = {
	name: 'Cleave',
	data: SavedTalents.create({
		talentsString: '111112',
		glyphs: Glyphs.create({
			...FireTalents.data.glyphs,
		}),
	}),
};

export const DefaultFireOptions = MageOptions.create({
	classOptions: {
		defaultMageArmor: MageArmor.MageArmorMoltenArmor,
	},
});

export const DefaultFireConsumables = ConsumesSpec.create({
	flaskId: 76085, // Flask of the Warm Sun
	foodId: 74650, // Mogu Fish Stew
	potId: 76093, // Potion of the Jade Serpent
	prepotId: 76093, // Potion of the Jade Serpent
});

export const ENCOUNTER_SINGLE_TARGET = PresetUtils.makePresetEncounter('Single Target', Encounter.defaultEncounterProto());
export const ENCOUNTER_CLEAVE = PresetUtils.makePresetEncounter('Cleave (3 targets)', Encounter.defaultEncounterProto(3));

export const P1_PRESET_BUILD_DEFAULT = PresetUtils.makePresetBuild('Single Target', {
	talents: FireTalents,
	rotation: FIRE_ROTATION_PRESET_DEFAULT,
	encounter: ENCOUNTER_SINGLE_TARGET,
});

export const P1_PRESET_BUILD_CLEAVE = PresetUtils.makePresetBuild('Cleave (3 targets)', {
	talents: FireTalentsCleave,
	rotation: FIRE_ROTATION_PRESET_DEFAULT,
	encounter: ENCOUNTER_CLEAVE,
});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
	race: Race.RaceTroll,
};

export const COMBUSTION_BREAKPOINT: UnitStatPresets = {
	unitStat: UnitStat.fromPseudoStat(PseudoStat.PseudoStatSpellHastePercent),
	presets: new Map([
		['11-tick - Combust', 4.986888],
		['12-tick - Combust', 15.008639],
		['13-tick - Combust', 25.07819],
		['14-tick - Combust', 35.043908],
		['15-tick - Combust', 45.032653],
		['16-tick - Combust', 54.918692],
		['17-tick - Combust', 64.880489],
		['18-tick - Combust', 74.978158],
		['19-tick - Combust', 85.01391],
		['20-tick - Combust', 95.121989],
		['21-tick - Combust', 105.128247],
		['22-tick - Combust', 114.822817],
		['23-tick - Combust', 124.971929],
		['24-tick - Combust', 135.017682],
		['25-tick - Combust', 144.798102],
		['26-tick - Combust', 154.777135],
		['27-tick - Combust', 164.900732],
		['28-tick - Combust', 175.103239],
		['29-tick - Combust', 185.306786],
	]),
};

export const GLYPHED_COMBUSTION_BREAKPOINT: UnitStatPresets = {
	unitStat: UnitStat.fromPseudoStat(PseudoStat.PseudoStatSpellHastePercent),
	presets: new Map([
		['21-tick - Combust (Glyph)', 2.511543],
		['22-tick - Combust (Glyph)', 7.469114],
		['23-tick - Combust (Glyph)', 12.549253],
		['24-tick - Combust (Glyph)', 17.439826],
		['25-tick - Combust (Glyph)', 22.473989],
		['26-tick - Combust (Glyph)', 27.469742],
		['27-tick - Combust (Glyph)', 32.538122],
		['28-tick - Combust (Glyph)', 37.457064],
		['29-tick - Combust (Glyph)', 42.551695],
		['30-tick - Combust (Glyph)', 47.601498],
		['31-tick - Combust (Glyph)', 52.555325],
		['32-tick - Combust (Glyph)', 57.604438],
		['33-tick - Combust (Glyph)', 62.469563],
		['34-tick - Combust (Glyph)', 67.364045],
		['35-tick - Combust (Glyph)', 72.562584],
		['36-tick - Combust (Glyph)', 77.462321],
		['37-tick - Combust (Glyph)', 82.648435],
		['38-tick - Combust (Glyph)', 87.44146],
		['39-tick - Combust (Glyph)', 92.492819],
	]),
};
