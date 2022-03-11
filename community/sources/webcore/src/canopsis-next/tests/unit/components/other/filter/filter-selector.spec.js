import { omit } from 'lodash';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { mockModals } from '@unit/utils/mock-hooks';
import { createSelectInputStub } from '@unit/stubs/input';
import { ENTITIES_TYPES, FILTER_MONGO_OPERATORS, MODALS } from '@/constants';

import FilterSelector from '@/components/other/filter/filter-selector.vue';

const localVue = createVueInstance();

const stubs = {
  'v-select': createSelectInputStub('v-select'),
  'c-operator-field': true,
  'c-enabled-field': true,
  'c-action-btn': true,
};

const snapshotStubs = {
  'c-operator-field': true,
  'c-enabled-field': true,
  'c-action-btn': true,
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
const selectListButton = wrapper => wrapper.find('c-action-btn-stub');
const selectMixFiltersField = wrapper => wrapper.find('c-enabled-field-stub');
const selectConditionField = wrapper => wrapper.find('c-operator-field-stub');

describe('filter-selector', () => {
  const $modals = mockModals();
  const lockedFilters = [
    {
      title: 'Locked filter 1',
      filter: { [FILTER_MONGO_OPERATORS.and]: { _id: 1 } },
    },
    {
      title: 'Locked filter 2',
      filter: { [FILTER_MONGO_OPERATORS.and]: { _id: 2 } },
    },
  ];
  const filters = [
    {
      title: 'Filter 1',
      filter: { [FILTER_MONGO_OPERATORS.or]: { _id: 1 } },
    },
    {
      title: 'Filter 2',
      filter: { [FILTER_MONGO_OPERATORS.or]: { _id: 2 } },
    },
  ];

  it('Filter changed after trigger the select field', () => {
    const wrapper = factory({
      propsData: {
        filters,
        value: filters[0],
      },
    });

    const filterSelectField = selectFilterSelectField(wrapper);

    const [, filter] = filters;

    filterSelectField.vm.$emit('input', filter);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(filter);
  });

  it('Filter list modal opened after click on the button', () => {
    const wrapper = factory({
      propsData: {
        filters,
        value: filters[0],
        hasAccessToUserFilter: true,
        entitiesType: ENTITIES_TYPES.entity,
      },
      mocks: {
        $modals,
      },
    });

    const listButton = selectListButton(wrapper);

    listButton.vm.$emit('click');

    const filtersWithSelected = filters.map((filter, index) => ({
      ...filter,
      selected: index === 0,
    }));

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.filtersList,
        config: {
          filters: filtersWithSelected,
          hasAccessToAddFilter: true,
          hasAccessToEditFilter: true,
          entitiesType: ENTITIES_TYPES.entity,
          action: expect.any(Function),
        },
      },
    );

    const [modalArguments] = $modals.show.mock.calls[0];

    const actionValue = [
      ...filtersWithSelected,
      {
        title: 'Filter 3',
        selected: false,
        filter: {},
      },
    ];

    modalArguments.config.action(actionValue);

    const updateFiltersEvents = wrapper.emitted('update:filters');

    expect(updateFiltersEvents).toHaveLength(1);

    const [changedFilters, newFilter] = updateFiltersEvents[0];

    expect(changedFilters).toEqual(
      actionValue.map(filter => omit(filter, ['selected'])),
    );
    expect(newFilter).toEqual(filters[0]);
  });

  it('Filter list modal opened after click on the button with multiple mode', () => {
    const wrapper = factory({
      propsData: {
        filters,
        value: filters,
        hasAccessToUserFilter: true,
        entitiesType: ENTITIES_TYPES.entity,
      },
      mocks: {
        $modals,
      },
    });

    const listButton = selectListButton(wrapper);

    listButton.vm.$emit('click');

    const filtersWithSelected = filters.map(filter => ({
      ...filter,
      selected: true,
    }));

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.filtersList,
        config: {
          filters: filtersWithSelected,
          hasAccessToAddFilter: true,
          hasAccessToEditFilter: true,
          entitiesType: ENTITIES_TYPES.entity,
          action: expect.any(Function),
        },
      },
    );

    const [modalArguments] = $modals.show.mock.calls[0];

    const actionValue = [
      ...filtersWithSelected,
      {
        title: 'Filter 3',
        selected: false,
        filter: {},
      },
    ];

    modalArguments.config.action(actionValue);

    const updateFiltersEvents = wrapper.emitted('update:filters');

    expect(updateFiltersEvents).toHaveLength(1);

    const [changedFilters, newFilter] = updateFiltersEvents[0];

    expect(changedFilters).toEqual(
      actionValue.map(filter => omit(filter, ['selected'])),
    );
    expect(newFilter).toEqual(filters);
  });

  it('Filter changed to array after enabled multiple', () => {
    const [selectedFilter] = filters;
    const wrapper = factory({
      propsData: {
        filters,
        value: selectedFilter,
        hasAccessToUserFilter: true,
        entitiesType: ENTITIES_TYPES.entity,
      },
      mocks: {
        $modals,
      },
    });

    const mixFiltersButton = selectMixFiltersField(wrapper);

    mixFiltersButton.vm.$emit('input', true);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual([selectedFilter]);
  });

  it('Filter changed to array after enabled multiple with empty value', () => {
    const wrapper = factory({
      propsData: {
        filters,
        value: null,
        hasAccessToUserFilter: true,
        entitiesType: ENTITIES_TYPES.entity,
      },
      mocks: {
        $modals,
      },
    });

    const mixFiltersButton = selectMixFiltersField(wrapper);

    mixFiltersButton.vm.$emit('input', true);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual([]);
  });

  it('Filter changed to object after disabled multiple', () => {
    const [selectedFilter] = filters;
    const wrapper = factory({
      propsData: {
        filters,
        value: [selectedFilter],
        hasAccessToUserFilter: true,
        entitiesType: ENTITIES_TYPES.entity,
      },
      mocks: {
        $modals,
      },
    });

    const mixFiltersButton = selectMixFiltersField(wrapper);

    mixFiltersButton.vm.$emit('input', false);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual(selectedFilter);
  });

  it('Filter changed to object after disabled multiple with empty value', () => {
    const wrapper = factory({
      propsData: {
        filters,
        value: [],
        hasAccessToUserFilter: true,
        entitiesType: ENTITIES_TYPES.entity,
      },
      mocks: {
        $modals,
      },
    });

    const mixFiltersButton = selectMixFiltersField(wrapper);

    mixFiltersButton.vm.$emit('input', false);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual(null);
  });

  it('Filter doesn\'t changed after disabled multiple with object value', () => {
    const [selectedFilter] = filters;
    const wrapper = factory({
      propsData: {
        filters,
        value: selectedFilter,
        hasAccessToUserFilter: true,
        entitiesType: ENTITIES_TYPES.entity,
      },
      mocks: {
        $modals,
      },
    });

    const mixFiltersButton = selectMixFiltersField(wrapper);

    mixFiltersButton.vm.$emit('input', false);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toBeFalsy();
  });

  it('Condition changed after trigger operator field', () => {
    const wrapper = factory({
      propsData: {
        filters,
        value: {},
        condition: FILTER_MONGO_OPERATORS.or,
      },
      mocks: {
        $modals,
      },
    });

    const conditionField = selectConditionField(wrapper);

    conditionField.vm.$emit('input', FILTER_MONGO_OPERATORS.and);

    const updateConditionEvents = wrapper.emitted('update:condition');

    expect(updateConditionEvents).toHaveLength(1);

    const [eventData] = updateConditionEvents[0];
    expect(eventData).toEqual(FILTER_MONGO_OPERATORS.and);
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
        value: filters[0],
        label: 'Custom label',
        itemText: 'item-text',
        itemValue: 'item-value',
        condition: FILTER_MONGO_OPERATORS.or,
        hideSelect: true,
        hideSelectIcon: true,
        hasAccessToListFilters: true,
        hasAccessToAddFilter: false,
        hasAccessToEditFilter: false,
        hasAccessToUserFilter: false,
        entitiesType: ENTITIES_TYPES.entity,
      },
    });

    const menuContent = wrapper.find('.v-menu__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `filter-selector` with custom all rights', () => {
    const wrapper = snapshotFactory({
      propsData: {
        filters,
        lockedFilters,
        long: true,
        value: filters[0],
        condition: FILTER_MONGO_OPERATORS.or,
        hideSelect: true,
        hideSelectIcon: true,
        hasAccessToListFilters: true,
        hasAccessToAddFilter: true,
        hasAccessToEditFilter: true,
        hasAccessToUserFilter: true,
        entitiesType: ENTITIES_TYPES.entity,
      },
    });

    const menuContent = wrapper.find('.v-menu__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `filter-selector` with value as array', () => {
    const wrapper = snapshotFactory({
      propsData: {
        filters,
        lockedFilters,
        value: filters,
        condition: FILTER_MONGO_OPERATORS.or,
        hideSelect: false,
      },
    });

    const menuContent = wrapper.find('.v-menu__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });
});
