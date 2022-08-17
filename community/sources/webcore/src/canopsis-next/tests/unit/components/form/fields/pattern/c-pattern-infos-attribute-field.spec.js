import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { createSelectInputStub } from '@unit/stubs/input';

import CPatternInfosAttributeField from '@/components/forms/fields/pattern/c-pattern-infos-attribute-field.vue';

const localVue = createVueInstance();

const stubs = {
  'v-select': createSelectInputStub('v-select'),
  'v-combobox': createSelectInputStub('v-combobox'),
};

const factory = (options = {}) => shallowMount(CPatternInfosAttributeField, {
  localVue,
  stubs,
  ...options,
});

const snapshotFactory = (options = {}) => mount(CPatternInfosAttributeField, {
  localVue,
  ...options,
});

const selectDictionarySelect = wrapper => wrapper.find('select.v-combobox');
const selectFieldSelect = wrapper => wrapper.find('select.v-select');

describe('c-pattern-infos-attribute-field', () => {
  it('Dictionary changed after trigger the dictionary select', () => {
    const value = {
      dictionary: 'test',
      field: 'name',
    };
    const wrapper = factory({
      propsData: {
        value,
      },
    });
    const dictionarySelect = selectDictionarySelect(wrapper);

    const newDictionary = 'newDictionary';

    dictionarySelect.vm.$emit('input', newDictionary);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual({
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
      },
    });
    const fieldSelect = selectFieldSelect(wrapper);

    const newField = 'newField';

    fieldSelect.vm.$emit('input', newField);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual({
      dictionary: value.dictionary,
      field: newField,
    });
  });

  it('Renders `c-pattern-infos-attribute-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          dictionary: '',
          field: '',
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `c-pattern-infos-attribute-field` with custom props', () => {
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
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
