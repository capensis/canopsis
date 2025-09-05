<template>
  <c-lazy-search-field
    :value="selectedItems?.[0]"
    :items="items"
    :label="$tc('templateTesting.testName')"
    :loading="wholePending"
    :name="name"
    :has-more="hasMoreItems"
    item-text="_id"
    item-value="name"
    required
    @input="changeSelectedItems"
    @fetch="fetchItems"
    @fetch:more="fetchMoreItems"
    @update:search="updateSearch"
  />
</template>

<script>
import { toRef } from 'vue';

import { useLazySearch } from '@/hooks/form/lazy-search';
import { useTemplateTest } from '@/hooks/store/modules/template-test';

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
    name: {
      type: String,
      default: 'test',
    },
  },
  setup(props, { emit }) {
    const { fetchTemplateTestListWithoutStore } = useTemplateTest();

    const {
      selectedItems,
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
      fetchHandler: fetchTemplateTestListWithoutStore,
      idParamsKey: 'ids',
      idKey: '_id',
      addable: true,
    }, emit);

    return {
      selectedItems,
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
