import Faker from 'faker';
import { Validator } from 'vee-validate';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { createInputStub, createSelectInputStub } from '@unit/stubs/input';
import { FILTER_INPUT_TYPES } from '@/constants';

import CMixedField from '@/components/forms/fields/c-mixed-field.vue';

const localVue = createVueInstance();

const stubs = {
  'v-select': createSelectInputStub('v-select'),
  'v-text-field': createInputStub('v-text-field'),
  'v-combobox': createInputStub('v-combobox'),
  'c-array-mixed-field': true,
};

const snapshotStubs = {
  'c-array-mixed-field': true,
};

const factory = (options = {}) => shallowMount(CMixedField, {
  localVue,
  stubs,
  provide: {
    $validator: new Validator(),
  },
  ...options,
});

describe('c-mixed-field', () => {
  it('Value set into select', () => {
    const value = Faker.datatype.number();
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    const selectElement = wrapper.find('select.v-select');

    expect(selectElement.element.value).toBe(FILTER_INPUT_TYPES.number);
  });

  it('Value changed after change the type select to string', () => {
    const value = Faker.datatype.number();
    const wrapper = factory({
      propsData: {
        value,
      },
    });
    const selectElement = wrapper.find('select.v-select');

    selectElement.setValue(FILTER_INPUT_TYPES.string);

    const inputEvents = wrapper.emitted('input');
    expect(inputEvents).toHaveLength(1);

    const [inputEventData] = inputEvents[0];
    expect(inputEventData).toEqual(`${value}`);
  });

  it('Value changed after change the type select to number', () => {
    const number = Faker.datatype.number();
    const wrapper = factory({
      propsData: {
        value: `${number}`,
      },
    });
    const selectElement = wrapper.find('select.v-select');

    selectElement.setValue(FILTER_INPUT_TYPES.number);

    const inputEvents = wrapper.emitted('input');
    expect(inputEvents).toHaveLength(1);

    const [inputEventData] = inputEvents[0];
    expect(inputEventData).toEqual(number);
  });

  it('Value changed after change the type select to boolean', () => {
    const wrapper = factory({
      propsData: {
        value: Faker.datatype.string(),
      },
    });
    const selectElement = wrapper.find('select.v-select');

    selectElement.setValue(FILTER_INPUT_TYPES.boolean);

    const inputEvents = wrapper.emitted('input');
    expect(inputEvents).toHaveLength(1);

    const [inputEventData] = inputEvents[0];
    expect(inputEventData).toEqual(true);
  });

  it('Value changed after change the type select to null', () => {
    const wrapper = factory({
      propsData: {
        value: Faker.datatype.string(),
      },
    });
    const selectElement = wrapper.find('select.v-select');

    selectElement.setValue(FILTER_INPUT_TYPES.null);

    const inputEvents = wrapper.emitted('input');
    expect(inputEvents).toHaveLength(1);

    const [inputEventData] = inputEvents[0];
    expect(inputEventData).toEqual(null);
  });

  it('Value changed after change the type select to array', () => {
    const value = Faker.datatype.string();
    const wrapper = factory({
      propsData: {
        value,
      },
    });
    const selectElement = wrapper.find('select.v-select');

    selectElement.setValue(FILTER_INPUT_TYPES.array);

    const inputEvents = wrapper.emitted('input');
    expect(inputEvents).toHaveLength(1);

    const [inputEventData] = inputEvents[0];
    expect(inputEventData).toEqual([value]);
  });

  it('Value changed on the first field after remove selected type', async () => {
    const value = Faker.datatype.number();
    const wrapper = factory({
      propsData: {
        value,
        types: [
          { value: FILTER_INPUT_TYPES.string },
          { value: FILTER_INPUT_TYPES.boolean },
          { value: FILTER_INPUT_TYPES.number },
        ],
      },
    });

    await wrapper.setProps({
      types: [
        { value: FILTER_INPUT_TYPES.string },
        { value: FILTER_INPUT_TYPES.boolean },
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
          { value: FILTER_INPUT_TYPES.string },
          { value: FILTER_INPUT_TYPES.boolean },
          { value: FILTER_INPUT_TYPES.number },
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

  it('Value changed to empty string after trigger the input with null value', () => {
    const wrapper = factory({
      propsData: {
        value: 'Value',
      },
    });
    const inputElement = wrapper.find('input.v-text-field');

    inputElement.vm.$emit('input', null);

    const inputEvents = wrapper.emitted('input');
    expect(inputEvents).toHaveLength(1);

    const [inputEventData] = inputEvents[0];
    expect(inputEventData).toEqual('');
  });

  it('Value changed after trigger the input with number value', () => {
    const newNumber = Faker.datatype.number();
    const wrapper = factory({
      propsData: {
        value: 12,
      },
    });
    const inputElement = wrapper.find('input.v-text-field');

    inputElement.setValue(`${newNumber}`);

    const inputEvents = wrapper.emitted('input');
    expect(inputEvents).toHaveLength(1);

    const [inputEventData] = inputEvents[0];
    expect(inputEventData).toEqual(newNumber);
  });

  it('Value changed after trigger the mixed field with array value', () => {
    const newValue = Faker.datatype.array();
    const wrapper = factory({
      propsData: {
        value: [],
      },
    });
    const arrayMixedFieldElement = wrapper.find('c-array-mixed-field-stub');

    arrayMixedFieldElement.vm.$emit('change', newValue);

    const inputEvents = wrapper.emitted('input');
    expect(inputEvents).toHaveLength(1);

    const [inputEventData] = inputEvents[0];
    expect(inputEventData).toEqual(newValue);
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
      },
    });
    const comboboxInputElement = wrapper.find('input.v-combobox');

    comboboxInputElement.setValue(item.value);

    const inputEvents = wrapper.emitted('input');
    expect(inputEvents).toHaveLength(1);

    const [inputEventData] = inputEvents[0];
    expect(inputEventData).toEqual(item.value);
  });

  it('v-validate works correctly with component', async () => {
    const name = Faker.datatype.string();
    const value = Faker.datatype.string();
    const validator = new Validator();

    mount({
      inject: ['$validator'],
      components: {
        CMixedField,
      },
      props: ['name', 'value'],
      template: `
        <c-mixed-field v-validate="'required'" :name="name" :value="value" />
      `,
    }, {
      localVue,
      stubs,
      mocks: { $t: () => {} },
      provide: {
        $validator: validator,
      },
      propsData: {
        value,
        name,
      },
    });

    await validator.validateAll();

    expect(validator.fields.find({ name })).toBeTruthy();
  });

  it('Renders `c-mixed-field` with default props correctly', () => {
    const wrapper = mount(CMixedField, {
      localVue,
      stubs: snapshotStubs,
      provide: {
        $validator: new Validator(),
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `c-mixed-field` with errors correctly', () => {
    const validator = new Validator();

    const wrapper = mount(CMixedField, {
      localVue,
      provide: {
        $validator: validator,
      },
      stubs: snapshotStubs,
      propsData: {
        errorMessages: ['First error message', 'Second error message'],
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `c-mixed-field` with custom props correctly', () => {
    const validator = new Validator();

    const wrapper = mount(CMixedField, {
      localVue,
      provide: {
        $validator: validator,
      },
      stubs: snapshotStubs,
      propsData: {
        value: 'Value',
        name: 'mixedFieldName',
        label: 'Mixed field label',
        disabled: true,
        soloInverted: true,
        flat: true,
        hideDetails: true,
        errorMessages: ['First error message', 'Second error message'],
        types: [
          { value: FILTER_INPUT_TYPES.string, text: 'Custom string' },
          { value: FILTER_INPUT_TYPES.number },
          { value: FILTER_INPUT_TYPES.boolean },
        ],
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `c-mixed-field` with string type and items correctly', () => {
    const validator = new Validator();

    const wrapper = mount(CMixedField, {
      localVue,
      provide: {
        $validator: validator,
      },
      stubs: snapshotStubs,
      propsData: {
        value: false,
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

    expect(wrapper.element).toMatchSnapshot();
    menuContents.wrappers.forEach((menuContent) => {
      expect(menuContent.element).toMatchSnapshot();
    });
  });

  it('Renders `c-mixed-field` with boolean type correctly', () => {
    const validator = new Validator();

    const wrapper = mount(CMixedField, {
      localVue,
      provide: {
        $validator: validator,
      },
      stubs: snapshotStubs,
      propsData: {
        value: false,
        name: 'mixedFieldName',
        label: 'Mixed field with boolean type',
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `c-mixed-field` with number type correctly', () => {
    const validator = new Validator();

    const wrapper = mount(CMixedField, {
      localVue,
      provide: {
        $validator: validator,
      },
      stubs: snapshotStubs,
      propsData: {
        value: 222,
        name: 'mixedFieldName',
        label: 'Mixed field with number type',
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `c-mixed-field` with null type correctly', () => {
    const validator = new Validator();

    const wrapper = mount(CMixedField, {
      localVue,
      provide: {
        $validator: validator,
      },
      stubs: snapshotStubs,
      propsData: {
        value: null,
        name: 'mixedFieldName',
        label: 'Mixed field with null type',
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `c-mixed-field` with null type correctly', () => {
    const validator = new Validator();

    const wrapper = mount(CMixedField, {
      localVue,
      provide: {
        $validator: validator,
      },
      stubs: snapshotStubs,
      propsData: {
        value: [0, '1', null, false, []],
        name: 'mixedFieldName',
        label: 'Mixed field with null type',
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });
});
