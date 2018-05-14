<template lang="pug">
  v-toolbar.toolbar(dense)
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
      span Aide sur la recherche avancée
</template>

<script>
import omit from 'lodash/omit';

export default {
  name: 'alarm-list-searching',
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
      this.$router.push({
        query: {
          ...this.$route.query,
          search,
        },

      });
    },
  },
};
</script>

<style scoped>
  .toolbar {
    background-color: rgb(251,247,247);
  }
</style>
