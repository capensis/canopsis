<template lang="pug">
  v-select.select(
    :items="items",
    v-model="itemPerPage",
    single-line,
    dense,
    hide-details,
    class="pa-0"
  )
</template>

<script>
import { PAGINATION_LIMIT } from '@/config';

/**
 * Component to select number of items per page on lists
 */
export default {
  props: {
    query: {
      type: Object,
      required: true,
    },
  },
  data: () => ({
    items: [5, 10, 20, 50],
  }),
  computed: {
    itemPerPage: {
      get() {
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
