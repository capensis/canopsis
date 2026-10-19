import { createNumberInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import CNumberField from '@/components/forms/fields/compat/c-number-field.vue';

const stubs = {
  'c-text-field': createNumberInputStub('c-text-field'),
};

const selectTextField = wrapper => wrapper.find('input.c-text-field');

describe('c-number-field', () => {
  const factory = generateShallowRenderer(CNumberField, { stubs });

  it('Value changed after trigger the input', () => {
    const wrapper = factory();
    const textField = selectTextField(wrapper);

    const newValue = Faker.datatype.number();

    textField.setValue(newValue);

    expect(wrapper).toEmitInput(newValue);
  });
});
