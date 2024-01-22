import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { PATTERN_FIELD_TYPES } from '@/constants';

import CMixedField from '@/components/forms/fields/c-mixed-field.vue';

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
  const snapshotFactory = generateRenderer(CMixedField, { stubs });

  it('Input type changed after trigger input type field with string value', () => {
    const wrapper = factory({
      propsData: {
        value: 12,
      },
    });

    const inputTypeField = selectInputTypeField(wrapper);

    inputTypeField.triggerCustomEvent('input', PATTERN_FIELD_TYPES.string);

    expect(wrapper).toEmit('input', '12');
  });

  it('Input type changed after trigger input type field with number value', () => {
    const wrapper = factory({
      propsData: {
        value: '12',
      },
    });

    const inputTypeField = selectInputTypeField(wrapper);

    inputTypeField.triggerCustomEvent('input', PATTERN_FIELD_TYPES.number);

    expect(wrapper).toEmit('input', 12);
  });

  it('Input type changed after trigger input type field with null value', () => {
    const wrapper = factory({
      propsData: {
        value: '12',
      },
    });

    const inputTypeField = selectInputTypeField(wrapper);

    inputTypeField.triggerCustomEvent('input', PATTERN_FIELD_TYPES.null);

    expect(wrapper).toEmit('input', null);
  });

  it('Input type changed after trigger input type field with boolean value', () => {
    const wrapper = factory({
      propsData: {
        value: 12,
      },
    });

    const inputTypeField = selectInputTypeField(wrapper);

    inputTypeField.triggerCustomEvent('input', PATTERN_FIELD_TYPES.boolean);

    expect(wrapper).toEmit('input', true);
  });

  it('Input type changed after trigger input type field with array value', () => {
    const wrapper = factory({
      propsData: {
        value: 12,
      },
    });

    const inputTypeField = selectInputTypeField(wrapper);

    inputTypeField.triggerCustomEvent('input', PATTERN_FIELD_TYPES.stringArray);

    expect(wrapper).toEmit('input', ['12']);
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
      types: [
        { value: PATTERN_FIELD_TYPES.string },
        { value: PATTERN_FIELD_TYPES.boolean },
      ],
    });

    expect(wrapper).toEmit('input', `${value}`);
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

    expect(wrapper).toEmit('input', undefined);
  });

  it('Renders `c-mixed-field` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-mixed-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 123,
        name: 'customName',
        label: 'Custom label',
        disabled: true,
        flat: true,
        hideDetails: true,
        errorMessages: ['Message'],
        items: [{ value2: 'value', text2: 'Value' }],
        itemText: 'text2',
        itemValue: 'value2',
        types: [
          { value: PATTERN_FIELD_TYPES.string },
          { value: PATTERN_FIELD_TYPES.number },
          { value: PATTERN_FIELD_TYPES.stringArray },
        ],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
