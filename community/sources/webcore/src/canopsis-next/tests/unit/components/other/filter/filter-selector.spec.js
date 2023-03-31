import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { createSelectInputStub } from '@unit/stubs/input';

import FilterSelector from '@/components/other/filter/filter-selector.vue';

const localVue = createVueInstance();

const stubs = {
  'v-select': createSelectInputStub('v-select'),
  'c-enabled-field': true,
};

const snapshotStubs = {
  'c-enabled-field': true,
};

const factory = (options = {}) => shallowMount(FilterSelector, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(FilterSelector, {
  localVue,
  stubs: snapshotStubs,

  parentComponent: {
    provide: {
      listClick: jest.fn(),
    },
  },

  ...options,
});

const selectFilterSelectField = wrapper => wrapper.find('select.v-select');

describe('filter-selector', () => {
  const lockedFilters = [
    {
      _id: '1',
      title: 'Locked filter 1',
      filter: { _id: 1 },
    },
    {
      _id: '2',
      title: 'Locked filter 2',
      filter: { _id: 2 },
    },
  ];

  const filters = [
    {
      _id: '3',
      title: 'Filter 1',
      filter: { _id: 1 },
    },
    {
      _id: '4',
      title: 'Filter 2',
      filter: { _id: 2 },
    },
  ];

  it('Filter changed after trigger the select field', () => {
    const wrapper = factory({
      propsData: {
        filters,
        value: filters[0]._id,
      },
    });

    const filterSelectField = selectFilterSelectField(wrapper);

    const [, filter] = filters;

    filterSelectField.vm.$emit('input', filter._id);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(filter._id);
  });

  it('Renders `filter-selector` with default props', () => {
    const wrapper = snapshotFactory();

    const menuContent = wrapper.find('.v-menu__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `filter-selector` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        filters,
        lockedFilters,
        long: true,
        value: filters[0]._id,
        label: 'Custom label',
      },
    });

    const menuContent = wrapper.find('.v-menu__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `filter-selector` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        filters,
        lockedFilters,
        hideMultiply: true,
        value: filters[0]._id,
        clearable: false,
      },
    });

    const menuContent = wrapper.find('.v-menu__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `filter-selector` with array value', () => {
    const wrapper = snapshotFactory({
      propsData: {
        filters,
        lockedFilters,
        value: [filters[0]._id],
      },
    });

    const menuContent = wrapper.find('.v-menu__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `filter-selector` with badges', () => {
    const wrapper = snapshotFactory({
      propsData: {
        filters: filters.map(filter => ({ ...filter, old_mongo_query: true })),
        lockedFilters: lockedFilters.map(filter => ({ ...filter, old_mongo_query: true })),
      },
    });

    const menuContent = wrapper.find('.v-menu__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });
});
