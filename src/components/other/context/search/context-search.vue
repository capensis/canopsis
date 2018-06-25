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
import searchMixin from '@/mixins/search';

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
