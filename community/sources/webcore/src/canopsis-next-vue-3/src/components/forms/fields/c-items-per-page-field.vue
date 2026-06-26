<template>
  <c-select
    :items="items"
    :model-value="modelValue"
    class="select"
    name="itemsPerPage"
    density="compact"
    hide-details
    single-line
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script>
import { PAGINATION_LIMIT, PAGINATION_PER_PAGE_VALUES } from '@/config';

/**
 * Component to select number of items per page on lists
 *
 * @prop {number} [modelValue=PAGINATION_LIMIT] - Current value
 * @prop {Array} [items=PAGINATION_PER_PAGE_VALUES] - Items for selection
 *
 * @event update:modelValue
 */
export default {
  props: {
    modelValue: {
      type: Number,
      default: PAGINATION_LIMIT,
    },
    items: {
      type: Array,
      default: () => [...PAGINATION_PER_PAGE_VALUES],
    },
  },

  emits: ['update:modelValue'],
};
</script>

<style scoped>
.select {
  max-width: 65px;
}
</style>
