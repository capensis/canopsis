<template>
  <c-card-iterator-item
    :drag-handle-class="dragHandleClass"
    small
    @remove="$emit('remove')"
  >
    <template #header>
      <v-layout class="gap-2" column>
        <v-select
          v-validate="'required'"
          :value="sortColumn.sort_by"
          :items="availableColumns"
          :label="$t('settings.columnName')"
          :error-messages="errors.collect(sortByName)"
          :name="sortByName"
          :menu-props="selectMenuProps"
          @change="changeSortBy"
        />
        <v-select
          v-validate="'required'"
          :value="sortColumn.sort"
          :items="sortOrderItems"
          :label="$t('common.sort')"
          :error-messages="errors.collect(sortName)"
          :name="sortName"
          :menu-props="selectMenuProps"
          @change="changeSort"
        />
      </v-layout>
    </template>
  </c-card-iterator-item>
</template>

<script>
import { computed, toRef } from 'vue';

import { ENTITIES_TYPES, SORT_ORDERS } from '@/constants';

import { useModelField } from '@/hooks/form/model-field';
import { useAvailableColumns } from '@/hooks/form/available-columns';

export default {
  inject: ['$validator'],
  model: {
    prop: 'sortColumn',
    event: 'input',
  },
  props: {
    sortColumn: {
      type: Object,
      default: () => ({}),
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
  setup(props, { emit }) {
    const { updateField } = useModelField(props, emit);

    const selectMenuProps = {
      contentClass: 'sort-column-field-menu',
    };

    const { availableColumns } = useAvailableColumns({
      type: ENTITIES_TYPES.alarm,
      items: toRef(props, 'items'),
    });

    const sortOrderItems = Object.values(SORT_ORDERS);

    const sortByName = computed(() => `${props.name}.sort_by`);
    const sortName = computed(() => `${props.name}.sort`);

    const changeSortBy = sortBy => updateField('sort_by', sortBy);
    const changeSort = sort => updateField('sort', sort);

    return {
      availableColumns,
      sortOrderItems,
      selectMenuProps,
      sortByName,
      sortName,
      changeSortBy,
      changeSort,
    };
  },
};
</script>
