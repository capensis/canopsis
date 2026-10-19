import { createInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import CIdField from '@/components/forms/fields/c-id-field.vue';

const stubs = {
  'v-text-field': createInputStub('v-text-field'),
  'c-help-icon': true,
};

const selectTextField = wrapper => wrapper.find('input.v-text-field');

describe('c-id-field', () => {
  const factory = generateShallowRenderer(CIdField, { stubs });

  test('Value changed after trigger the text field', () => {
    const wrapper = factory({
      propsData: {
        value: '',
      },
    });

    const textField = selectTextField(wrapper);

    const newValue = Faker.datatype.string();

    textField.setValue(newValue);

    expect(wrapper).toEmitInput(newValue);
  });
});
