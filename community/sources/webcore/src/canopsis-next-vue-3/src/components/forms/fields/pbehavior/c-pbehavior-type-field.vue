<template>
  <c-lazy-search-field
    :value="selectedValue"
    :items="items"
    :label="label || t('common.type')"
    :loading="wholePending"
    :has-more="hasMoreItems"
    :name="name"
    :required="required"
    :disabled="disabled"
    :clearable="clearable"
    item-text="name"
    item-value="_id"
    autocomplete
    return-object
    @update:value="changeSelectedItems"
    @fetch="fetchItems"
    @fetch:more="fetchMoreItems"
    @update:search="updateSearch"
  />
</template>

<script>
import { computed, toRef, watch } from 'vue';

import { useLazySearch } from '@/hooks/form/lazy-search';
import { useI18n } from '@/hooks/i18n';
import { usePbehaviorType } from '@/hooks/store/modules/pbehavior-type';

export default {
  props: {
    modelValue: {
      type: [String, Object, Array],
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'type',
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
    multiple: {
      type: Boolean,
      default: false,
    },
    chips: {
      type: Boolean,
      default: false,
    },
    clearable: {
      type: Boolean,
      default: false,
    },
    types: {
      type: Array,
      required: false,
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { fetchPbehaviorTypesFieldList } = usePbehaviorType();

    const fetchHandler = ({ params }) =>
      fetchPbehaviorTypesFieldList({
        params: {
          ...params,
          with_hidden: true,
          ...(props.types ? { types: props.types } : {}),
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
        value: toRef(props, 'modelValue'),
        idKey: '_id',
        idParamsKey: 'ids',
        fetchHandler,
        multiple: toRef(props, 'multiple'),
        returnObject: toRef(props, 'returnObject'),
        attachValue: true,
      },
      emit,
    );

    const selectedValue = computed(() => (props.multiple ? selectedItems.value : selectedItem.value));

    watch(() => props.types, () => fetchItems(), { deep: true });

    return {
      t,
      selectedValue,
      items,
      wholePending,
      hasMoreItems,
      fetchItems,
      fetchMoreItems,
      changeSelectedItems,
      updateSearch,
    };
  },
};
</script>
