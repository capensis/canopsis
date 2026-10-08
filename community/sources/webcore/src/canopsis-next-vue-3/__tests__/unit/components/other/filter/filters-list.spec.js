import { mockModals } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';

import { ENTITIES_TYPES } from '@/constants';

import FiltersList from '@/components/other/filter/filters-list.vue';

const stubs = {
  'c-draggable-list-field': true,
  'filter-tile': true,
  'c-alert': true,
};

const selectAddButton = wrapper => wrapper.find('v-btn-stub');

describe('filters-list', () => {
  const $modals = mockModals();
  const filters = [
    { _id: '1', title: 'Filter 1' },
    { _id: '2', title: 'Filter 2' },
  ];

  const factory = generateShallowRenderer(FiltersList, { stubs });

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

    selectAddButton(wrapper).triggerCustomEvent('click', new Event('click'));

    expect(wrapper).toHaveBeenEmit('add');
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

    wrapper.vm.$emit('edit', filters[editedIndex], editedIndex);

    expect(wrapper.emitted('edit')).toEqual([[filters[editedIndex], editedIndex]]);
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

    wrapper.vm.$emit('delete', filters[deleteIndex], deleteIndex);

    expect(wrapper.emitted('delete')).toEqual([[filters[deleteIndex], deleteIndex]]);
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

    wrapper.vm.$emit('update:filters', updatedFilters);

    expect(wrapper.emitted('update:filters')).toEqual([[updatedFilters]]);
  });
});
