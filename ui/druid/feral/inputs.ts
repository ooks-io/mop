import * as InputHelpers from '../../core/components/input_helpers.js';
import { Player } from '../../core/player.js';
import { APLRotation_Type } from '../../core/proto/apl.js';
import { Spec } from '../../core/proto/common.js';
import { FeralDruid_Rotation_AplType as AplType } from '../../core/proto/druid.js';
import { TypedEvent } from '../../core/typed_event.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const AssumeBleedActive = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecFeralDruid>({
	fieldName: 'assumeBleedActive',
	label: 'Assume Bleed Always Active',
	labelTooltip: "Assume bleed always exists for 'Rend and Tear' calculations. Otherwise will only calculate based on own rip/rake/lacerate.",
	extraCssClasses: ['within-raid-sim-hide'],
});

export const CannotShredTarget = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecFeralDruid>({
	fieldName: 'cannotShredTarget',
	label: 'Cannot Shred Target',
	labelTooltip: 'Alternative to "In Front of Target" for modeling bosses that do not Parry or Block, but which you still cannot Shred.',
});

function ShouldShowAdvParamST(player: Player<Spec.SpecFeralDruid>): boolean {
	const rot = player.getSimpleRotation();
	return rot.manualParams && rot.rotationType == AplType.SingleTarget;
}

function ShouldShowAdvParamAoe(player: Player<Spec.SpecFeralDruid>): boolean {
	const rot = player.getSimpleRotation();
	return rot.manualParams && rot.rotationType == AplType.Aoe;
}

export const FeralDruidRotationConfig = {
	inputs: [
		InputHelpers.makeRotationEnumInput<Spec.SpecFeralDruid, AplType>({
			fieldName: 'rotationType',
			label: 'Type',
			values: [
				{ name: 'Single Target', value: AplType.SingleTarget },
				// { name: 'AOE', value: AplType.Aoe },
			],
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'bearWeave',
			label: 'Enable bear-weaving',
			labelTooltip: 'Weave into Bear Form while pooling Energy',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'snekWeave',
			label: 'Use Albino Snake',
			labelTooltip: 'Reset swing timer at the end of bear-weaves using Albino Snake pet',
			showWhen: (player: Player<Spec.SpecFeralDruid>) => player.getSimpleRotation().bearWeave,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'useNs',
			label: "Use Nature's Swiftness",
			labelTooltip: "Use Nature's Swiftness to fill gaps in Predatory Swiftness uptime",
			showWhen: (player: Player<Spec.SpecFeralDruid>) => player.getTalents().dreamOfCenarius,
			changeEmitter: (player: Player<Spec.SpecFeralDruid>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'allowAoeBerserk',
			label: 'Allow AoE Berserk',
			labelTooltip: 'Allow Berserk usage in AoE rotation',
			showWhen: (player: Player<Spec.SpecFeralDruid>) => player.getSimpleRotation().rotationType == AplType.Aoe,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'manualParams',
			label: 'Manual Advanced Parameters',
			labelTooltip: 'Manually specify advanced parameters, otherwise will use preset defaults',
			showWhen: (player: Player<Spec.SpecFeralDruid>) => player.getSimpleRotation().rotationType == AplType.SingleTarget,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralDruid>({
			fieldName: 'minRoarOffset',
			label: 'Roar Offset',
			labelTooltip: 'Minimum offset gained between current Rip and new Roar to allow a Roar clip',
			showWhen: ShouldShowAdvParamST,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralDruid>({
			fieldName: 'ripLeeway',
			label: 'Rip Leeway',
			labelTooltip: 'Minimum tolerated gap between current Rip and current Roar expiration times before kicking in emergency measures',
			showWhen: ShouldShowAdvParamST,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'useBite',
			label: 'Bite during rotation',
			labelTooltip: 'Use Bite during rotation rather than just for Rip maintenance during Execute',
			showWhen: ShouldShowAdvParamST,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralDruid>({
			fieldName: 'biteTime',
			label: 'Bite Time',
			labelTooltip: 'Minimum seconds remaining before Rip or Roar should ideally be refreshed (including planned early clips) to allow a Bite',
			showWhen: (player: Player<Spec.SpecFeralDruid>) =>
				ShouldShowAdvParamST(player) && player.getSimpleRotation().useBite,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralDruid>({
			fieldName: 'berserkBiteTime',
			label: 'Bite Time during Berserk',
			labelTooltip: 'More aggressive threshold when Berserk is active',
			showWhen: (player: Player<Spec.SpecFeralDruid>) =>
				ShouldShowAdvParamST(player) && player.getSimpleRotation().useBite,
		}),
	],
};
