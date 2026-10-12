<template>
  <c-lazy-search-field
    :value="fieldValue"
    v-bind="$attrs"
    :label="selectLabel"
    :loading="wholePending"
    :items="items"
    :name="name"
    :has-more="hasMoreItems"
    :item-text="itemText"
    :item-value="itemValue"
    :item-disabled="itemDisabled"
    :disabled="disabled"
    :required="required"
    :clearable="clearable"
    :autocomplete="autocomplete"
    return-object
    with-type
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
  </c-lazy-search-field>
</template>

<script setup>
import { isArray } from 'lodash';
import { computed, toRef, watch } from 'vue';

import { PAGINATION_LIMIT } from '@/config';
import { BASIC_ENTITY_TYPES } from '@/constants';

import { useLazySearch } from '@/hooks/form/lazy-search';
import { useI18n } from '@/hooks/i18n';
import { useEntity } from '@/hooks/store/modules/entity';

defineOptions({ inheritAttrs: false });

const props = defineProps({
  value: {
    type: [Array, String, Object],
    default: '',
  },
  name: {
    type: String,
    default: 'entities',
  },
  label: {
    type: String,
    required: false,
  },
  itemText: {
    type: [String, Function],
    default: '_id',
  },
  itemValue: {
    type: String,
    default: '_id',
  },
  limit: {
    type: Number,
    default: PAGINATION_LIMIT,
  },
  entityTypes: {
    type: Array,
    default: () => Object.values(BASIC_ENTITY_TYPES),
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
  autocomplete: {
    type: Boolean,
    default: false,
  },
  itemDisabled: {
    type: [String, Array, Function],
    required: false,
  },
  required: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(['input', 'update:value', 'update:modelValue', 'update:selected-items']);

const { t } = useI18n();
const { fetchContextEntitiesListWithoutStore } = useEntity();

const multiple = computed(() => isArray(props.value));
const selectLabel = computed(() => {
  if (props.label) {
    return props.label;
  }

  if (multiple.value) {
    return t('common.entity', props.value.length);
  }

  return t('common.entity');
});

const emitModel = (event, value) => {
  emit(event, value);

  if (event === 'update:value') {
    emit('input', value);
    emit('update:modelValue', value);
  }
};

const fetchHandler = ({ params }) => fetchContextEntitiesListWithoutStore({
    params: {
      ...params,
      type: props.entityTypes,
    },
  });

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
    fetchHandler,
    idKey: toRef(props, 'itemValue'),
    idParamsKey: '_id',
    value: toRef(props, 'value'),
    multiple,
    addable: computed(() => !props.autocomplete),
    returnObject: toRef(props, 'returnObject'),
    attachValue: true,
    initialQuery: {
      page: 1,
      limit: props.limit,
      search: '',
    },
  },
  emitModel,
);

const fieldValue = computed(() => (multiple.value ? selectedItems.value : selectedItem.value));

watch(
  () => props.entityTypes,
  () => {
    changeSelectedItems(multiple.value ? [] : '');
    fetchItems();
  },
  { deep: true },
);
</script>
