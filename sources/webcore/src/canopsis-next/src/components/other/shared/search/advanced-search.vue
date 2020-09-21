<template lang="pug">
  search-field(v-model="searchingText", @submit="submit", @clear="clear")
    v-tooltip(v-if="tooltip", bottom)
      v-btn(
        data-test="tableSearchHelp",
        icon,
        slot="activator"
      )
        v-icon help_outline
      div(data-test="tableSearchHelpInfo", v-html="tooltip")
</template>

<script>
import { omit } from 'lodash';

import { replaceTextNotInQuotes } from '@/helpers/searching';

import SearchField from '@/components/forms/fields/search-field.vue';

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
    field: {
      type: String,
      default: 'search',
    },
    tooltip: {
      type: String,
      default: '',
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

      this.$emit('update:query', omit(this.query, [this.field]));
    },

    submit() {
      const requestData = this.getRequestData();

      if (requestData || this.query[this.field]) {
        this.$emit('update:query', {
          ...this.query,

          page: 1,
          [this.field]: requestData,
        });
      }
    },
  },
};
</script>
