<template>
  <c-card-iterator-form
    :value="modelValue"
    :name="name || undefined"
    :handle="`.${dragItemHandleClass}`"
    :iterator-class="{ empty: isColumnsEmpty }"
    @add="add"
    @update:value="$emit('update:modelValue', $event)"
  >
    <template #item="{ item: column, index }">
      <column-field
        :key="column.key"
        :model-value="modelValue[index]"
        :name="getColumnFieldName(column.key)"
        :type="type"
        :drag-handle-class="dragItemHandleClass"
        :with-html="withHtml"
        :with-template="withTemplate"
        :with-color-indicator="withColorIndicator"
        :with-instructions="withInstructions"
        :optional-infos-attributes="optionalInfosAttributes"
        :with-simple-template="withSimpleTemplate"
        :without-infos-attributes="withoutInfosAttributes"
        :without-custom-label="withoutCustomLabel"
        :with-filter-on-click="withFilterOnClick"
        :items="items"
        :variables="variables"
        :excluded-columns="excludedColumns"
        class="mb-3"
        @update:model-value="updateItemInArray(index, $event)"
        @remove="remove(index)"
      />
    </template>
  </c-card-iterator-form>
</template>

<script>
import { computed } from 'vue';

import { ENTITIES_TYPES } from '@/constants';

import { widgetColumnToForm } from '@/helpers/entities/widget/column/form';

import { useArrayModelField } from '@/hooks/form/array-model-field';

import ColumnField from './partials/column-field.vue';

export default {
  components: { ColumnField },

  props: {
    modelValue: {
      type: Array,
      default: () => [],
    },
    type: {
      type: String,
      default: ENTITIES_TYPES.alarm,
    },
    withTemplate: {
      type: Boolean,
      default: false,
    },
    withHtml: {
      type: Boolean,
      default: false,
    },
    withColorIndicator: {
      type: Boolean,
      default: false,
    },
    withInstructions: {
      type: Boolean,
      default: false,
    },
    withoutCustomLabel: {
      type: Boolean,
      default: false,
    },
    optionalInfosAttributes: {
      type: Boolean,
      default: false,
    },
    withSimpleTemplate: {
      type: Boolean,
      default: false,
    },
    variables: {
      type: Array,
      required: false,
    },
    withoutInfosAttributes: {
      type: Boolean,
      default: false,
    },
    withFilterOnClick: {
      type: Boolean,
      default: false,
    },
    items: {
      type: Array,
      required: false,
    },
    excludedColumns: {
      type: Array,
      required: false,
    },
    name: {
      type: String,
      default: '',
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const dragItemHandleClass = 'column-drag-handle';
    const isColumnsEmpty = computed(() => !props.modelValue.length);
    const { addItemIntoArray, removeItemFromArray, updateItemInArray } = useArrayModelField(
      props,
      emit,
    );
    const add = () => addItemIntoArray(widgetColumnToForm());
    const getColumnFieldName = key => [props.name, key].filter(Boolean).join('.');

    return {
      dragItemHandleClass,
      isColumnsEmpty,
      add,
      remove: removeItemFromArray,
      updateItemInArray,
      getColumnFieldName,
    };
  },
};
</script>
