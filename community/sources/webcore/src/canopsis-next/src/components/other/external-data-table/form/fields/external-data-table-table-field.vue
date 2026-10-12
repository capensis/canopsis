<template>
  <c-lazy-search-field
    :value="selectedItem"
    :items="items"
    :label="label || $t('externalData.tableField')"
    :disabled="disabled"
    :required="required"
    :loading="wholePending"
    :has-more="hasMoreItems"
    item-text="name"
    item-value="_id"
    return-object
    @input="changeSelectedItems"
    @fetch="fetchItems"
    @fetch:more="fetchMoreItems"
    @update:search="updateSearch"
  />
</template>
<script>
import { toRef } from 'vue';

import { useLazySearch } from '@/hooks/form/lazy-search';
import { useExternalDataTable } from '@/hooks/store/modules/external-data-table';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { fetchExternalDataTablesListWithoutStore } = useExternalDataTable();

    const {
      selectedItem,
      items,
      wholePending,
      hasMoreItems,
      fetchItems,
      fetchMoreItems,
      changeSelectedItems,
      removeItemFromSelectedItemsByIndex,
      updateSearch,
    } = useLazySearch({
      value: toRef(props, 'value'),
      idKey: '_id',
      idParamsKey: 'ids',
      fetchHandler: fetchExternalDataTablesListWithoutStore,
    }, emit);

    return {
      selectedItem,
      items,
      wholePending,
      hasMoreItems,
      fetchItems,
      fetchMoreItems,
      changeSelectedItems,
      updateSearch,
      removeItemFromSelectedItemsByIndex,
    };
  },
};
</script>
