<template>
  <c-lazy-search-field
    :value="selectedItem"
    :items="items"
    :label="$t('alarm.alarmDisplayName')"
    :disabled="disabled"
    :required="required"
    :loading="wholePending"
    :has-more="hasMoreItems"
    :name="name"
    :item-text="itemText"
    :item-value="itemValue"
    :no-data-text="$t('alarm.noAlarmFound')"
    return-object
    clearable
    autocomplete
    with-type
    @input="changeSelectedItems"
    @fetch="fetchItems"
    @fetch:more="fetchMoreItems"
    @update:search="updateSearch"
  />
</template>

<script>
import { merge } from 'lodash';
import { toRef, watch } from 'vue';

import { useLazySearch } from '@/hooks/form/lazy-search';
import { useAlarmStoreModule } from '@/hooks/store/modules/alarm';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [Array, String, Object],
      default: '',
    },
    name: {
      type: String,
      default: 'alarm',
    },
    itemText: {
      type: String,
      default: 'display_name',
    },
    itemValue: {
      type: String,
      default: '_id',
    },
    limit: {
      type: Number,
      default: 20,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    params: {
      type: Object,
      default: () => {},
    },
  },
  setup(props, { emit }) {
    const { useActions } = useAlarmStoreModule();

    const { fetchDisplayNamesWithoutStore } = useActions({
      fetchDisplayNamesWithoutStore: 'fetchDisplayNamesWithoutStore',
    });

    const fetchHandler = ({ params }) => fetchDisplayNamesWithoutStore({
      params: merge(params, props.params),
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
    } = useLazySearch({
      fetchHandler,
      value: toRef(props, 'value'),
      idKey: props.itemValue,
      idParamsKey: 'ids',
      attachValue: true,
    }, emit);

    watch(() => props.params, () => {
      fetchItems();
    }, { deep: true });

    return {
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
