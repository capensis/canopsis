import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import CDescriptionField from '@/components/forms/fields/c-description-field.vue';

const stubs = {
  'v-textarea': true,
};

const selectTextareaNode = wrapper => wrapper.vm.$children[0];

describe('c-description-field', () => {
  const factory = generateShallowRenderer(CDescriptionField, { stubs });

  test('Value changed after trigger the textarea', () => {
    const wrapper = factory({
      propsData: {
        value: '',
      },
    });

    const newValue = Faker.datatype.string();

    selectTextareaNode(wrapper).$emit('input', newValue);

    expect(wrapper).toEmitInput(newValue);
  });
});
