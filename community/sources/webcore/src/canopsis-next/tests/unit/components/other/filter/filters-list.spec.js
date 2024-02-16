import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { mockModals } from '@unit/utils/mock-hooks';
import { ENTITIES_TYPES } from '@/constants';

import FiltersList from '@/components/other/filter/filters-list.vue';

const stubs = {
  'c-draggable-list-field': true,
  'filter-tile': true,
  'c-alert': true,
};

const snapshotStubs = {
  'c-draggable-list-field': true,
  'filter-tile': true,
  'c-alert': true,
};

const selectAddButton = wrapper => wrapper.find('v-btn-stub');
const selectFilterTiles = wrapper => wrapper.findAll('filter-tile-stub');
const selectDraggableField = wrapper => wrapper.find('c-draggable-list-field-stub');

describe('filters-list', () => {
  const $modals = mockModals();
  const filters = [
    { _id: '1', title: 'Filter 1' },
    { _id: '2', title: 'Filter 2' },
  ];

  const factory = generateShallowRenderer(FiltersList, { stubs });
  const snapshotFactory = generateRenderer(FiltersList, { stubs: snapshotStubs });

  it('Create filter modal opened after trigger add button', () => {
    const wrapper = factory({
      propsData: {
        filters,
        addable: true,
        withAlarm: true,
        withEntity: true,
      },
      mocks: {
        $modals,
      },
    });

    const addButton = selectAddButton(wrapper);

    addButton.vm.$emit('click', new Event('click'));

    const addEvents = wrapper.emitted('add');

    expect(addEvents).toHaveLength(1);
  });

  it('Edit filter modal opened after trigger edit button', () => {
    const wrapper = factory({
      propsData: {
        filters,
        addable: true,
        withEntity: true,
      },
      mocks: {
        $modals,
      },
    });

    const editedIndex = 1;
    const filterTile = selectFilterTiles(wrapper).at(editedIndex);
    filterTile.vm.$emit('edit');

    const editEvents = wrapper.emitted('edit');

    expect(editEvents).toHaveLength(1);
    expect(editEvents[0]).toEqual([filters[editedIndex], editedIndex]);
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

    const deleteIndex = 1;
    const filterTile = selectFilterTiles(wrapper).at(deleteIndex);

    filterTile.vm.$emit('delete');

    const deleteEvents = wrapper.emitted('delete');

    expect(deleteEvents).toHaveLength(1);
    expect(deleteEvents[0]).toEqual([filters[deleteIndex], deleteIndex]);
  });

  it('Should send updated filters array on dragging', () => {
    const wrapper = factory({
      propsData: {
        filters,
      },
      mocks: {
        $modals,
      },
    });

    const updatedFilters = [...filters].reverse();
    const draggableField = selectDraggableField(wrapper);

    draggableField.vm.$emit('input', updatedFilters);

    expect(wrapper).toEmit('input', updatedFilters);
  });

  it('Renders `filters-list` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `filters-list` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        filters,
        pending: true,
        addable: true,
        editable: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
