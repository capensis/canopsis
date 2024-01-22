import { range } from 'lodash';

import { flushPromises, generateRenderer } from '@unit/utils/vue';
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

const stubs = {
  'c-advanced-data-table': CAdvancedDataTable,
  'c-search-field': true,
  'v-checkbox': true,
  'c-expand-btn': true,
  'c-action-btn': true,
  'c-table-pagination': true,
  'maps-list-expand-item': true,
};

describe('maps-list', () => {
  const totalItems = 11;

  const types = Object.values(MAP_TYPES);

  const maps = range(totalItems).map(index => ({
    _id: `map-id-${index}`,
    type: types[index % 3],
    author: {
      display_name: `author-${index}`,
    },
    deletable: true,
    updated: 1614861888 + index,
  }));

  const snapshotFactory = generateRenderer(MapsList, { stubs });

  test('Selected items removed after trigger mass remove button', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        maps,
        removable: true,
        options: {
          page: 1,
          itemsPerPage: 10,
        },
        totalItems,
      },
    });

    await flushPromises();

    await selectRowCheckboxByIndex(wrapper, 0).trigger('click');
    await selectRowCheckboxByIndex(wrapper, 2).trigger('click');
    await selectMassRemoveButton(wrapper).triggerCustomEvent('click');

    expect(wrapper).toEmit('remove-selected', [maps[0], maps[2]]);
  });

  test('Remove event emitted after trigger click on the remove button', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        maps,
        removable: true,
        options: {
          page: 1,
          itemsPerPage: 10,
        },
        totalItems,
      },
    });

    const removableRowIndex = 2;

    const removeButton = selectRowRemoveButtonByIndex(wrapper, removableRowIndex);
    await removeButton.triggerCustomEvent('click');

    expect(wrapper).toEmit('remove', maps[removableRowIndex]._id);
  });

  test('Update event emitted after trigger click on the edit button', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        maps,
        updatable: true,
        options: {
          page: 1,
          itemsPerPage: 10,
        },
        totalItems,
      },
    });

    const editableRowIndex = 5;

    const editButton = selectRowEditButtonByIndex(wrapper, editableRowIndex);
    await editButton.triggerCustomEvent('click');

    expect(wrapper).toEmit('edit', maps[editableRowIndex]);
  });

  test('Duplicate event emitted after trigger click on the duplicate button', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        maps,
        duplicable: true,
        options: {
          page: 1,
          itemsPerPage: 10,
        },
        totalItems,
      },
    });

    const editableRowIndex = 5;

    const duplicateButton = selectRowDuplicateButtonByIndex(wrapper, editableRowIndex);
    await duplicateButton.triggerCustomEvent('click');

    expect(wrapper).toEmit('duplicate', maps[editableRowIndex]);
  });

  test('Renders `maps-list` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        maps,
        options: {},
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `maps-list` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        maps,
        options: {
          page: 2,
          itemsPerPage: 10,
          search: 'Filter',
          sortBy: ['updated'],
          sortDesc: [true],
        },
        totalItems,
        pending: true,
        removable: true,
        updatable: true,
        duplicable: true,
      },
    });

    await selectRowExpandButtonByIndex(wrapper, 1).trigger('click');

    expect(wrapper).toMatchSnapshot();
  });
});
