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
import { getContextSearchByText } from '@/helpers/widget-search';
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
      requestParam: 'searchFilter',
    };
  },
  computed: {
    requestData() {
      return getContextSearchByText(this.searchingText);
    },
  },
  methods: {
    submit() {
      this.$emit('update:query', {
        ...this.query,

        page: 1,
        [this.requestParam]: getContextSearchByText(this.searchingText),
      });
    },
  },
};
</script>

<style scoped>
  .toolbar {
    background-color: white;
  }
</style>
