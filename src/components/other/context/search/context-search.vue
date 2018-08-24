<template lang="pug">
  v-toolbar.toolbar(dense, flat)
    v-text-field(
    label="Search",
    v-model="searchingText",
    @keyup.enter="submit",
    hide-details,
    single-line,
    )
    v-btn(icon @click="submit")
      v-icon search
    v-btn(icon @click="clear")
      v-icon clear
</template>

<script>
import searchMixin from '@/mixins/search';

/**
 * Search component for the entities list
 *
 * @module context
 */
export default {
  mixins: [searchMixin],
  data() {
    return {
      searchingText: '',
      requestParam: '_filter',
    };
  },
  computed: {
    requestData() {
      return this.$options.filters.formatContextSearch(this.searchingText);
    },
  },
};
</script>

<style scoped>
  .toolbar {
    background-color: white;
  }
</style>
