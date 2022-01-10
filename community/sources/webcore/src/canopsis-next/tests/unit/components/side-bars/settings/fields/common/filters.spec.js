import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { ENTITIES_TYPES, FILTER_MONGO_OPERATORS } from '@/constants';

import Filters from '@/components/side-bars/settings/fields/common/filters.vue';

const localVue = createVueInstance();

const stubs = {
  'filter-selector': true,
};

const factory = (options = {}) => shallowMount(Filters, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(Filters, {
  localVue,
  stubs,

  parentComponent: {
    provide: {
      list: {
        register: jest.fn(),
        unregister: jest.fn(),
      },
      listClick: jest.fn(),
    },
  },

  ...options,
});

const selectFilterSelectorField = wrapper => wrapper.find('filter-selector-stub');

describe('filters', () => {
  const filters = [{
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

  it('Conditions updated after trigger update:condition on the filter selector field', () => {
    const wrapper = factory();

    const filterSelectField = selectFilterSelectorField(wrapper);

    filterSelectField.vm.$emit('update:condition', FILTER_MONGO_OPERATORS.or);

    const updateConditionEvents = wrapper.emitted('update:condition');

    expect(updateConditionEvents).toHaveLength(1);

    const [eventData] = updateConditionEvents[0];
    expect(eventData).toBe(FILTER_MONGO_OPERATORS.or);
  });

  it('Filters updated after trigger update:filters on the filter selector field', () => {
    const wrapper = factory();

    const filterSelectField = selectFilterSelectorField(wrapper);

    filterSelectField.vm.$emit('update:filters', filters);

    const updateConditionEvents = wrapper.emitted('update:filters');

    expect(updateConditionEvents).toHaveLength(1);

    const [eventData] = updateConditionEvents[0];
    expect(eventData).toBe(filters);
  });

  it('Renders `filters` with default and required props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `filters` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        filters,
        value: {
          filter: { $and: [{ test: { test: 'test' } }] },
          title: 'Filter',
        },
        condition: FILTER_MONGO_OPERATORS.and,
        hideSelect: true,
        addable: false,
        editable: false,
        entitiesType: ENTITIES_TYPES.entity,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
