import { createNumberInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import { TREE_OF_DEPENDENCIES_SHOW_TYPES } from '@/constants';

import ServiceDependenciesShowTypeField from '@/components/other/service/form/fields/service-dependencies-show-type-field.vue';

const stubs = {
  'widget-settings-item': true,
  'v-radio-group': createNumberInputStub('v-radio-group'),
};

const selectRadioGroup = wrapper => wrapper.find('.v-radio-group');

describe('service-dependencies-show-type-field', () => {
  const factory = generateShallowRenderer(ServiceDependenciesShowTypeField, { stubs });

  test('Value changed after trigger radio group field', () => {
    const wrapper = factory({
      propsData: {
        value: TREE_OF_DEPENDENCIES_SHOW_TYPES.allDependencies,
      },
    });

    selectRadioGroup(wrapper).setValue(
      TREE_OF_DEPENDENCIES_SHOW_TYPES.dependenciesDefiningTheState,
    );

    expect(wrapper).toEmitInput(TREE_OF_DEPENDENCIES_SHOW_TYPES.dependenciesDefiningTheState);
  });
});
