import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import CColumnSizeField from '@/components/forms/fields/c-column-size-field.vue';

const stubs = {
  'c-radio-group': {
    props: ['modelValue', 'value'],
    emits: ['update:modelValue'],
    template: `
      <input
        :value="modelValue ?? value"
        class="c-radio-group"
        @input="$emit('update:modelValue', +$event.target.value)"
      />
    `,
  },
};

const selectRadioGroup = wrapper => wrapper.find('.c-radio-group');

describe('c-column-size-field', () => {
  const factory = generateShallowRenderer(CColumnSizeField, { stubs });

  test('Renders `c-column-size-field` with required props', () => {
    const wrapper = factory();
    const newSize = Faker.datatype.number({ min: 1, max: 12 });

    selectRadioGroup(wrapper).triggerCustomEvent('update:modelValue', newSize);

    expect(wrapper).toEmitInput(newSize);
  });
});
