import { generateRenderer } from '@unit/utils/vue';
import { range } from 'lodash';

import { TAG_TYPES } from '@/constants';

import TagsList from '@/components/other/tag/tags-list.vue';

const CAdvancedDataTableStub = {
  props: ['items'],
  data: () => ({ selected: [] }),
  methods: {
    toggle(item) {
      this.selected = this.selected.includes(item)
        ? this.selected.filter(selectedItem => selectedItem !== item)
        : [...this.selected, item];
    },
  },
  template: `
    <div class="advanced-table">
      <div class="mass-actions">
        <slot name="mass-actions" :selected="selected" :clear-selected="() => selected = []" />
      </div>
      <div v-for="(item, index) in items" :key="item._id" class="table-row">
        <button class="row-checkbox" @click="toggle(item)">{{ index }}</button>
        <slot name="actions" :item="item" />
      </div>
    </div>
  `,
};

const CActionBtnStub = {
  props: ['type'],
  emits: ['click'],
  template: '<button :class="[`action-${type}`]" @click="$emit(\'click\', $event)" />',
};

const stubs = {
  'c-advanced-data-table': CAdvancedDataTableStub,
  'c-search': true,
  'v-checkbox': true,
  'c-expand-btn': true,
  'c-action-btn': CActionBtnStub,
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

    await wrapper.findAll('button.row-checkbox').at(1).trigger('click');
    await wrapper.findAll('button.row-checkbox').at(2).trigger('click');
    await wrapper.find('.mass-actions button.action-delete').trigger('click');

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

    await wrapper
      .findAll('.table-row')
      .at(removeRowIndex)
      .find('button.action-delete')
      .trigger('click');

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

    await wrapper
      .findAll('.table-row')
      .at(editRowIndex)
      .find('button.action-edit')
      .trigger('click');

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

    await wrapper
      .findAll('.table-row')
      .at(duplicateRowIndex)
      .find('button.action-duplicate')
      .trigger('click');

    expect(wrapper).toEmit('duplicate', tags[duplicateRowIndex]);
  });
});
