import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createSelectInputStub } from '@unit/stubs/input';

import { PATTERN_FIELD_TYPES, PATTERN_RULE_INFOS_FIELDS } from '@/constants';

import CInfosAttributeField from '@/components/forms/fields/c-infos-attribute-field.vue';

const stubs = {
  'v-select': createSelectInputStub('v-select'),
  'v-combobox': createSelectInputStub('v-combobox'),
};

const selectDictionarySelect = wrapper => wrapper.find('select.v-combobox');
const selectFieldSelect = wrapper => wrapper.find('select.v-select');

describe('c-infos-attribute-field', () => {
  const factory = generateShallowRenderer(CInfosAttributeField, { stubs });
  const snapshotFactory = generateRenderer(CInfosAttributeField);

  it('Dictionary changed after trigger the dictionary select', () => {
    const value = {
      dictionary: 'test',
      field: 'name',
    };
    const wrapper = factory({
      propsData: {
        value,
        combobox: true,
      },
    });
    const dictionarySelect = selectDictionarySelect(wrapper);

    const newDictionary = 'newDictionary';

    dictionarySelect.triggerCustomEvent('input', newDictionary);

    expect(wrapper).toEmitInput({
      dictionary: newDictionary,
      field: value.field,
      value: '',
      fieldType: PATTERN_FIELD_TYPES.string,
    });
  });

  it('Dictionary with object changed after trigger the dictionary select in combobox mode', () => {
    const value = {
      dictionary: 'test',
      field: 'name',
      value: 'test-value',
      fieldType: PATTERN_FIELD_TYPES.number,
    };
    const items = [
      {
        value: 'new-dictionary',
        definedType: PATTERN_FIELD_TYPES.boolean,
      },
    ];
    const wrapper = factory({
      propsData: {
        value,
        items,
        combobox: true,
      },
    });
    const dictionarySelect = selectDictionarySelect(wrapper);

    dictionarySelect.triggerCustomEvent('input', items[0]);

    expect(wrapper).toEmitInput({
      dictionary: items[0].value,
      field: value.field,
      value: true,
      fieldType: items[0].definedType,
    });
  });

  it('Dictionary changed in text field mode', () => {
    const value = {
      dictionary: 'test',
      field: 'name',
      value: 'test-value',
    };
    const wrapper = factory({
      propsData: {
        value,
        combobox: false,
      },
    });

    const newDictionary = Faker.datatype.string();

    wrapper.vm.updateInfosDictionary(newDictionary);

    expect(wrapper).toEmitInput({
      dictionary: newDictionary,
      field: value.field,
      value: 'test-value',
      fieldType: PATTERN_FIELD_TYPES.string,
    });
  });

  it('Field changed after trigger the field select', () => {
    const value = {
      dictionary: 'test',
      field: 'name',
    };
    const wrapper = factory({
      propsData: {
        value,
        combobox: true,
      },
    });
    const fieldSelect = selectFieldSelect(wrapper);

    const newField = 'newField';

    fieldSelect.triggerCustomEvent('input', newField);

    expect(wrapper).toEmitInput({
      dictionary: value.dictionary,
      field: newField,
      value: '',
      fieldType: PATTERN_FIELD_TYPES.string,
    });
  });

  it('Field changed to value with definedType from items', () => {
    const value = {
      dictionary: 'test-dictionary',
      field: 'name',
      value: 'test-value',
      fieldType: PATTERN_FIELD_TYPES.string,
    };
    const items = [
      {
        value: 'test-dictionary',
        definedType: PATTERN_FIELD_TYPES.number,
      },
    ];
    const wrapper = factory({
      propsData: {
        value,
        items,
        combobox: true,
      },
    });
    const fieldSelect = selectFieldSelect(wrapper);

    fieldSelect.triggerCustomEvent('input', PATTERN_RULE_INFOS_FIELDS.value);

    expect(wrapper).toEmitInput({
      dictionary: value.dictionary,
      field: PATTERN_RULE_INFOS_FIELDS.value,
      value: 0,
      fieldType: PATTERN_FIELD_TYPES.number,
    });
  });

  it('Field changed to value without definedType defaults to string', () => {
    const value = {
      dictionary: 'test-dictionary',
      field: 'name',
      value: 123,
      fieldType: PATTERN_FIELD_TYPES.number,
    };
    const items = [
      {
        value: 'test-dictionary',
      },
    ];
    const wrapper = factory({
      propsData: {
        value,
        items,
        combobox: true,
      },
    });
    const fieldSelect = selectFieldSelect(wrapper);

    fieldSelect.triggerCustomEvent('input', PATTERN_RULE_INFOS_FIELDS.value);

    expect(wrapper).toEmitInput({
      dictionary: value.dictionary,
      field: PATTERN_RULE_INFOS_FIELDS.value,
      value: '123',
      fieldType: PATTERN_FIELD_TYPES.string,
    });
  });

  it('Field changed to name always sets string fieldType', () => {
    const value = {
      dictionary: 'test-dictionary',
      field: PATTERN_RULE_INFOS_FIELDS.value,
      value: true,
      fieldType: PATTERN_FIELD_TYPES.boolean,
    };
    const items = [
      {
        value: 'test-dictionary',
        definedType: PATTERN_FIELD_TYPES.boolean,
      },
    ];
    const wrapper = factory({
      propsData: {
        value,
        items,
        combobox: true,
      },
    });
    const fieldSelect = selectFieldSelect(wrapper);

    fieldSelect.triggerCustomEvent('input', PATTERN_RULE_INFOS_FIELDS.name);

    expect(wrapper).toEmitInput({
      dictionary: value.dictionary,
      field: PATTERN_RULE_INFOS_FIELDS.name,
      value: 'true',
      fieldType: PATTERN_FIELD_TYPES.string,
    });
  });

  it('Field select is disabled when dictionary is empty', () => {
    const value = {
      dictionary: '',
      field: 'name',
    };
    const wrapper = factory({
      propsData: {
        value,
      },
    });
    const fieldSelect = selectFieldSelect(wrapper);

    expect(fieldSelect.attributes('disabled')).toBe('disabled');
  });

  it('Field select is enabled when dictionary is provided', () => {
    const value = {
      dictionary: 'test-dictionary',
      field: 'name',
    };
    const wrapper = factory({
      propsData: {
        value,
      },
    });
    const fieldSelect = selectFieldSelect(wrapper);

    expect(fieldSelect.attributes('disabled')).toBeFalsy();
  });

  it('Shows pending loading state in combobox', () => {
    const value = {
      dictionary: 'test',
      field: 'name',
    };
    const wrapper = factory({
      propsData: {
        value,
        combobox: true,
        pending: true,
      },
    });
    const dictionarySelect = selectDictionarySelect(wrapper);

    expect(dictionarySelect.attributes('loading')).toBe('true');
  });

  it('Uses row layout when row prop is true', () => {
    const value = {
      dictionary: 'test',
      field: 'name',
    };
    const wrapper = factory({
      propsData: {
        value,
        row: true,
      },
    });

    expect(wrapper.vm.row).toBe(true);
    expect(wrapper.html()).toContain('xs6');
    expect(wrapper.html()).toContain('pl-3');
  });

  it('Uses column layout when column prop is true', () => {
    const value = {
      dictionary: 'test',
      field: 'name',
    };
    const wrapper = factory({
      propsData: {
        value,
        column: true,
      },
    });

    expect(wrapper.vm.column).toBe(true);
    expect(wrapper.html()).toContain('column');
  });

  it('Renders `c-pattern-infos-attribute-field` with default props and combobox', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          dictionary: '',
          field: '',
        },
        combobox: true,
      },
    });

    await wrapper.activateAllMenus();

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `c-pattern-infos-attribute-field` with default props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          dictionary: '',
          field: '',
        },
      },
    });

    await wrapper.activateAllMenus();

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `c-pattern-infos-attribute-field` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          dictionary: 'Test text',
          field: 'test',
        },
        label: 'Custom label',
        items: [
          {
            value: 'Test text',
          },
        ],
        name: 'custom_filter_infos_attribute_name',
        disabled: true,
        combobox: true,
      },
    });

    await wrapper.activateAllMenus();

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `c-pattern-infos-attribute-field` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          dictionary: 'Test text',
          field: 'test',
        },
        label: 'Custom label',
        items: [
          {
            value: 'Test text',
          },
        ],
        name: 'custom_filter_infos_attribute_name',
        disabled: true,
        combobox: true,
      },
    });

    await wrapper.activateAllMenus();

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
