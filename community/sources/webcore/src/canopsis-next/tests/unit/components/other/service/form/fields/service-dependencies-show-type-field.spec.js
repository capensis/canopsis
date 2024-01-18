import { TREE_OF_DEPENDENCIES_SHOW_TYPES } from '@/constants';
import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createNumberInputStub } from '@unit/stubs/input';

import ServiceDependenciesShowTypeField from '@/components/other/service/form/fields/service-dependencies-show-type-field.vue';

const stubs = {
  'widget-settings-item': true,
  'v-radio-group': createNumberInputStub('v-radio-group'),
};

const snapshotStubs = {
  'widget-settings-item': true,
};

const selectRadioGroup = wrapper => wrapper.find('.v-radio-group');

describe('service-dependencies-show-type-field', () => {
  const factory = generateShallowRenderer(ServiceDependenciesShowTypeField, { stubs });
  const snapshotFactory = generateRenderer(ServiceDependenciesShowTypeField, { stubs: snapshotStubs });

  test('Value changed after trigger radio group field', () => {
    const wrapper = factory({
      propsData: {
        value: TREE_OF_DEPENDENCIES_SHOW_TYPES.allDependencies,
      },
    });

    selectRadioGroup(wrapper).setValue(TREE_OF_DEPENDENCIES_SHOW_TYPES.dependenciesDefiningTheState);

    expect(wrapper).toEmit('input', TREE_OF_DEPENDENCIES_SHOW_TYPES.dependenciesDefiningTheState);
  });

  test('Renders `service-dependencies-show-type-field` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `service-dependencies-show-type-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: TREE_OF_DEPENDENCIES_SHOW_TYPES.allDependencies,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
