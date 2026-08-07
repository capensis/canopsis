<template>
  <c-card-iterator-item
    :drag-handle-class="dragHandleClass"
    class="sort-column-field"
    small
    @remove="$emit('remove')"
  >
    <template #header>
      <c-layout
        class="sort-column-field__inputs"
        gap="2"
        column
      >
        <c-select
          :model-value="currentSortColumn.sort_by"
          :items="availableColumns"
          :label="$t('settings.columnName')"
          :name="sortByName"
          item-title="text"
          item-value="value"
          @update:model-value="changeSortBy"
        />
        <c-select
          :model-value="currentSortColumn.sort"
          :items="sortOrderItems"
          :label="$t('common.sort')"
          :name="sortName"
          @update:model-value="changeSort"
        />
      </c-layout>
    </template>
  </c-card-iterator-item>
</template>

<script>
import { computed, toRef } from 'vue';

import { ENTITIES_TYPES, SORT_ORDERS } from '@/constants';

import { useAvailableColumns } from '@/hooks/form/available-columns';
import { useModelField } from '@/hooks/form/model-field';

export default {
  props: {
    modelValue: {
      type: Object,
      default: () => ({}),
    },
    sortColumn: {
      type: Object,
      default: undefined,
    },
    name: {
      type: String,
      default: '',
    },
    dragHandleClass: {
      type: String,
      default: 'item-drag-handler',
    },
    items: {
      type: Array,
      required: false,
    },
  },

  emits: ['update:modelValue', 'remove'],

  setup(props, { emit }) {
    const { updateField } = useModelField(props, emit);

    const { availableColumns } = useAvailableColumns({
      type: ENTITIES_TYPES.alarm,
      items: toRef(props, 'items'),
    });

    const sortOrderItems = Object.values(SORT_ORDERS);
    const sortByName = computed(() => `${props.name}.sort_by`);
    const sortName = computed(() => `${props.name}.sort`);
    const currentSortColumn = computed(() => props.sortColumn ?? props.modelValue);

    const changeSortBy = sortBy => updateField('sort_by', sortBy);
    const changeSort = sort => updateField('sort', sort);

    return {
      availableColumns,
      sortOrderItems,
      sortByName,
      sortName,
      currentSortColumn,
      changeSortBy,
      changeSort,
    };
  },
};
</script>

<style lang="scss" scoped>
.sort-column-field {
  :deep(.c-card-iterator-item__row) {
    align-items: center;
  }

  :deep(.c-card-iterator-item__actions) {
    align-items: center;
    align-self: stretch;
  }

  &__inputs {
    flex: 1 1 auto;
    min-width: 0;
    width: 100%;

    :deep(.v-input) {
      width: 100%;
    }
  }
}
</style>
