import { range } from 'lodash';

import { mount, createVueInstance } from '@unit/utils/vue';
import {
  selectRowRemoveButtonByIndex,
  selectRowEditButtonByIndex,
  selectRowCheckboxByIndex,
  selectMassRemoveButton,
} from '@unit/utils/table';

import { PATTERN_TYPES } from '@/constants';
import PatternsList from '@/components/other/pattern/patterns-list.vue';
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
};

const snapshotFactory = (options = {}) => mount(PatternsList, {
  localVue,
  stubs,

  ...options,
});

describe('patterns-list', () => {
  const totalItems = 11;

  const types = [PATTERN_TYPES.alarm, PATTERN_TYPES.entity, PATTERN_TYPES.pbehavior];

  const patterns = range(totalItems).map(index => ({
    _id: `pattern-id-${index}`,
    type: types[index % 3],
    enabled: !!(index % 2),
    author: {
      name: `author-${index}`,
    },
    updated: 1614861888 + index,
  }));

  it('Selected items removed after trigger mass remove button', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        patterns,
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

    const removeSelectedEvent = wrapper.emitted('remove-selected');

    expect(removeSelectedEvent).toHaveLength(1);

    const [eventData] = removeSelectedEvent[0];
    expect(eventData).toEqual(
      [patterns[0], patterns[2]],
    );
  });

  it('Remove event emitted after trigger click on the remove button', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        patterns,
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

    const editEvent = wrapper.emitted('remove');

    expect(editEvent).toHaveLength(1);

    const [eventData] = editEvent[0];
    expect(eventData).toEqual(patterns[removableRowIndex]._id);
  });

  it('Update event emitted after trigger click on the remove button', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        patterns,
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

    const editEvent = wrapper.emitted('edit');

    expect(editEvent).toHaveLength(1);

    const [eventData] = editEvent[0];
    expect(eventData).toEqual(patterns[editableRowIndex]);
  });

  it('Renders `patterns-list` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        patterns,
        pagination: {},
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `patterns-list` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        patterns,
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
        corporate: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
