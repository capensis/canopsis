<template>
  <c-lazy-search-field
    :value="selectedItem"
    :items="items"
    :label="label || t('common.reason')"
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
import { toRef } from 'vue';

import { useLazySearch } from '@/hooks/form/lazy-search';
import { useI18n } from '@/hooks/i18n';
import { usePbehaviorReason } from '@/hooks/store/modules/pbehavior-reason';

export default {
  props: {
    modelValue: {
      type: [String, Object],
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'reason',
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
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { fetchPbehaviorReasonsListWithoutStore } = usePbehaviorReason();

    const fetchHandler = ({ params }) =>
      fetchPbehaviorReasonsListWithoutStore({
        params: {
          ...params,
          with_hidden: true,
        },
      });

    const {
      selectedItem,
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
        returnObject: toRef(props, 'returnObject'),
        attachValue: true,
      },
      emit,
    );

    return {
      t,
      selectedItem,
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
