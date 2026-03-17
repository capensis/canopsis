<template>
  <c-advanced-search-field
    v-model="rules"
    :searches="searches"
    :attributes="attributes"
    :pending="pending"
    :allow-or="allowOr"
    with-history
    alarm-pattern
    v-on="$listeners"
  />
</template>

<script>
import { ref } from 'vue';

import { advancedSearchRuleItemToFormItem } from '@/helpers/search/advanced-search';

import { useAlarmAdvancedSearchAttributes } from './hooks/advanced-search';

export default {
  props: {
    searches: {
      type: Array,
      default: () => [],
    },
  },
  setup() {
    const rules = ref([advancedSearchRuleItemToFormItem()]);

    const { attributes, pending, allowOr } = useAlarmAdvancedSearchAttributes({ rules });

    return {
      rules,
      attributes,
      pending,
      allowOr,
    };
  },
};
</script>
