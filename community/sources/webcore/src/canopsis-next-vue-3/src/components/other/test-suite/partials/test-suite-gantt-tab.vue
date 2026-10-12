<template>
  <div class="position-relative">
    <c-progress-overlay :pending="pending" />
    <test-suite-historical-data-months-field
      :model-value="query.months"
      @update:model-value="updateMonths"
    />
    <junit-gantt-chart
      :items="ganttIntervals"
      :historical="historical"
      :total-items="meta.total_count"
      :query="query"
      :width="840"
      @update:query="query = $event"
    />
  </div>
</template>

<script>
import { mapActions } from 'pinia';

import { PAGINATION_LIMIT } from '@/config';
import { TEST_SUITE_STATUSES } from '@/constants';

import { useTestSuiteStore } from '@/stores/modules/test-suite';

import JunitGanttChart from '@/components/common/chart/junit-gantt-chart.vue';

import TestSuiteHistoricalDataMonthsField from './test-suite-historical-data-months-field.vue';

export default {
  components: { JunitGanttChart, TestSuiteHistoricalDataMonthsField },
  props: {
    testSuite: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      pending: false,
      historical: false,
      ganttIntervals: [],
      meta: {},
      query: {
        page: 1,
        itemsPerPage: PAGINATION_LIMIT,
        months: 0,
      },
    };
  },
  watch: {
    query: {
      deep: true,
      handler() {
        this.fetchList();
      },
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    ...mapActions(useTestSuiteStore, {
      fetchTestSuiteItemGanttIntervalsWithoutStore: 'fetchItemGanttIntervalsWithoutStore',
    }),

    updateMonths(months) {
      this.query = {
        ...this.query,
        months,
        page: 1,
      };
    },

    async fetchList() {
      try {
        this.pending = true;

        const { page, itemsPerPage: limit, months } = this.query;
        const { data = [], meta } = await this.fetchTestSuiteItemGanttIntervalsWithoutStore({
          id: this.testSuite._id,
          params: { page, limit, months },
        });

        const totalTimeInterval = {
          name: this.$t(`testSuite.statuses.${TEST_SUITE_STATUSES.total}`),
          status: TEST_SUITE_STATUSES.total,
          avg_status: TEST_SUITE_STATUSES.total,
          from: 0,
          to: meta.time,
          time: meta.time,
          avg_time: meta.avg_time,
          avg_to: meta.avg_time,
        };

        this.ganttIntervals = [...data, totalTimeInterval];
        this.meta = meta;
        this.historical = Boolean(months);
      } catch (err) {
        console.error(err);

        this.$popups.error({ text: err.message || err.description || this.$t('errors.default') });
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>
