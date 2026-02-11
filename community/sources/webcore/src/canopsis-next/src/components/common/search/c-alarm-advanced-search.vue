<template>
  <c-advanced-search-field
    v-model="rules"
    :searches="searches"
    :attributes="attributes"
    :allow-or="allowOr"
    :basic-field="basicField"
    with-history
    v-on="$listeners"
  />
</template>

<script>
import { ref } from 'vue';

import { ADVANCED_SEARCH_FIELDS } from '@/constants';

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
    const basicField = ADVANCED_SEARCH_FIELDS.alarm;

    const rules = ref([advancedSearchRuleItemToFormItem()]);

    const { attributes, allowOr } = useAlarmAdvancedSearchAttributes({ rules });

    return {
      basicField,
      rules,
      attributes,
      allowOr,
    };
  },
};
</script>
