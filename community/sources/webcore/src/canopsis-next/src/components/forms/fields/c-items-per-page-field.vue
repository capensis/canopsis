<template>
  <v-select
    :items="items"
    :value="value"
    class="select"
    hide-details
    single-line
    dense
    @input="$emit('input', $event)"
  />
</template>

<script>
import { PAGINATION_LIMIT, PAGINATION_PER_PAGE_VALUES } from '@/config';

/**
 * Component to select number of items per page on lists
 *
 * @prop {number} [value=PAGINATION_LIMIT] - Current value
 * @prop {Array} [items=PAGINATION_PER_PAGE_VALUES] - Items for selection
 *
 * @event query#update
 */
export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Number,
      default: PAGINATION_LIMIT,
    },
    items: {
      type: Array,
      default: () => [...PAGINATION_PER_PAGE_VALUES],
    },
  },
};
</script>

<style scoped>
  .select {
    max-width: 65px;
  }
</style>
