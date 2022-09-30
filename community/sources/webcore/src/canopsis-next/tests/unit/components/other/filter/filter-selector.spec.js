import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { createSelectInputStub } from '@unit/stubs/input';
import { FILTER_MONGO_OPERATORS } from '@/constants';

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
      filter: { [FILTER_MONGO_OPERATORS.and]: { _id: 1 } },
    },
    {
      _id: '2',
      title: 'Locked filter 2',
      filter: { [FILTER_MONGO_OPERATORS.and]: { _id: 2 } },
    },
  ];
  const filters = [
    {
      _id: '3',
      title: 'Filter 1',
      filter: { [FILTER_MONGO_OPERATORS.or]: { _id: 1 } },
    },
    {
      _id: '4',
      title: 'Filter 2',
      filter: { [FILTER_MONGO_OPERATORS.or]: { _id: 2 } },
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
        hideIcon: true,
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
        hideIcon: false,
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
});
