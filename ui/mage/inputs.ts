import * as InputHelpers from '../core/components/input_helpers';
import { MageArmor } from '../core/proto/mage';
import { ActionId } from '../core/proto_utils/action_id';
import { MageSpecs } from '../core/proto_utils/utils';

// Configuration for class-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.
export const MageArmorInputs = <SpecType extends MageSpecs>() =>
	InputHelpers.makeClassOptionsEnumIconInput<SpecType, MageArmor>({
		fieldName: 'defaultMageArmor',
		values: [
			{ value: MageArmor.MageArmorNone, tooltip: 'No Armor' },
			{ actionId: ActionId.fromSpellId(7302), value: MageArmor.MageArmorFrostArmor },
			{ actionId: ActionId.fromSpellId(6117), value: MageArmor.MageArmorMageArmor },
			{ actionId: ActionId.fromSpellId(30482), value: MageArmor.MageArmorMoltenArmor },
		],
	});
