<template lang="pug">
  v-select.select.pa-0(
    :items="items",
    v-model="itemPerPage",
    single-line,
    dense,
    hide-details
  )
</template>

<script>
import { PAGINATION_LIMIT } from '@/config';

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
