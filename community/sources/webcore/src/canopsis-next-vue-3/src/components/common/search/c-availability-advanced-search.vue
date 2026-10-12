<template>
  <c-advanced-search
    :searches="searches"
    :attributes="attributes"
    :pending="pending"
    allow-or
    with-history
    @submit="$emit('submit', $event)"
    @reset="$emit('reset')"
    @remove:search="$emit('remove:search', $event)"
    @toggle-pin:search="$emit('toggle-pin:search', $event)"
  />
</template>

<script>
import { useAvailabilityAdvancedSearchAttributes } from './hooks/advanced-search';

export default {
  props: {
    searches: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['submit', 'reset', 'remove:search', 'toggle-pin:search'],

  setup() {
    const { attributes, pending } = useAvailabilityAdvancedSearchAttributes();

    return {
      attributes,
      pending,
    };
  },
};
</script>
