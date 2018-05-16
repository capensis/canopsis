<template lang="pug">
transition(name="fade" mode="out-in")
  v-btn(v-if="sortingDirection === 'DESC'" flat icon @click.prevent="sortingDirection = 'ASC'")
    v-icon arrow_drop_down
  v-btn(v-if="sortingDirection === 'ASC'" flat icon @click.prevent="sortingDirection = 'DESC'")
    v-icon arrow_drop_up
</template>

<script>
export default {
  name: 'alarm-list-sorting',
  props: {
    columnToSort: {
      type: String,
      required: true,
    },
  },
  computed: {
    sortingDirection: {
      get() {
        return this.$route.query.sort_dir || 'DESC';
      },
      set(sortDir) {
        this.$router.push({
          query: {
            ...this.$route.query,
            sort_dir: sortDir,
            sort_key: this.columnToSort,
          },
        });
      },
    },
  },
};
</script>

<style scoped>

</style>
