<template>
  <c-select-field
    v-field="value"
    :label="$t('common.separator')"
    :items="separators"
  />
</template>

<script>
import { computed } from 'vue';

import { CSV_SEPARATORS } from '@/constants';

export default {
  props: {
    value: {
      type: String,
      default: CSV_SEPARATORS.comma,
    },
  },
  setup() {
    const separators = computed(() => Object.values(CSV_SEPARATORS));

    return {
      separators,
    };
  },
};
</script>
