<template>
  <c-lazy-search-field
    :value="fieldValue"
    v-bind="$attrs"
    :items="items"
    :label="label || t('common.pbehavior', multiple ? 2 : 1)"
    :loading="wholePending"
    :has-more="hasMoreItems"
    :name="name"
    :required="required"
    :disabled="disabled"
    :clearable="clearable"
    :item-title="itemTitle"
    :item-value="itemValue"
    autocomplete
    return-object
    @update:value="changeSelectedItems"
    @fetch="fetchItems"
    @fetch:more="fetchMoreItems"
    @update:search="updateSearch"
  >
    <template v-if="$slots.selection" #selection="slotProps">
      <slot name="selection" v-bind="slotProps" />
    </template>
    <template v-if="$slots.item" #item="slotProps">
      <slot name="item" v-bind="slotProps" />
    </template>
    <template v-if="$slots.icon" #icon="slotProps">
      <slot name="icon" v-bind="slotProps" />
    </template>
    <template v-if="$slots['no-data']" #no-data="">
      <slot name="no-data" />
    </template>
  </c-lazy-search-field>
</template>

<script setup>
import { isArray } from 'lodash';
import { computed, toRef } from 'vue';

import { useLazySearch } from '@/hooks/form/lazy-search';
import { useI18n } from '@/hooks/i18n';
import { usePbehavior } from '@/hooks/store/modules/pbehavior';

defineOptions({ inheritAttrs: false });

const props = defineProps({
  value: {
    type: [Array, String, Object],
    default: undefined,
  },
  modelValue: {
    type: [Array, String, Object],
    default: undefined,
  },
  label: {
    type: String,
    default: '',
  },
  name: {
    type: String,
    default: 'pbehavior',
  },
  itemTitle: {
    type: [String, Function],
    default: 'name',
  },
  itemValue: {
    type: String,
    default: '_id',
  },
  required: {
    type: Boolean,
    default: false,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  returnObject: {
    type: Boolean,
    default: false,
  },
  clearable: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(['update:value', 'update:modelValue', 'update:selected-items']);

const { t } = useI18n();
const { fetchPbehaviorsListWithoutStore } = usePbehavior();

const model = computed(() => props.modelValue ?? props.value ?? '');
const multiple = computed(() => isArray(model.value));

const emitModel = (event, value) => {
  emit(event, value);

  if (event === 'update:value') {
    emit('update:modelValue', value);
  }
};

const {
  selectedItem,
  selectedItems,
  items,
  wholePending,
  hasMoreItems,
  fetchItems,
  fetchMoreItems,
  changeSelectedItems,
  updateSearch,
} = useLazySearch(
  {
    value: model,
    idKey: toRef(props, 'itemValue'),
    idParamsKey: 'ids',
    fetchHandler: fetchPbehaviorsListWithoutStore,
    multiple,
    returnObject: toRef(props, 'returnObject'),
    attachValue: true,
  },
  emitModel,
);

const fieldValue = computed(() => (multiple.value ? selectedItems.value : selectedItem.value));
</script>
