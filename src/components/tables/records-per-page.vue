<template lang="pug">
  v-select.select.pa-0(
  :items="$config.PAGINATION_PER_PAGE_VALUES",
  v-model="limit",
  hide-details,
  single-line,
  dense,
  )
</template>

<script>
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
  computed: {
    limit: {
      get() {
        return this.query.limit || this.$config.PAGINATION_LIMIT;
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
