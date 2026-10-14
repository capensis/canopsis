<template>
  <c-select
    :model-value="value"
    :label="$t('common.separator')"
    :items="separators"
    :disabled="disabled"
    @update:model-value="updateModel"
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
    disabled: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['input', 'update:value'],

  setup(props, { emit }) {
    const separators = computed(() => Object.values(CSV_SEPARATORS));

    const updateModel = (value) => {
      emit('update:value', value);
      emit('input', value);
    };

    return {
      separators,
      updateModel,
    };
  },
};
</script>
