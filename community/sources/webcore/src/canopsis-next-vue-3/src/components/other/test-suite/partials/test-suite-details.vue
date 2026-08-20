<template>
  <c-advanced-data-table
    v-model:options="options"
    :items="detailsItems"
    :headers="headers"
    :loading="pending"
    :total-items="totalItems"
    advanced-pagination
    expand
  >
    <template #status="{ item }">
      <v-tooltip
        :disabled="!item.message"
        location="bottom"
      >
        <template #activator="{ props: activatorProps }">
          <c-test-suite-chip
            :value="item.status"
            v-bind="activatorProps"
          />
        </template>
        <span>{{ item.message }}</span>
      </v-tooltip>
    </template>
    <template #time="{ item }">
      <span>{{ convertNumberToFixedString(item.time) }}{{ $constants.TIME_UNITS.second }}</span>
    </template>
    <template #expand="{ item }">
      <test-suite-details-expand-panel :test-suite-detail="item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { convertNumberToFixedString } from '@/helpers/string';

import { entitiesTestSuitesMixin } from '@/mixins/entities/test-suite';
import { localQueryMixin } from '@/mixins/query/query';

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
    convertNumberToFixedString,

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
