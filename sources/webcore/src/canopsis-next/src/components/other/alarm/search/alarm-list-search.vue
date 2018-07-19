<template lang="pug">
  v-toolbar.white(dense, flat)
    v-text-field(
      label="Search",
      v-model="searchingText",
      hide-details,
      single-line,
      @keyup.enter="submit",
      @keyup.delete="clear",
    )
    v-btn(icon @click="submit")
      v-icon search
    v-btn(icon @click="clear")
      v-icon clear
    v-tooltip(bottom)
      v-btn(icon slot="activator")
        v-icon help_outline
      div(v-html="$t('search.advancedSearch')",)
</template>

<script>
import searchMixin from '@/mixins/search';

/**
 * Search component for the alarms list
 *
 * @module alarm
 */
export default {
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

