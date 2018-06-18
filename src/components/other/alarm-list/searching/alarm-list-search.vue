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
    v-tooltip(bottom)
      v-btn(icon slot="activator")
        v-icon help_outline
      div(v-html="$t('search.advancedSearch')")
</template>

<script>
import omit from 'lodash/omit';

export default {
  data() {
    return {
      searchingText: '',
    };
  },
  methods: {
    clear() {
      const query = omit(this.$route.query, ['search']);
      this.$router.push({
        query: {
          ...query,
        },
      });
    },
    submit() {
      const search = this.searchingText;
      const { query } = this.$route;
      this.$router.push({ query, search });
    },
  },
};
</script>

<style scoped>
  .toolbar {
    background-color: white;
  }
</style>
