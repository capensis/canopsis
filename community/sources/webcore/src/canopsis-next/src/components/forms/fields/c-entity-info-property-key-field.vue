<template>
  <c-lazy-search-field
    :value="selectedItem"
    :items="items"
    :label="label || $t('externalData.tableField')"
    :disabled="disabled"
    :required="required"
    :loading="wholePending"
    :has-more="hasMoreItems"
    :name="name"
    item-text="value"
    item-value="value"
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
import { useService } from '@/hooks/store/modules/service';

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
    name: {
      type: String,
      default: 'name',
    },
  },
  setup(props, { emit }) {
    const { fetchEntityInfosKeysWithoutStore } = useService();

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
      idKey: 'value',
      idParamsKey: 'ids',
      fetchHandler: fetchEntityInfosKeysWithoutStore,
      addable: true,
    }, emit);

    return {
      selectedItem,
      items,
      wholePending,
      hasMoreItems,
      fetchItems,
      fetchMoreItems,
      changeSelectedItems,
      removeItemFromSelectedItemsByIndex,
      updateSearch,
    };
  },
};
</script>
