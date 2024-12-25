<template>

</template>

<script>
import { keyBy, isArray } from 'lodash';
import { computed, toRef, ref, onMounted } from 'vue';

import { usePendingWithLocalQuery } from '@/hooks/query/shared';
import { useLazySearch } from '@/hooks/form/lazy-search';

export default {
  props: {
    value: {
      type: [String, Array],
      default: '',
    },
    fetch: {
      type: Function,
      required: true,
    },
  },
  setup(props, { emit }) {
    const itemsById = ref({});
    const pageCount = ref(1);

    const {
      selectedItems,
      items,
      hasMoreItems,
      fetchItems,
      fetchMoreItems,
      changeSelectedItems,
      removeItemFromSelectedItemsByIndex,
      updateSearch,
      wholePending: pending,
    } = useLazySearch({
      value: toRef(props, 'value'),
      addable: true,
      idKey: '_id',
      idParamsKey: 'ids',
      fetchHandler: props.fetch,
    }, emit);

    return {
      selectedItems,
      pending,
      items,
      hasMoreItems,

      fetchMoreItems,
      updateSearch,
    };
  },
};
</script>
