import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { createSelectInputStub } from '@unit/stubs/input';
import { SAMPLINGS } from '@/constants';

import CEntityField from '@/components/forms/fields/entity/c-entity-field.vue';
import CSelectField from '@/components/forms/fields/c-select-field';

const localVue = createVueInstance();

const stubs = {
  'c-select-field': createSelectInputStub('c-select-field'),
};

const snapshotStubs = {
  'c-select-field': CSelectField,
};

const factory = (options = {}) => shallowMount(CEntityField, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(CEntityField, {
  localVue,
  stubs: snapshotStubs,

  ...options,
});

const selectAutocomplete = wrapper => wrapper.find('.c-select-field');

describe('c-entity-field', () => {
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

  it('Renders `c-entity-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: SAMPLINGS.day,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `c-entity-field` with custom props', () => {
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

  it('Renders `c-entity-field` with array value', () => {
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
