<template>
  <c-layout column>
    <c-card-iterator-field
      :data="modelValue"
      :class="{ empty: isSortColumnsEmpty }"
      :handle="`.${dragItemHandleClass}`"
      item-key="key"
      @update:data="$emit('update:modelValue', $event)"
    >
      <template #item="{ item: sortColumn, index }">
        <sort-column-field
          :key="sortColumn.key"
          :model-value="modelValue[index]"
          :name="sortColumn.key"
          :drag-handle-class="dragItemHandleClass"
          :items="items"
          class="mb-3"
          @update:model-value="updateItemInArray(index, $event)"
          @remove="remove(index)"
        />
      </template>
    </c-card-iterator-field>

    <c-layout justify="end">
      <v-tooltip location="left">
        <template #activator="{ props }">
          <v-btn
            class="mt-3"
            color="primary"
            size="small"
            variant="outlined"
            icon
            v-bind="props"
            @click.prevent="add"
          >
            <v-icon>add</v-icon>
          </v-btn>
        </template>
        <span>{{ $t('common.add') }}</span>
      </v-tooltip>
    </c-layout>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { widgetSortColumnToForm } from '@/helpers/entities/widget/sort-column/form';

import { useArrayModelField } from '@/hooks/form/array-model-field';

import SortColumnField from './partials/sort-column-field.vue';

export default {
  components: {
    SortColumnField,
  },

  props: {
    modelValue: {
      type: Array,
      default: () => [],
    },
    items: {
      type: Array,
      required: false,
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const dragItemHandleClass = 'sort-column-drag-handler';
    const isSortColumnsEmpty = computed(() => !props.modelValue?.length);
    const { addItemIntoArray, removeItemFromArray, updateItemInArray } = useArrayModelField(props, emit);
    const add = () => addItemIntoArray(widgetSortColumnToForm());

    return {
      dragItemHandleClass,
      isSortColumnsEmpty,
      add,
      remove: removeItemFromArray,
      updateItemInArray,
    };
  },
};
</script>
