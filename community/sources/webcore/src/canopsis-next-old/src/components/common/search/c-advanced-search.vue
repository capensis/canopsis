<template>
  <c-advanced-search-field
    v-model="rules"
    :searches="searches"
    :attributes="attributes"
    :pending="pending"
    :allow-or="allowOr"
    :alarm-pattern="alarmPattern"
    :with-history="withHistory"
    v-on="$listeners"
  />
</template>

<script>
import { ref } from 'vue';

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
    alarmPattern: {
      type: Boolean,
      default: false,
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
