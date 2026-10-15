import { createInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import { PATTERN_FIELD_TYPES } from '@/constants';

import CMixedInputField from '@/components/forms/fields/compat/c-mixed-input-field.vue';

const stubs = {
  'c-text-field': createInputStub('c-text-field'),
  'c-combobox': createInputStub('c-combobox'),
  'c-array-text-field': true,
};

const selectTextField = wrapper => wrapper.find('input.c-text-field');
const selectCombobox = wrapper => wrapper.find('input.c-combobox');

describe('c-mixed-input-field', () => {
  const factory = generateShallowRenderer(CMixedInputField, { stubs });

  it('Value changed to empty string after trigger the input with null value', () => {
    const wrapper = factory({
      propsData: {
        value: 'Value',
        inputType: PATTERN_FIELD_TYPES.string,
      },
    });

    selectTextField(wrapper).triggerCustomEvent('update:modelValue', null);

    expect(wrapper).toEmitInput('');
  });

  it('Value changed after trigger the input with number value', () => {
    const newNumber = Faker.datatype.number();
    const wrapper = factory({
      propsData: {
        value: 12,
        inputType: PATTERN_FIELD_TYPES.number,
      },
    });

    selectTextField(wrapper).triggerCustomEvent('update:modelValue', `${newNumber}`);

    expect(wrapper).toEmitInput(newNumber);
  });

  it('Value changed after trigger the select with items', () => {
    const item = {
      text: Faker.datatype.string(),
      value: Faker.datatype.string(),
    };
    const wrapper = factory({
      propsData: {
        value: '',
        items: [item],
        inputType: PATTERN_FIELD_TYPES.string,
      },
    });

    selectCombobox(wrapper).triggerCustomEvent('update:modelValue', item.value);

    expect(wrapper).toEmitInput(item.value);
  });
});
