<template lang="pug">
  v-select.select.pa-0(
  :items="items",
  v-model="rowsPerPage",
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
    rowsPerPage: {
      get() {
        return this.query.rowsPerPage || PAGINATION_LIMIT;
      },
      set(rowsPerPage) {
        this.$emit('update:query', { ...this.query, page: 1, rowsPerPage });
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
