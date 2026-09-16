import {
  selectRowRemoveButtonByIndex,
  selectRowEditButtonByIndex,
  selectRowCheckboxByIndex,
  selectMassRemoveButton,
} from '@unit/utils/table';
import { generateRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import { range } from 'lodash';

import { PATTERN_TYPES } from '@/constants';

import CAdvancedDataTable from '@/components/common/table/c-advanced-data-table.vue';
import PatternsList from '@/components/other/pattern/patterns-list.vue';

const cMassActionsPanelStub = {
  name: 'c-mass-actions-panel-stub',
  template:
    '<div class="c-mass-actions-panel-stub"><div class="layout"><slot name="actions" /></div></div>',
};

const stubs = {
  'c-advanced-data-table': CAdvancedDataTable,
  'c-mass-actions-panel': cMassActionsPanelStub,
  'c-search': true,
  'v-checkbox': true,
  'c-expand-btn': true,
  'c-action-btn': true,
  'c-table-pagination': true,
};

describe('patterns-list', () => {
  const totalItems = 11;

  const types = [PATTERN_TYPES.alarm, PATTERN_TYPES.entity, PATTERN_TYPES.pbehavior];

  const patterns = range(totalItems).map(index => ({
    _id: `pattern-id-${index}`,
    type: types[index % 3],
    enabled: !!(index % 2),
    author: {
      display_name: `author-${index}`,
    },
    updated: 1614861888 + index,
  }));

  const factory = generateRenderer(PatternsList, { stubs });

  it('Selected items removed after trigger mass remove button', async () => {
    const wrapper = factory({
      propsData: {
        patterns,
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

    const massRemoveButton = selectMassRemoveButton(wrapper);
    await massRemoveButton.triggerCustomEvent('click');

    expect(wrapper).toEmit('remove-selected', [patterns[0], patterns[2]]);
  });

  it('Remove event emitted after trigger click on the remove button', async () => {
    const wrapper = factory({
      propsData: {
        patterns,
        options: {
          page: 1,
          itemsPerPage: 10,
        },
        totalItems,
      },
    });

    await flushPromises();

    const removableRowIndex = 2;

    await selectRowRemoveButtonByIndex(wrapper, removableRowIndex).triggerCustomEvent('click');

    expect(wrapper).toEmit('remove', patterns[removableRowIndex]._id);
  });

  it('Update event emitted after trigger click on the remove button', async () => {
    const wrapper = factory({
      propsData: {
        patterns,
        options: {
          page: 1,
          itemsPerPage: 10,
        },
        totalItems,
      },
    });

    await flushPromises();

    const editableRowIndex = 5;

    await selectRowEditButtonByIndex(wrapper, editableRowIndex).triggerCustomEvent('click');

    expect(wrapper).toEmit('edit', patterns[editableRowIndex]);
  });
});
