<template lang="pug">
  v-select.select.pa-0(
  :items="items",
  v-model="limit",
  hide-details,
  single-line,
  dense,
  )
</template>

<script>
import { PAGINATION_LIMIT, PAGINATION_PER_PAGE_VALUES } from '@/config';

/**
 * Component to select number of items per page on lists
 *
 * @prop {Object} query - Object containing widget query information
 *
 * @event query#update
 */
export default {
  props: {
    query: {
      type: Object,
      required: true,
    },
  },
  data: () => ({
    items: PAGINATION_PER_PAGE_VALUES,
  }),
  computed: {
    limit: {
      get() {
        console.log(this.query);
        return this.query.limit || PAGINATION_LIMIT;
      },
      set(limit) {
        this.$emit('update:query', { ...this.query, page: 1, limit });
      },
    },
  },
};
</script>

<style scoped>
  .select {
    max-width: 65px;
  }
</style>
