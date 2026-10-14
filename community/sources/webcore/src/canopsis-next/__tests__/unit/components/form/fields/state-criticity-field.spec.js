import { generateShallowRenderer } from '@unit/utils/vue';

import { ALARM_STATES } from '@/constants';

import StateCriticityField from '@/components/forms/fields/state-criticity-field.vue';

const stubs = {
  'v-btn-toggle': {
    props: ['modelValue', 'value'],
    emits: ['update:modelValue'],
    template: `
      <input
        :value="modelValue ?? value"
        class="v-btn-toggle"
        @input="$emit('update:modelValue', +$event.target.value)"
      />
    `,
  },
};

const selectBtnToggle = wrapper => wrapper.find('input.v-btn-toggle');

describe('state-criticity-field', () => {
  const factory = generateShallowRenderer(StateCriticityField, { stubs });

  it('Value changed after trigger click on the button', () => {
    const wrapper = factory({
      propsData: {
        modelValue: ALARM_STATES.major,
      },
    });

    selectBtnToggle(wrapper).triggerCustomEvent('update:modelValue', ALARM_STATES.ok);

    expect(wrapper).toEmitInput(ALARM_STATES.ok);
  });
});
