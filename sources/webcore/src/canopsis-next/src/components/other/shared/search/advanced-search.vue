<template lang="pug">
  search-field(v-model="searchingText", @submit="submit", @clear="clear")
    v-tooltip(bottom)
      v-btn(
        data-test="tableSearchHelp",
        icon,
        slot="activator"
      )
        v-icon help_outline
      div(data-test="tableSearchHelpInfo", v-html="$t('search.advancedSearch')")
</template>

<script>
import { omit } from 'lodash';

import { replaceTextNotInQuotes } from '@/helpers/searching';

import SearchField from '@/components/forms/fields/search-field.vue';

/**
   * Search component for the entities list
   *
   * @module context
   */
export default {
  components: { SearchField },
  props: {
    query: {
      type: Object,
      required: true,
    },
    columns: {
      type: Array,
      default: () => [],
    },
    parameter: {
      type: String,
      default: 'search',
    },
  },
  data() {
    return {
      searchingText: '',
    };
  },
  methods: {
    getRequestData() {
      if (this.searchingText.startsWith('-')) {
        return this.columns.reduce((acc, { text, value }) =>
          replaceTextNotInQuotes(acc, text, value), this.searchingText.replace(/^-(\s*)/, ''));
      }

      return this.searchingText;
    },

    clear() {
      this.searchingText = '';

      this.$emit('update:query', omit(this.query, [this.parameter]));
    },

    submit() {
      const requestData = this.getRequestData();

      if (requestData || this.query[this.parameter]) {
        this.$emit('update:query', {
          ...this.query,

          page: 1,
          [this.parameter]: requestData,
        });
      }
    },
  },
};
</script>
