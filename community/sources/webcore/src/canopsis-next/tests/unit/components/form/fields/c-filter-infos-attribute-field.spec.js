import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { createSelectInputStub } from '@unit/stubs/input';

import CFilterInfosAttributeField from '@/components/forms/fields/c-filter-infos-attribute-field.vue';

const localVue = createVueInstance();

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const factory = (options = {}) => shallowMount(CFilterInfosAttributeField, {
  localVue,
  stubs,
  ...options,
});

const snapshotFactory = (options = {}) => mount(CFilterInfosAttributeField, {
  localVue,
  ...options,
});

const selectSelectFields = wrapper => wrapper.findAll('select.v-select');
const selectDictionarySelect = wrapper => selectSelectFields(wrapper).at(0);
const selectFieldSelect = wrapper => selectSelectFields(wrapper).at(1);

describe('c-filter-infos-attribute-field', () => {
  it('Dictionary changed after trigger the dictionary select', () => {
    const wrapper = factory({
      propsData: {
        value: 'infos.test.name',
      },
    });
    const dictionarySelect = selectDictionarySelect(wrapper);

    const newDictionary = 'newDictionary';

    dictionarySelect.vm.$emit('input', newDictionary);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(`infos.${newDictionary}.name`);
  });

  it('Field changed after trigger the field select', () => {
    const wrapper = factory({
      propsData: {
        value: 'infos.test.name',
      },
    });
    const fieldSelect = selectFieldSelect(wrapper);

    const newField = 'newField';

    fieldSelect.vm.$emit('input', newField);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(`infos.test.${newField}`);
  });

  it('Renders `c-filter-infos-attribute-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 'infos.',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `c-filter-infos-attribute-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 'infos.test.test',
        label: 'Custom label',
        items: [
          {
            value: 'test',
            text: 'Test text',
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
