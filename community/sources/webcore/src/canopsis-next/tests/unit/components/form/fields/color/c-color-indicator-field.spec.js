import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import CColorIndicatorField from '@/components/forms/fields/color/c-color-indicator-field.vue';
import { COLOR_INDICATOR_TYPES } from '@/constants';

const stubs = {
  'v-radio-group': {
    props: ['value'],
    template: `
      <input :value="value" class="v-radio-group" @input="$listeners.input($event.target.value)" />
    `,
  },
};

describe('c-color-indicator-field', () => {
  const factory = generateShallowRenderer(CColorIndicatorField, { stubs });
  const snapshotFactory = generateRenderer(CColorIndicatorField);

  it('The value did set in the input', () => {
    const wrapper = factory({ propsData: { value: COLOR_INDICATOR_TYPES.state } });
    const input = wrapper.find('input.v-radio-group');

    expect(input.element.value).toBe(COLOR_INDICATOR_TYPES.state);
  });

  it('Value changed after trigger the input', () => {
    const wrapper = factory({ propsData: { value: COLOR_INDICATOR_TYPES.state } });
    const input = wrapper.find('input.v-radio-group');

    input.setValue(COLOR_INDICATOR_TYPES.impactState);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);
    expect(inputEvents[0]).toEqual([COLOR_INDICATOR_TYPES.impactState]);
  });

  it('Renders `c-color-indicator-field` with state value correctly', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-color-indicator-field` with impact state value correctly', () => {
    const wrapper = snapshotFactory({

      propsData: {
        value: COLOR_INDICATOR_TYPES.impactState,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
