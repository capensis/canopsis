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

    inputTypeField.vm.$emit('input', PATTERN_FIELD_TYPES.string);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe('12');
  });

  it('Input type changed after trigger input type field with number value', () => {
    const wrapper = factory({
      propsData: {
        value: '12',
      },
    });

    const inputTypeField = selectInputTypeField(wrapper);

    inputTypeField.vm.$emit('input', PATTERN_FIELD_TYPES.number);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(12);
  });

  it('Input type changed after trigger input type field with null value', () => {
    const wrapper = factory({
      propsData: {
        value: '12',
      },
    });

    const inputTypeField = selectInputTypeField(wrapper);

    inputTypeField.vm.$emit('input', PATTERN_FIELD_TYPES.null);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(null);
  });

  it('Input type changed after trigger input type field with boolean value', () => {
    const wrapper = factory({
      propsData: {
        value: 12,
      },
    });

    const inputTypeField = selectInputTypeField(wrapper);

    inputTypeField.vm.$emit('input', PATTERN_FIELD_TYPES.boolean);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(true);
  });

  it('Input type changed after trigger input type field with array value', () => {
    const wrapper = factory({
      propsData: {
        value: 12,
      },
    });

    const inputTypeField = selectInputTypeField(wrapper);

    inputTypeField.vm.$emit('input', PATTERN_FIELD_TYPES.stringArray);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual(['12']);
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

    const inputEvents = wrapper.emitted('input');
    expect(inputEvents).toHaveLength(1);

    const [inputEventData] = inputEvents[0];
    expect(inputEventData).toEqual(`${value}`);
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

    const inputEvents = wrapper.emitted('input');
    expect(inputEvents).toHaveLength(1);

    const [inputEventData] = inputEvents[0];
    expect(inputEventData).toEqual(undefined);
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
