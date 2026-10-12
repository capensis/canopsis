import { createNumberInputStub } from '@unit/stubs/input';
import { randomArrayItem } from '@unit/utils/array';
import { generateShallowRenderer } from '@unit/utils/vue';

import { TREE_OF_DEPENDENCIES_SHOW_TYPES } from '@/constants';

import TreeOfDependenciesSettings from '@/components/sidebars/form/fields/tree-of-dependencies-settings.vue';

const stubs = {
  'widget-settings-item': true,
  'v-radio-group': createNumberInputStub('v-radio-group'),
};

const selectRadioGroup = wrapper => wrapper.find('.v-radio-group');

describe('tree-of-dependencies-settings', () => {
  const factory = generateShallowRenderer(TreeOfDependenciesSettings, { stubs });

  test('Value changed after trigger radio group field', () => {
    const wrapper = factory({
      propsData: {
        value: TREE_OF_DEPENDENCIES_SHOW_TYPES.custom,
      },
    });

    const newValue = randomArrayItem(Object.values(TREE_OF_DEPENDENCIES_SHOW_TYPES));

    selectRadioGroup(wrapper).setValue(newValue);

    expect(wrapper).toEmitInput(newValue);
  });
});
