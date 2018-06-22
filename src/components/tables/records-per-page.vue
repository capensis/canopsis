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
// OTHERS
import { PAGINATION_LIMIT } from '@/config';

export default {
  data: () => ({
    items: [5, 10, 20, 50],
  }),
  computed: {
    itemPerPage: {
      get() {
        return parseInt(this.$route.query.limit, 10) || PAGINATION_LIMIT;
      },
      set(limit) {
        this.$router.push({
          query: {
            ...this.$route.query,
            limit,
            page: 1,
          },
        });
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
