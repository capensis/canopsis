<template>
  <c-select
    :label="$t('common.filterByType')"
    :model-value="value"
    :items="types"
    :name="name"
    clearable
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script>
import { EVENT_FILTER_FAILURE_TYPES } from '@/constants';

export default {
  emits: ['update:modelValue'],
  props: {
    value: {
      type: Number,
      required: false,
    },
    name: {
      type: String,
      default: 'type',
    },
  },
  computed: {
    types() {
      return Object.values(EVENT_FILTER_FAILURE_TYPES).map(value => ({
        value,
        text: this.$t(`eventFilter.failureTypes.${value}`),
      }));
    },
  },
};
</script>
