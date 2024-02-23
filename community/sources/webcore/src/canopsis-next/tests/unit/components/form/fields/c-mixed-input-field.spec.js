import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createInputStub } from '@unit/stubs/input';

import { PATTERN_FIELD_TYPES } from '@/constants';

import CMixedInputField from '@/components/forms/fields/c-mixed-input-field.vue';

const stubs = {
  'v-text-field': createInputStub('v-text-field'),
  'v-combobox': createInputStub('v-combobox'),
  'c-array-text-field': true,
};

const snapshotStubs = {
  'c-array-text-field': true,
};

const selectTextField = wrapper => wrapper.find('input.v-text-field');
const selectCombobox = wrapper => wrapper.find('input.v-combobox');

describe('c-mixed-input-field', () => {
  const factory = generateShallowRenderer(CMixedInputField, {
    stubs,
    parentComponent: {
      $_veeValidate: {
        validator: 'new',
      },
    },
  });
  const snapshotFactory = generateRenderer(CMixedInputField, {
    stubs: snapshotStubs,
    parentComponent: {
      $_veeValidate: {
        validator: 'new',
      },
    },
  });

  it('Value changed to empty string after trigger the input with null value', () => {
    const wrapper = factory({
      propsData: {
        value: 'Value',
        inputType: PATTERN_FIELD_TYPES.string,
      },
    });

    selectTextField(wrapper).setValue(null);

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
    selectTextField(wrapper).setValue(`${newNumber}`);

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
    const combobox = selectCombobox(wrapper);

    combobox.setValue(item.value);

    expect(wrapper).toEmitInput(item.value);
  });

  it('Renders `c-mixed-input-field` with default props correctly', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `c-mixed-input-field` with errors correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        errorMessages: ['First error message', 'Second error message'],
      },
    });

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `c-mixed-input-field` with custom props correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 'Value',
        inputType: PATTERN_FIELD_TYPES.string,
        name: 'mixedFieldName',
        label: 'Mixed field label',
        disabled: true,
        flat: true,
        hideDetails: true,
        errorMessages: ['First error message', 'Second error message'],
        types: [
          { value: PATTERN_FIELD_TYPES.string, text: 'Custom string' },
          { value: PATTERN_FIELD_TYPES.number },
          { value: PATTERN_FIELD_TYPES.boolean },
        ],
      },
    });

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `c-mixed-input-field` with string type and items correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: false,
        inputType: PATTERN_FIELD_TYPES.boolean,
        name: 'mixedFieldName',
        label: 'Mixed field with boolean type',
        items: [{
          customText: 'Custom item text',
          customValue: 'Custom item value',
        }],
        itemText: 'customText',
        itemValue: 'customValue',
      },
    });

    const menuContents = wrapper.findAllMenus();

    expect(wrapper).toMatchSnapshot();
    menuContents.wrappers.forEach((menuContent) => {
      expect(menuContent.element).toMatchSnapshot();
    });
  });

  it('Renders `c-mixed-input-field` with boolean type correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: false,
        inputType: PATTERN_FIELD_TYPES.boolean,
        name: 'mixedFieldName',
        label: 'Mixed field with boolean type',
      },
    });

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `c-mixed-input-field` with number type correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 222,
        inputType: PATTERN_FIELD_TYPES.number,
        name: 'mixedFieldName',
        label: 'Mixed field with number type',
      },
    });

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `c-mixed-input-field` with null type correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        inputType: PATTERN_FIELD_TYPES.null,
        value: null,
        name: 'mixedFieldName',
        label: 'Mixed field with null type',
      },
    });

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `c-mixed-input-field` with null type correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: [0, '1', null, false, []],
        inputType: PATTERN_FIELD_TYPES.stringArray,
        name: 'mixedFieldName',
        label: 'Mixed field with null type',
      },
    });

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
