<template>
  <c-lazy-search-field
    v-bind="$attrs"
    :value="multiple ? selectedItems : selectedItems[0]"
    :label="label || $t('alarm.metaAlarm')"
    :loading="wholePending"
    :items="items"
    :name="name"
    :has-more="hasMoreItems"
    :required="required"
    :item-text="itemText"
    :item-value="itemValue"
    :disabled="disabled"
    :autocomplete="autocomplete"
    :clearable="clearable"
    :multiple="multiple"
    with-type
    return-object
    @input="changeSelectedItems"
    @fetch="fetchItems"
    @fetch:more="fetchMoreItems"
    @update:search="updateSearch"
  >
    <template #no-data="">
      <slot name="no-data" />
    </template>
  </c-lazy-search-field>
</template>

<script>
import { toRef } from 'vue';

import { useLazySearch } from '@/hooks/form/lazy-search';
import { useMetaAlarm } from '@/hooks/store/modules/meta-alarm';

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
      default: 'meta_alarm',
    },
    itemText: {
      type: String,
      default: 'name',
    },
    itemValue: {
      type: String,
      default: '_id',
    },
    label: {
      type: String,
      default: '',
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
    autocomplete: {
      type: Boolean,
      default: true,
    },
    addable: {
      type: Boolean,
      default: false,
    },
    multiple: {
      type: Boolean,
      default: false,
    },
    clearable: {
      type: Boolean,
      default: false,
    },
    returnObject: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { fetchMetaAlarmsListWithoutStore } = useMetaAlarm();

    const {
      selectedItems,
      items,
      wholePending,
      hasMoreItems,
      fetchItems,
      fetchMoreItems,
      changeSelectedItems,
      updateSearch,
    } = useLazySearch({
      idParamsKey: 'ids',
      fetchHandler: async (fetchParams) => {
        const data = await fetchMetaAlarmsListWithoutStore(fetchParams);

        return {
          data,
          meta: {
            total_count: data.length,
            page_count: 1,
          },
        };
      },
      idKey: toRef(props, 'itemValue'),
      value: toRef(props, 'value'),
      addable: toRef(props, 'addable'),
      multiple: toRef(props, 'multiple'),
      returnObject: toRef(props, 'returnObject'),
      limit: 100,
    }, emit);

    return {
      selectedItems,
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
