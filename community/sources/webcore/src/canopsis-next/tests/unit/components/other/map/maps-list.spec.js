import { range } from 'lodash';

import { mount, createVueInstance } from '@unit/utils/vue';
import {
  selectRowRemoveButtonByIndex,
  selectRowEditButtonByIndex,
  selectRowCheckboxByIndex,
  selectMassRemoveButton,
  selectRowDuplicateButtonByIndex,
  selectRowExpandButtonByIndex,
} from '@unit/utils/table';

import { MAP_TYPES } from '@/constants';
import MapsList from '@/components/other/map/maps-list.vue';
import CAdvancedDataTable from '@/components/common/table/c-advanced-data-table.vue';

const localVue = createVueInstance();

const stubs = {
  'c-advanced-data-table': CAdvancedDataTable,
  'c-search-field': true,
  'v-checkbox': true,
  'v-checkbox-functional': true,
  'c-expand-btn': true,
  'c-action-btn': true,
  'c-table-pagination': true,
  'maps-list-expand-item': true,
};

const snapshotFactory = (options = {}) => mount(MapsList, {
  localVue,
  stubs,

  ...options,
});

describe('maps-list', () => {
  const totalItems = 11;

  const types = Object.values(MAP_TYPES);

  const maps = range(totalItems).map(index => ({
    _id: `map-id-${index}`,
    type: types[index % 3],
    author: {
      name: `author-${index}`,
    },
    deletable: true,
    updated: 1614861888 + index,
  }));

  test('Selected items removed after trigger mass remove button', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        maps,
        removable: true,
        pagination: {
          page: 1,
          rowsPerPage: 10,
        },
        totalItems,
      },
    });

    const firstCheckbox = selectRowCheckboxByIndex(wrapper, 0);
    await firstCheckbox.vm.$emit('change', true);

    const thirdCheckbox = selectRowCheckboxByIndex(wrapper, 2);
    await thirdCheckbox.vm.$emit('change', true);

    const massRemoveButton = selectMassRemoveButton(wrapper);

    await massRemoveButton.vm.$emit('click');

    expect(wrapper).toEmit('remove-selected', [maps[0], maps[2]]);
  });

  test('Remove event emitted after trigger click on the remove button', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        maps,
        removable: true,
        pagination: {
          page: 1,
          rowsPerPage: 10,
        },
        totalItems,
      },
    });

    const removableRowIndex = 2;

    const removeButton = selectRowRemoveButtonByIndex(wrapper, removableRowIndex);
    await removeButton.vm.$emit('click');

    expect(wrapper).toEmit('remove', maps[removableRowIndex]._id);
  });

  test('Update event emitted after trigger click on the edit button', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        maps,
        updatable: true,
        pagination: {
          page: 1,
          rowsPerPage: 10,
        },
        totalItems,
      },
    });

    const editableRowIndex = 5;

    const editButton = selectRowEditButtonByIndex(wrapper, editableRowIndex);
    await editButton.vm.$emit('click');

    expect(wrapper).toEmit('edit', maps[editableRowIndex]);
  });

  test('Duplicate event emitted after trigger click on the duplicate button', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        maps,
        duplicable: true,
        pagination: {
          page: 1,
          rowsPerPage: 10,
        },
        totalItems,
      },
    });

    const editableRowIndex = 5;

    const duplicateButton = selectRowDuplicateButtonByIndex(wrapper, editableRowIndex);
    await duplicateButton.vm.$emit('click');

    expect(wrapper).toEmit('duplicate', maps[editableRowIndex]);
  });

  test('Renders `maps-list` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        maps,
        pagination: {},
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `maps-list` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        maps,
        pagination: {
          page: 2,
          rowsPerPage: 10,
          search: 'Filter',
          sortBy: 'updated',
          descending: true,
        },
        totalItems,
        pending: true,
        removable: true,
        updatable: true,
        duplicable: true,
      },
    });

    const expandButton = selectRowExpandButtonByIndex(wrapper, 1);

    await expandButton.vm.$emit('expand');

    expect(wrapper.element).toMatchSnapshot();
  });
});
