<template>
  <c-lazy-search-field
    :value="multiple ? selectedItems : selectedItem"
    :label="label || t('alarm.metaAlarm')"
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
    @update:value="changeSelectedItems"
    @fetch="fetchItems"
    @fetch:more="fetchMoreItems"
    @update:search="updateSearch"
  >
    <template #no-data>
      <slot name="no-data" />
    </template>
  </c-lazy-search-field>
</template>

<script>
import { toRef } from 'vue';

import { useLazySearch } from '@/hooks/form/lazy-search';
import { useI18n } from '@/hooks/i18n';
import { useMetaAlarm } from '@/hooks/store/modules/meta-alarm';

export default {
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
      default: 100,
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

  emits: ['update:value', 'update:modelValue'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { fetchMetaAlarmsListWithoutStore } = useMetaAlarm();

    const {
      selectedItems,
      selectedItem,
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

        if (Array.isArray(data)) {
          return {
            data,
            meta: {
              total_count: data.length,
              page_count: 1,
            },
          };
        }

        return data;
      },
      idKey: toRef(props, 'itemValue'),
      value: toRef(props, 'value'),
      addable: toRef(props, 'addable'),
      multiple: toRef(props, 'multiple'),
      returnObject: toRef(props, 'returnObject'),
      initialQuery: {
        page: 1,
        limit: props.limit,
        search: '',
      },
    }, emit);

    return {
      t,
      selectedItems,
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
