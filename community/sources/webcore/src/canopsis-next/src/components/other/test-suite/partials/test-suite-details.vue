<template>
  <c-advanced-data-table
    :items="detailsItems"
    :headers="headers"
    :loading="pending"
    :options.sync="options"
    :total-items="totalItems"
    advanced-pagination
    expand
  >
    <template #status="{ item }">
      <v-tooltip
        :disabled="!item.message"
        bottom
      >
        <template #activator="{ on }">
          <c-test-suite-chip
            :value="item.status"
            v-on="on"
          />
        </template>
        <span>{{ item.message }}</span>
      </v-tooltip>
    </template>
    <template #time="{ item }">
      <span>{{ item.time | fixed }}{{ $constants.TIME_UNITS.second }}</span>
    </template>
    <template #expand="{ item }">
      <test-suite-details-expand-panel :test-suite-detail="item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { localQueryMixin } from '@/mixins/query/query';
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
        { text: this.$t('common.timeTaken'), value: 'time' },
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
