import { generateShallowRenderer } from '@unit/utils/vue';
import { range } from 'lodash';

import { MAP_TYPES } from '@/constants';

import MapsList from '@/components/other/map/maps-list.vue';

const stubs = {
  'c-advanced-data-table': {
    name: 'CAdvancedDataTable',
    props: ['items'],
    template: `
      <div>
        <div v-for="item in items" :key="item._id" class="map-actions">
          <slot name="actions" :item="item" />
        </div>
      </div>
    `,
  },
  'c-action-btn': {
    name: 'CActionBtn',
    props: ['type'],
    emits: ['click'],
    template: '<button @click="$emit(\'click\', $event)" />',
  },
  'c-table-mass-actions-panel': true,
};

const selectActionButtons = wrapper => wrapper.findAllComponents({ name: 'CActionBtn' });
const selectActionButtonByIndex = (wrapper, index) => selectActionButtons(wrapper).at(index);

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

  const factory = generateShallowRenderer(MapsList, { stubs });

  test('Remove event emitted after trigger click on the remove button', async () => {
    const wrapper = factory({
      props: {
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

    await selectActionButtonByIndex(wrapper, removableRowIndex).trigger('click');

    expect(wrapper.emitted('remove')).toEqual([[maps[removableRowIndex]._id]]);
  });

  test('Update event emitted after trigger click on the edit button', async () => {
    const wrapper = factory({
      props: {
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

    await selectActionButtonByIndex(wrapper, editableRowIndex).trigger('click');

    expect(wrapper.emitted('edit')).toEqual([[maps[editableRowIndex]]]);
  });

  test('Duplicate event emitted after trigger click on the duplicate button', async () => {
    const wrapper = factory({
      props: {
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

    await selectActionButtonByIndex(wrapper, editableRowIndex).trigger('click');

    expect(wrapper.emitted('duplicate')).toEqual([[maps[editableRowIndex]]]);
  });
});
