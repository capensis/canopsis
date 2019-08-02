<template lang="pug">
  search-field(v-model="searchingText", @submit="submit", @clear="clear")
</template>

<script>
import { getContextWidgetSearchByText } from '@/helpers/entities-search';

import searchMixin from '@/mixins/search';

import SearchField from '@/components/forms/fields/search-field.vue';

/**
 * Search component for the entities list
 *
 * @module context
 */
export default {
  components: { SearchField },
  mixins: [searchMixin],
  data() {
    return {
      searchingText: '',
      requestParam: 'searchFilter',
    };
  },
  computed: {
    requestData() {
      return getContextWidgetSearchByText(this.searchingText);
    },
  },
};
</script>
