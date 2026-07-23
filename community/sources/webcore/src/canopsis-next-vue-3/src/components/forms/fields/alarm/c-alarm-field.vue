<template>
  <c-lazy-search-field
    :value="selectedItem"
    :items="items"
    :label="t('alarm.alarmDisplayName')"
    :disabled="disabled"
    :required="required"
    :loading="wholePending"
    :has-more="hasMoreItems"
    :name="name"
    :item-text="itemText"
    :item-value="itemValue"
    :no-data-text="t('alarm.noAlarmFound')"
    return-object
    clearable
    autocomplete
    with-type
    @update:value="changeSelectedItems"
    @fetch="fetchItems"
    @fetch:more="fetchMoreItems"
    @update:search="updateSearch"
  />
</template>

<script>
import { merge } from 'lodash';
import { toRef, watch } from 'vue';

import { useI18n } from '@/hooks/i18n';
import { useLazySearch } from '@/hooks/form/lazy-search';
import { useAlarm } from '@/hooks/store/modules/alarm';

export default {
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
      default: () => ({}),
    },
  },

  emits: ['update:value'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { fetchListWithoutStore } = useAlarm();

    const fetchHandler = ({ params }) =>
      fetchListWithoutStore({
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
    } = useLazySearch(
      {
        fetchHandler,
        value: toRef(props, 'value'),
        idKey: props.itemValue,
        idParamsKey: 'ids',
        attachValue: true,
      },
      emit,
    );

    watch(
      () => props.params,
      () => fetchItems(),
      { deep: true },
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
