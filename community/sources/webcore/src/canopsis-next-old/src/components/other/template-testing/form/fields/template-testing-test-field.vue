<template>
  <c-lazy-search-field
    :value="selectedItems?.[0]"
    :items="items"
    :label="$tc('templateTesting.testName')"
    :loading="wholePending"
    :name="name"
    :has-more="hasMoreItems"
    :required="required"
    :return-object="returnObject"
    item-text="name"
    item-value="_id"
    autocomplete
    @input="changeSelectedItems"
    @fetch="fetchItems"
    @fetch:more="fetchMoreItems"
    @update:search="updateSearch"
  />
</template>

<script>
import { merge } from 'lodash';
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
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'name',
    },
    ruleId: {
      type: String,
      required: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    returnObject: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { fetchTemplateTestListWithoutStore } = useTemplateTest();

    const fetchHandler = ({ params }) => fetchTemplateTestListWithoutStore({
      params: merge(params, props.params, { rule: props.ruleId }),
    });

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
      fetchHandler,
      value: toRef(props, 'value'),
      returnObject: toRef(props, 'returnObject'),
      idParamsKey: 'ids',
      idKey: '_id',
      attachValue: true,
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
