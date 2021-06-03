<template lang="pug">
  c-advanced-data-table(
    :items="detailsItems",
    :headers="headers",
    :loading="pending",
    :pagination.sync="pagination",
    :total-items="totalItems",
    advanced-pagination,
    expand
  )
    template(slot="status", slot-scope="props")
      v-tooltip(:disabled="!props.item.message", bottom)
        c-test-suite-chip(slot="activator", :value="props.item.status")
        span {{ props.item.message }}
    template(slot="time", slot-scope="props")
      span {{ props.item.time | fixed }}{{ $constants.TIME_UNITS.second }}
    template(slot="expand", slot-scope="props")
      test-suite-details-expand-panel(:test-suite-detail="props.item")
</template>

<script>
import { localQueryMixin } from '@/mixins/query-local/query';
import { entitiesTestSuitesMixin } from '@/mixins/entities/test-suite';

import TestSuiteDetailsExpandPanel from './test-suite-details-expand-panel.vue';

export default {
  components: { TestSuiteDetailsExpandPanel },
  mixins: [localQueryMixin, entitiesTestSuitesMixin],
  props: {
    testSuite: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      pending: false,
      detailsItems: [],
      totalItems: 0,
    };
  },
  computed: {
    headers() {
      return [
        { text: this.$t('common.name'), value: 'name' },
        { text: this.$t('common.status'), value: 'status' },
        { text: this.$t('testSuite.timeTaken'), value: 'time' },
      ];
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    async fetchList() {
      this.pending = true;

      const { data: details, meta } = await this.fetchTestSuiteDetailsWithoutStore({
        id: this.testSuite._id,
        params: this.getQuery(),
      });

      this.detailsItems = details;
      this.totalItems = meta.total_count;
      this.pending = false;
    },
  },
};
</script>