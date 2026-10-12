import { generateShallowRenderer } from '@unit/utils/vue';

import CEnabledField from '@/components/forms/fields/c-enabled-field.vue';

const stubs = {
  'v-switch': {
    props: ['modelValue', 'value'],
    emits: ['update:modelValue'],
    template: `
      <input
        :checked="modelValue ?? value"
        type="checkbox"
        class="v-switch"
        @change="$emit('update:modelValue', $event.target.checked)"
      />
    `,
  },
};

const selectSwitchField = wrapper => wrapper.find('input.v-switch');

describe('c-enabled-field', () => {
  const factory = generateShallowRenderer(CEnabledField, { stubs });

  it('Value set in the input', () => {
    const wrapper = factory({ propsData: { modelValue: false } });
    const input = selectSwitchField(wrapper);

    expect(input.element.checked).toBe(false);
  });

  it('Value changed in the input', () => {
    const wrapper = factory({ propsData: { modelValue: false } });

    selectSwitchField(wrapper).triggerCustomEvent('update:modelValue', true);

    expect(wrapper).toEmitInput(true);
  });
});
