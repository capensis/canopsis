<template>
  <c-lazy-search-field
    v-field="value"
    :label="$t('metaAlarmRule.field.title')"
    :loading="pending"
    :items="rules"
    :name="name"
    :has-more="hasMore"
    :required="required"
    :item-text="itemText"
    :item-value="itemValue"
    :disabled="disabled"
    :no-data-text="$t('metaAlarmRule.field.noData')"
    clearable
    autocomplete
    @input="clearQuerySearch"
    @fetch="fetchRules"
    @fetch:more="fetchMoreRules"
    @update:search="updateQuerySearch"
  />
</template>

<script>
import { ref, computed, onMounted } from 'vue';

import { prepareDataForItemsById } from '@/helpers/search/lazy-search';

import { usePendingWithLocalQuery } from '@/hooks/query/shared';
import { useMetaAlarmRule } from '@/hooks/store/modules/meta-alarm-rule';

export default {
  props: {
    value: {
      type: [Array, String, Object],
      default: '',
    },
    name: {
      type: String,
      default: 'meta_alarm_rule',
    },
    itemText: {
      type: String,
      default: 'name',
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
  },
  setup(props) {
    const rulesById = ref({});
    const meta = ref({});

    const rules = computed(() => Object.values(rulesById.value));

    const {
      fetchMetaAlarmRulesListWithoutStore,
    } = useMetaAlarmRule();

    const {
      pending,
      query,
      fetchHandlerWithQuery: fetchRules,
      updateQueryPage,
      updateQuerySearch,
    } = usePendingWithLocalQuery({
      initialQuery: {
        page: 1,
        limit: props.limit,
        search: props.value,
      },
      fetchHandler: async (params) => {
        const response = await fetchMetaAlarmRulesListWithoutStore({ params });

        rulesById.value = prepareDataForItemsById(
          query.value.page !== 1 ? rulesById.value : {},
          props.value,
          response.data,
        );
        meta.value = response.meta;
      },
    });

    const hasMore = computed(() => meta.value.page_count > query.value.page);

    const fetchMoreRules = () => updateQueryPage(query.value.page + 1);

    const clearQuerySearch = value => !value && updateQuerySearch('');

    onMounted(fetchRules);

    return {
      rules,
      meta,
      pending,
      hasMore,

      fetchRules,
      fetchMoreRules,
      clearQuerySearch,
      updateQuerySearch,
    };
  },
};
</script>
