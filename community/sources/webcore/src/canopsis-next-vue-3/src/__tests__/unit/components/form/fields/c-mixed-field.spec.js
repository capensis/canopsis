import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import { PATTERN_FIELD_TYPES } from '@/constants';

import CMixedField from '@/components/forms/fields/compat/c-mixed-field.vue';

const stubs = {
  'c-input-type-field': true,
  'c-mixed-input-field': true,
};

const selectInputTypeField = wrapper => wrapper.find('c-input-type-field-stub');

describe('c-mixed-field', () => {
  const factory = generateShallowRenderer(CMixedField, {
    stubs,
    parentComponent: {
      $_veeValidate: {
        validator: 'new',
      },
    },
  });

  it('Input type changed after trigger input type field with string value', () => {
    const wrapper = factory({
      propsData: {
        value: 12,
      },
    });

    const inputTypeField = selectInputTypeField(wrapper);

    inputTypeField.triggerCustomEvent('input', PATTERN_FIELD_TYPES.string);

    expect(wrapper).toEmitInput('12');
  });

  it('Input type changed after trigger input type field with number value', () => {
    const wrapper = factory({
      propsData: {
        value: '12',
      },
    });

    const inputTypeField = selectInputTypeField(wrapper);

    inputTypeField.triggerCustomEvent('input', PATTERN_FIELD_TYPES.number);

    expect(wrapper).toEmitInput(12);
  });

  it('Input type changed after trigger input type field with null value', () => {
    const wrapper = factory({
      propsData: {
        value: '12',
      },
    });

    const inputTypeField = selectInputTypeField(wrapper);

    inputTypeField.triggerCustomEvent('input', PATTERN_FIELD_TYPES.null);

    expect(wrapper).toEmitInput(null);
  });

  it('Input type changed after trigger input type field with boolean value', () => {
    const wrapper = factory({
      propsData: {
        value: 12,
      },
    });

    const inputTypeField = selectInputTypeField(wrapper);

    inputTypeField.triggerCustomEvent('input', PATTERN_FIELD_TYPES.boolean);

    expect(wrapper).toEmitInput(true);
  });

  it('Input type changed after trigger input type field with array value', () => {
    const wrapper = factory({
      propsData: {
        value: 12,
      },
    });

    const inputTypeField = selectInputTypeField(wrapper);

    inputTypeField.triggerCustomEvent('input', PATTERN_FIELD_TYPES.stringArray);

    expect(wrapper).toEmitInput([
      expect.objectContaining({
        key: expect.any(String),
        value: '12',
      }),
    ]);
  });

  it('Value changed on the first field after remove selected type', async () => {
    const value = Faker.datatype.number();
    const wrapper = factory({
      propsData: {
        value,
        types: [
          { value: PATTERN_FIELD_TYPES.string },
          { value: PATTERN_FIELD_TYPES.boolean },
          { value: PATTERN_FIELD_TYPES.number },
        ],
      },
    });

    await wrapper.setProps({
      types: [{ value: PATTERN_FIELD_TYPES.string }, { value: PATTERN_FIELD_TYPES.boolean }],
    });

    expect(wrapper).toEmitInput(`${value}`);
  });

  it('Value change on undefined after remove all types', async () => {
    const value = Faker.datatype.number();
    const wrapper = factory({
      propsData: {
        value,
        types: [
          { value: PATTERN_FIELD_TYPES.string },
          { value: PATTERN_FIELD_TYPES.boolean },
          { value: PATTERN_FIELD_TYPES.number },
        ],
      },
    });

    await wrapper.setProps({
      types: [],
    });

    expect(wrapper).toEmitInput(undefined);
  });
});
