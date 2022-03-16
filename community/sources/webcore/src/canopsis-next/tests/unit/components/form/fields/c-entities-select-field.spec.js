import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { createSelectInputStub } from '@unit/stubs/input';
import { SAMPLINGS } from '@/constants';

import CEntitiesSelectField from '@/components/forms/fields/entity/c-entities-select-field.vue';

const localVue = createVueInstance();

const stubs = {
  'v-autocomplete': createSelectInputStub('v-autocomplete'),
};

const factory = (options = {}) => shallowMount(CEntitiesSelectField, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(CEntitiesSelectField, {
  localVue,

  ...options,
});

const selectAutocomplete = wrapper => wrapper.find('select.v-autocomplete');

describe('c-entities-select-field', () => {
  const items = [
    {
      value: 'value',
      text: 'Text',
    },
    {
      value: 'value 2',
      text: 'Text 2',
    },
    {
      value: 'value 3',
      text: 'Text 3',
    },
  ];

  it('Value changed after trigger the input', () => {
    const wrapper = factory({
      propsData: {
        items,
        itemText: 'text',
        itemValue: 'value',
      },
    });
    const autocompleteElement = selectAutocomplete(wrapper);

    autocompleteElement.setValue(items[0].value);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(items[0].value);
  });

  it('Renders `c-entities-select-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: SAMPLINGS.day,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `c-entities-select-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: items[2].text,
        search: items[1].text,
        items,
        label: 'Custom label',
        name: 'customName',
        itemText: 'text',
        itemValue: 'value',
        disabled: true,
        loading: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `c-entities-select-field` with array value', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: items.map(({ text }) => text),
        items,
        itemText: 'text',
        itemValue: 'value',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
