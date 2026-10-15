import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { COLOR_INDICATOR_TYPES } from '@/constants';

import FieldColorIndicator from '@/components/sidebars/form/fields/color-indicator.vue';

const stubs = {
  'widget-settings-item': true,
  'c-color-indicator-field': true,
};

const selectNumberField = wrapper => wrapper.find('c-color-indicator-field-stub');

describe('field-color-indicator', () => {
  const factory = generateShallowRenderer(FieldColorIndicator, { stubs });
  const snapshotFactory = generateRenderer(FieldColorIndicator, { stubs });

  test('Value changed after trigger color indicator field', () => {
    const wrapper = factory();

    selectNumberField(wrapper).triggerCustomEvent('input', COLOR_INDICATOR_TYPES.impactState);

    expect(wrapper).toEmitInput(COLOR_INDICATOR_TYPES.impactState);
  });

  test('Renders `field-color-indicator` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `field-color-indicator` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: COLOR_INDICATOR_TYPES.impactState,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
