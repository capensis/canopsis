import { generateShallowRenderer } from '@unit/utils/vue';

import { COLOR_INDICATOR_TYPES } from '@/constants';

import CColorIndicatorField from '@/components/forms/fields/color/c-color-indicator-field.vue';

const stubs = {
  'c-radio-group': {
    props: ['modelValue', 'value'],
    emits: ['update:modelValue'],
    template: `
      <input
        :value="modelValue ?? value"
        class="c-radio-group"
        @input="$emit('update:modelValue', $event.target.value)"
      />
    `,
  },
};

const selectRadioGroup = wrapper => wrapper.find('input.c-radio-group');

describe('c-color-indicator-field', () => {
  const factory = generateShallowRenderer(CColorIndicatorField, { stubs });

  it('The value did set in the input', () => {
    const wrapper = factory({ propsData: { modelValue: COLOR_INDICATOR_TYPES.state } });

    expect(selectRadioGroup(wrapper).element.value).toBe(COLOR_INDICATOR_TYPES.state);
  });

  it('Value changed after trigger the input', () => {
    const wrapper = factory({ propsData: { modelValue: COLOR_INDICATOR_TYPES.state } });

    selectRadioGroup(wrapper).triggerCustomEvent(
      'update:modelValue',
      COLOR_INDICATOR_TYPES.impactState,
    );

    expect(wrapper).toEmitInput(COLOR_INDICATOR_TYPES.impactState);
  });
});
