<template>
  <patterns-list
    :options.sync="options"
    :patterns="corporatePatterns"
    :total-items="corporatePatternsMeta.total_count"
    :pending="corporatePatternsPending"
    corporate
    @edit="$listeners.edit"
    @remove="$listeners.remove"
    @remove-selected="$listeners['remove-selected']"
  />
</template>

<script>
import { localQueryMixin } from '@/mixins/query/query';
import { entitiesCorporatePatternsMixin } from '@/mixins/entities/pattern/corporate';

import PatternsList from './patterns-list.vue';

export default {
  components: {
    PatternsList,
  },
  mixins: [
    localQueryMixin,
    entitiesCorporatePatternsMixin,
  ],
  mounted() {
    this.fetchList();
  },
  methods: {
    fetchList() {
      return this.fetchCorporatePatternsList({ params: this.getQuery() });
    },
  },
};
</script>
