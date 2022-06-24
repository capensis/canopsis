import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { mockModals } from '@unit/utils/mock-hooks';
import { ENTITIES_TYPES } from '@/constants';

import FiltersList from '@/components/other/filter/filters-list.vue';

const localVue = createVueInstance();

const stubs = {
  'c-draggable-list-field': true,
  'filter-tile': true,
};

const snapshotStubs = {
  'c-draggable-list-field': true,
  'filter-tile': true,
};

const factory = (options = {}) => shallowMount(FiltersList, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(FiltersList, {
  localVue,
  stubs: snapshotStubs,

  ...options,
});

const selectAddButton = wrapper => wrapper.find('v-btn-stub');
const selectFilterTiles = wrapper => wrapper.findAll('filter-tile-stub');

describe('filters-list', () => {
  const $modals = mockModals();
  const filters = [
    { _id: '1', title: 'Filter 1' },
    { _id: '2', title: 'Filter 2' },
  ];

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
    expect(editEvents[0]).toEqual([filters[editedIndex]]);
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
    expect(deleteEvents[0]).toEqual([filters[deleteIndex]]);
  });

  it('Renders `filters-list` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
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

    expect(wrapper.element).toMatchSnapshot();
  });
});
