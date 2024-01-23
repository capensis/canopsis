<template>
  <c-search-field
    :combobox="combobox"
    :items="items"
    @submit="submit"
    @clear="clear"
    @toggle-pin="togglePin"
    @remove="remove"
  >
    <v-tooltip
      v-if="tooltip"
      bottom
    >
      <template #activator="{ on }">
        <v-btn
          icon
          v-on="on"
        >
          <v-icon>help_outline</v-icon>
        </v-btn>
      </template>
      <div v-html="tooltip" />
    </v-tooltip>
  </c-search-field>
</template>

<script>
import { omit } from 'lodash';

import { replaceTextNotInQuotes } from '@/helpers/search/quotes';

export default {
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
    combobox: {
      type: Boolean,
      default: false,
    },
    items: {
      type: Array,
      required: false,
    },
  },
  methods: {
    prepareRequestData(search = '') {
      if (!search.startsWith('-')) {
        return search;
      }

      const preparedSearch = search.replace(/^-(\s*)/, '');

      if (this.columns.length) {
        return this.columns.reduce(
          (acc, { text, value }) => replaceTextNotInQuotes(acc, text, value),
          preparedSearch,
        );
      }

      return preparedSearch;
    },

    remove(search) {
      this.$emit('remove', search);
    },

    togglePin(search) {
      this.$emit('toggle-pin', search);
    },

    clear() {
      const newQuery = omit(this.query, [this.field]);

      newQuery.page = 1;

      this.$emit('update:query', newQuery);
    },

    submit(search) {
      const requestData = this.prepareRequestData(search);

      this.$emit('submit', search);

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
