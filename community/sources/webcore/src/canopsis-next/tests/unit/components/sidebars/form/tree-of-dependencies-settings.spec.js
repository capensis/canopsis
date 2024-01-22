import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { randomArrayItem } from '@unit/utils/array';
import { createNumberInputStub } from '@unit/stubs/input';

import { TREE_OF_DEPENDENCIES_SHOW_TYPES } from '@/constants';

import TreeOfDependenciesSettings from '@/components/sidebars/form/fields/tree-of-dependencies-settings.vue';

const stubs = {
  'widget-settings-item': true,
  'v-radio-group': createNumberInputStub('v-radio-group'),
};

const snapshotStubs = {
  'widget-settings-item': true,
};

const selectRadioGroup = wrapper => wrapper.find('.v-radio-group');

describe('tree-of-dependencies-settings', () => {
  const factory = generateShallowRenderer(TreeOfDependenciesSettings, { stubs });
  const snapshotFactory = generateRenderer(TreeOfDependenciesSettings, { stubs: snapshotStubs });

  test('Value changed after trigger radio group field', () => {
    const wrapper = factory({
      propsData: {
        value: TREE_OF_DEPENDENCIES_SHOW_TYPES.custom,
      },
    });

    const newValue = randomArrayItem(Object.values(TREE_OF_DEPENDENCIES_SHOW_TYPES));

    selectRadioGroup(wrapper).setValue(newValue);

    expect(wrapper).toEmit('input', newValue);
  });

  test('Renders `tree-of-dependencies-settings` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `tree-of-dependencies-settings` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: TREE_OF_DEPENDENCIES_SHOW_TYPES.allDependencies,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
