<template>
  <c-lazy-search-field
    :value="selectedItem"
    :items="items"
    :label="label || $t('externalAuthToken.tokenValue')"
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
import { useWebhookTokenRule } from '@/hooks/store/modules/webhook-token-rule';

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
    const { fetchWebhookTokenRulesListWithoutStore } = useWebhookTokenRule();

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
      fetchHandler: fetchWebhookTokenRulesListWithoutStore,
      attachValue: true,
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
