import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createSelectInputStub } from '@unit/stubs/input';

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
    });
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
