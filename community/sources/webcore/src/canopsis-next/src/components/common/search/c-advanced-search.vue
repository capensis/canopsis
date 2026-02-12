<template>
  <c-advanced-search-field
    v-model="rules"
    :searches="searches"
    :attributes="attributes"
    :pending="pending"
    :allow-or="allowOr"
    :basic-field="basicField"
    :with-history="withHistory"
    v-on="$listeners"
  />
</template>

<script>
import { ref } from 'vue';

import { ADVANCED_SEARCH_FIELDS } from '@/constants';

import { advancedSearchRuleItemToFormItem } from '@/helpers/search/advanced-search';

export default {
  props: {
    attributes: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: false,
    },
    searches: {
      type: Array,
      default: () => [],
    },
    allowOr: {
      type: Boolean,
      default: true,
    },
    withHistory: {
      type: Boolean,
      default: false,
    },
    basicField: {
      type: String,
      default: ADVANCED_SEARCH_FIELDS.search,
    },
  },
  setup() {
    const rules = ref([advancedSearchRuleItemToFormItem()]);

    return {
      rules,
    };
  },
};
</script>
