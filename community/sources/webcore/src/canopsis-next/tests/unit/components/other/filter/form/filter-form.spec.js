import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { mockModals } from '@unit/utils/mock-hooks';
import { ENTITIES_TYPES, MODALS } from '@/constants';

import FiltersForm from '@/components/other/filter/form/filters-form.vue';

const localVue = createVueInstance();

const stubs = {
  'c-draggable-list-field': true,
  'filter-tile': true,
};

const snapshotStubs = {
  'c-draggable-list-field': true,
  'filter-tile': true,
};

const factory = (options = {}) => shallowMount(FiltersForm, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(FiltersForm, {
  localVue,
  stubs: snapshotStubs,

  ...options,
});

const selectAddButton = wrapper => wrapper.find('v-btn-stub');
const selectFilterTiles = wrapper => wrapper.findAll('filter-tile-stub');
const selectDraggableField = wrapper => wrapper.find('c-draggable-list-field-stub');

describe('filters-form', () => {
  const $modals = mockModals();
  const filters = [
    { title: 'Filter 1' },
    { title: 'Filter 2' },
  ];
  const filtersTitles = filters.map(({ title }) => title);

  it('Create filter modal opened after trigger add button', () => {
    const wrapper = factory({
      propsData: {
        filters,
        addable: true,
        entitiesType: ENTITIES_TYPES.entity,
      },
      mocks: {
        $modals,
      },
    });

    const addButton = selectAddButton(wrapper);

    addButton.vm.$emit('click', new Event('click'));

    expect($modals.show).toBeCalledWith({
      name: MODALS.createFilter,
      config: {
        title: 'modals.filter.create.title',
        entitiesType: ENTITIES_TYPES.entity,
        existingTitles: filtersTitles,
        action: expect.any(Function),
      },
    });
    const [{ config }] = $modals.show.mock.calls[0];

    const newFilter = {
      title: Faker.datatype.string(),
    };

    config.action(newFilter);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);
    const [eventData] = inputEvents[0];
    expect(eventData).toEqual([
      ...filters,
      newFilter,
    ]);
  });

  it('Edit filter modal opened after trigger edit button', () => {
    const wrapper = factory({
      propsData: {
        filters,
        addable: true,
        entitiesType: ENTITIES_TYPES.entity,
      },
      mocks: {
        $modals,
      },
    });

    const editedIndex = 1;

    const filterTile = selectFilterTiles(wrapper).at(editedIndex);

    filterTile.vm.$emit('edit');

    expect($modals.show).toBeCalledWith({
      name: MODALS.createFilter,
      config: {
        filter: filters[editedIndex],
        title: 'modals.filter.edit.title',
        entitiesType: ENTITIES_TYPES.entity,
        existingTitles: filtersTitles,
        action: expect.any(Function),
      },
    });
    const [{ config }] = $modals.show.mock.calls[0];

    const updatedFilter = {
      title: Faker.datatype.string(),
    };

    config.action(updatedFilter);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);
    const [eventData] = inputEvents[0];

    expect(eventData).toEqual([
      ...filters.slice(0, 1),
      updatedFilter,
    ]);
  });

  it('Confirmation modal opened after trigger delete button', () => {
    const wrapper = factory({
      propsData: {
        filters,
        addable: true,
        entitiesType: ENTITIES_TYPES.entity,
      },
      mocks: {
        $modals,
      },
    });

    const filterTile = selectFilterTiles(wrapper).at(1);

    filterTile.vm.$emit('delete');

    expect($modals.show).toBeCalledWith({
      name: MODALS.confirmation,
      config: {
        action: expect.any(Function),
      },
    });
    const [{ config }] = $modals.show.mock.calls[0];

    config.action();

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);
    const [eventData] = inputEvents[0];

    expect(eventData).toEqual(filters.slice(0, 1));
  });

  it('Filters changed after trigger draggable list field', () => {
    const wrapper = factory({
      propsData: {
        filters,
        entitiesType: ENTITIES_TYPES.entity,
      },
    });

    const draggableField = selectDraggableField(wrapper);

    const updatedFilters = filters.slice().reverse();

    draggableField.vm.$emit('input', updatedFilters);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);
    const [eventData] = inputEvents[0];

    expect(eventData).toEqual(updatedFilters);
  });

  it('Renders `filters-form` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `filters-form` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        filters,
        addable: true,
        editable: true,
        entitiesType: ENTITIES_TYPES.entity,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
