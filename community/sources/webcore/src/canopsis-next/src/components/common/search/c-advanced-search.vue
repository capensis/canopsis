<template>
  <c-advanced-search-field
    v-model:rules="rules"
    :searches="searches"
    :attributes="attributes"
    :pending="pending"
    :allow-or="allowOr"
    :alarm-pattern="alarmPattern"
    :with-history="withHistory"
    :name="name"
    :disabled="disabled"
    @submit="$emit('submit', $event)"
    @reset="$emit('reset')"
    @remove:search="$emit('remove:search', $event)"
    @toggle-pin:search="$emit('toggle-pin:search', $event)"
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
    name: {
      type: String,
      default: 'advancedSearch',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['submit', 'reset', 'remove:search', 'toggle-pin:search'],

  setup() {
    const rules = ref([advancedSearchRuleItemToFormItem()]);

    return {
      rules,
    };
  },
};
</script>
