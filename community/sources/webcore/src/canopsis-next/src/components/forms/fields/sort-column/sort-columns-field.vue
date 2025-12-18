<template>
  <v-layout column>
    <c-card-iterator-field
      v-field="sortColumns"
      :class="{ empty: isSortColumnsEmpty }"
      :handle="`.${dragItemHandleClass}`"
      item-key="key"
    >
      <template #item="{ item: sortColumn, index }">
        <sort-column-field
          v-field="sortColumns[index]"
          :key="sortColumn.key"
          :name="sortColumn.key"
          :drag-handle-class="dragItemHandleClass"
          :items="items"
          class="mb-3"
          @remove="remove(index)"
        />
      </template>
    </c-card-iterator-field>
    <v-layout justify-end>
      <v-tooltip left>
        <template #activator="{ on }">
          <v-btn
            class="mt-3"
            color="primary"
            fab
            small
            v-on="on"
            @click.prevent="add"
          >
            <v-icon>add</v-icon>
          </v-btn>
        </template>
        <span>{{ $t('common.add') }}</span>
      </v-tooltip>
    </v-layout>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { widgetSortColumnToForm } from '@/helpers/entities/widget/sort-column/form';

import { useArrayModelField } from '@/hooks/form/array-model-field';

import SortColumnField from './partials/sort-column-field.vue';

export default {
  inject: ['$validator'],
  components: {
    SortColumnField,
  },
  model: {
    prop: 'sortColumns',
    event: 'input',
  },
  props: {
    sortColumns: {
      type: Array,
      default: () => [],
    },
    items: {
      type: Array,
      required: false,
    },
  },
  setup(props, { emit }) {
    const dragItemHandleClass = 'sort-column-drag-handler';

    const isSortColumnsEmpty = computed(() => !props.sortColumns?.length);

    const { addItemIntoArray, removeItemFromArray } = useArrayModelField(props, emit);

    const add = () => addItemIntoArray(widgetSortColumnToForm());

    const remove = index => removeItemFromArray(index);

    return {
      dragItemHandleClass,
      isSortColumnsEmpty,
      add,
      remove,
    };
  },
};
</script>
