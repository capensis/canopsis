import { generateShallowRenderer } from '@unit/utils/vue';
import { range } from 'lodash';

import { PATTERN_TYPES } from '@/constants';

import PatternsList from '@/components/other/pattern/patterns-list.vue';

const stubs = {
  'c-advanced-data-table': {
    name: 'CAdvancedDataTable',
    props: ['items'],
    template: `
      <div>
        <div class="mass-actions">
          <slot name="mass-actions" :selected="[items[0], items[2]]" />
        </div>
        <div v-for="item in items" :key="item._id" class="row-actions">
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
};

const selectMassRemoveButton = wrapper =>
  wrapper.find('.mass-actions').findComponent({ name: 'CActionBtn' });
const selectRowActionButtons = (wrapper, index) =>
  wrapper.findAll('.row-actions').at(index).findAllComponents({ name: 'CActionBtn' });
const selectRowEditButtonByIndex = (wrapper, index) => selectRowActionButtons(wrapper, index).at(0);
const selectRowRemoveButtonByIndex = (wrapper, index) =>
  selectRowActionButtons(wrapper, index).at(1);

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

  const factory = generateShallowRenderer(PatternsList, { stubs });

  it('Selected items removed after trigger mass remove button', async () => {
    const wrapper = factory({
      props: {
        patterns,
        options: {
          page: 1,
          itemsPerPage: 10,
        },
        totalItems,
      },
    });

    selectMassRemoveButton(wrapper).vm.$emit('click');
    await wrapper.vm.$nextTick();

    expect(wrapper.emitted('remove-selected')).toEqual([[[patterns[0], patterns[2]]]]);
  });

  it('Remove event emitted after trigger click on the remove button', async () => {
    const wrapper = factory({
      props: {
        patterns,
        options: {
          page: 1,
          itemsPerPage: 10,
        },
        totalItems,
      },
    });

    const removableRowIndex = 2;

    selectRowRemoveButtonByIndex(wrapper, removableRowIndex).vm.$emit('click');
    await wrapper.vm.$nextTick();

    expect(wrapper.emitted('remove')).toEqual([[patterns[removableRowIndex]._id]]);
  });

  it('Update event emitted after trigger click on the remove button', async () => {
    const wrapper = factory({
      props: {
        patterns,
        options: {
          page: 1,
          itemsPerPage: 10,
        },
        totalItems,
      },
    });

    const editableRowIndex = 5;

    selectRowEditButtonByIndex(wrapper, editableRowIndex).vm.$emit('click');
    await wrapper.vm.$nextTick();

    expect(wrapper.emitted('edit')).toEqual([[patterns[editableRowIndex]]]);
  });
});
