import {
  selectRowRemoveButtonByIndex,
  selectRowEditButtonByIndex,
  selectRowDuplicateButtonByIndex,
} from '@unit/utils/table';
import { generateRenderer } from '@unit/utils/vue';
import { range } from 'lodash';

import { MAP_TYPES } from '@/constants';

import CAdvancedDataTable from '@/components/common/table/c-advanced-data-table.vue';
import MapsList from '@/components/other/map/maps-list.vue';

const stubs = {
  'c-advanced-data-table': CAdvancedDataTable,
  'c-table-mass-actions-panel': true,
  'c-search': true,
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

  const factory = generateRenderer(MapsList, { stubs });

  test('Remove event emitted after trigger click on the remove button', async () => {
    const wrapper = factory({
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
    const wrapper = factory({
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
    const wrapper = factory({
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
});
