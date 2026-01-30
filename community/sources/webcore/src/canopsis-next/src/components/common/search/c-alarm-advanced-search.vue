<template>
  <c-advanced-search
    v-model="rules"
    ref="advancedSearchElement"
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
    const advancedSearchElement = ref(null);
    const rules = ref([advancedSearchRuleItemToFormItem()]);

    const { attributes, allowOr } = useAlarmAdvancedSearchAttributes({ rules });

    return {
      advancedSearchElement,
      rules,
      attributes,
      allowOr,
    };
  },
};
</script>
