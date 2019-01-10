<template lang="pug">
  search-field(v-model="searchingText", @submit="submit", @clear="clear")
    v-tooltip(bottom)
      v-btn(icon slot="activator")
        v-icon help_outline
      div(v-html="$t('search.advancedSearch')",)
</template>

<script>
import searchMixin from '@/mixins/search';

import SearchField from '@/components/forms/fields/search-field.vue';

/**
 * Search component for the alarms list
 *
 * @module alarm
 */
export default {
  components: { SearchField },
  mixins: [searchMixin],
  data() {
    return {
      searchingText: '',
      requestParam: 'search',
    };
  },
  computed: {
    requestData() {
      return this.searchingText;
    },
  },
};
</script>
