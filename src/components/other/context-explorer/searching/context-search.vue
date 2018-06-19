<template lang="pug">
  v-toolbar.toolbar(dense, flat)
    v-text-field(
      label="Search"
      v-model="searchingText"
      hide-details
      single-line
      @keyup.enter="submit"
      @keyup.delete="clear"
    )
    v-btn(icon @click="submit")
      v-icon search
    v-btn(icon @click="clear")
      v-icon clear
</template>

<script>
import omit from 'lodash/omit';
import merge from 'lodash/merge';

export default {
  filters: {
    formatedSearching(text) {
      return `{"$and":[{},{"$or":[{"name":{"$regex":"${text}","$options":"i"}},
      {"type":{"$regex":"${text}","$options":"i"}}]},{}]}`;
    },
  },
  data() {
    return {
      searchingText: '',
    };
  },
  methods: {
    clear() {
      const query = omit(this.$route.query, ['_filter']);
      this.$router.push({ query });
    },
    submit() {
      const filter = this.$options.filters.formatedSearching(this.searchingText);
      let { query } = this.$route;
      query = merge(query, { _filter: filter });
      this.$router.push({ query });
    },
  },
};
</script>

<style scoped>
  .toolbar {
    background-color: white;
  }
</style>
