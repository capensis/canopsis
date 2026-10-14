<template>
  <c-advanced-search-field
    v-model:rules="rules"
    :searches="searches"
    :attributes="attributes"
    :pending="pending"
    :allow-or="allowOr"
    with-history
    alarm-pattern
    @submit="$emit('submit', $event)"
    @reset="$emit('reset')"
    @remove:search="$emit('remove:search', $event)"
    @toggle-pin:search="$emit('toggle-pin:search', $event)"
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

  emits: ['submit', 'reset', 'remove:search', 'toggle-pin:search'],

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
