import {
  selectRowRemoveButtonByIndex,
  selectRowEditButtonByIndex,
  selectRowCheckboxByIndex,
  selectMassRemoveButton,
  selectRowDuplicateButtonByIndex,
} from '@unit/utils/table';
import { generateRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import { range } from 'lodash';

import { TAG_TYPES } from '@/constants';

import CAdvancedDataTable from '@/components/common/table/c-advanced-data-table.vue';
import TagsList from '@/components/other/tag/tags-list.vue';

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
  'c-chip': true,
  'tags-list-expand-panel': true,
};

describe('tags-list', () => {
  const totalItems = 11;

  const types = [TAG_TYPES.created, TAG_TYPES.imported];

  const tags = range(totalItems).map(index => ({
    _id: `pattern-id-${index}`,
    value: `pattern-value-${index}`,
    type: types[index % 2],
    enabled: !!(index % 2),
    deletable: true,
    author: {
      display_name: `author-${index}`,
    },
    created: 1614851888 + index,
    updated: 1614861888 + index,
  }));

  const factory = generateRenderer(TagsList, { stubs });

  it('Selected items removed after trigger mass remove button', async () => {
    const wrapper = factory({
      propsData: {
        tags,
        options: {
          page: 1,
          itemsPerPage: 10,
        },
        totalItems,
        removable: true,
      },
    });

    await flushPromises();

    await selectRowCheckboxByIndex(wrapper, 1).trigger('click', true);
    await selectRowCheckboxByIndex(wrapper, 2).trigger('click', true);
    await selectMassRemoveButton(wrapper).triggerCustomEvent('click');

    expect(wrapper).toEmit('remove-selected', [tags[1], tags[2]]);
  });

  it('Remove event emitted after trigger click on the remove button', async () => {
    const wrapper = factory({
      propsData: {
        tags,
        options: {
          page: 1,
          itemsPerPage: 10,
        },
        totalItems,
        removable: true,
      },
    });

    const removeRowIndex = 2;

    await selectRowRemoveButtonByIndex(wrapper, removeRowIndex).triggerCustomEvent('click');

    expect(wrapper).toEmit('remove', tags[removeRowIndex]._id);
  });

  it('Update event emitted after trigger click on the remove button', async () => {
    const wrapper = factory({
      propsData: {
        tags,
        options: {
          page: 1,
          itemsPerPage: 10,
        },
        totalItems,
        updatable: true,
      },
    });

    const editRowIndex = 5;

    await selectRowEditButtonByIndex(wrapper, editRowIndex).triggerCustomEvent('click');

    expect(wrapper).toEmit('edit', tags[editRowIndex]);
  });

  it('Duplicate event emitted after trigger click on the duplicate button', async () => {
    const wrapper = factory({
      propsData: {
        tags,
        options: {
          page: 1,
          itemsPerPage: 10,
        },
        totalItems,
        duplicable: true,
      },
    });

    const duplicateRowIndex = 5;

    await selectRowDuplicateButtonByIndex(wrapper, duplicateRowIndex).triggerCustomEvent('click');

    expect(wrapper).toEmit('duplicate', tags[duplicateRowIndex]);
  });
});
