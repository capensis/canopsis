<template>
  <c-advanced-search
    v-model="rules"
    :searches="searches"
    :attributes="attributes"
    :allow-or="allowOr"
    with-history
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

    const { attributes, allowOr } = useAlarmAdvancedSearchAttributes({ rules });

    return {
      rules,
      attributes,
      allowOr,
    };
  },
};
</script>
