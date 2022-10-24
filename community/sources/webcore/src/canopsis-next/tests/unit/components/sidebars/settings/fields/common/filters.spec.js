import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import Filters from '@/components/sidebars/settings/fields/common/filters.vue';

const localVue = createVueInstance();

const stubs = {
  'widget-settings-item': true,
  'filter-selector': true,
  'filters-list': true,
};

const factory = (options = {}) => shallowMount(Filters, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(Filters, {
  localVue,
  stubs,

  ...options,
});

const selectFilterSelectorField = wrapper => wrapper.find('filter-selector-stub');

describe('filters', () => {
  const filters = [{
    _id: 'filter-id',
    filter: { $and: [{ test: { test: 'test' } }] },
    title: 'Filter',
  }];

  it('Filter updated after trigger input on the filter selector field', () => {
    const wrapper = factory();

    const filterSelectField = selectFilterSelectorField(wrapper);

    filterSelectField.vm.$emit('input', filters[0]);

    const updateConditionEvents = wrapper.emitted('input');

    expect(updateConditionEvents).toHaveLength(1);

    const [eventData] = updateConditionEvents[0];
    expect(eventData).toBe(filters[0]);
  });
  /* TODO: fix it
  it('Filters updated after trigger update:filters on the filter selector field', () => {
    const wrapper = factory();

    const filterSelectField = selectFilterSelectorField(wrapper);

    filterSelectField.vm.$emit('update:filters', filters);

    const updateConditionEvents = wrapper.emitted('update:filters');

    expect(updateConditionEvents).toHaveLength(1);

    const [eventData] = updateConditionEvents[0];
    expect(eventData).toBe(filters);
  }); */

  it('Renders `filters` with default and required props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `filters` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        filters,
        widgetId: Faker.datatype.string(),
        value: filters[0]._id,
        addable: false,
        editable: false,
        withAlarm: false,
        withEntity: false,
        withPbehavior: false,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
