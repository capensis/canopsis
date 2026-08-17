<template>
  <c-layout class="position-relative" justify="center">
    <c-progress-overlay :pending="pending" />
    <junit-gantt-chart
      :items="ganttIntervals"
      :total-items="meta.total_count"
      :query="query"
      :width="840"
      @update:query="query = $event"
    />
  </c-layout>
</template>

<script>
import { mapActions } from 'pinia';

import { PAGINATION_LIMIT } from '@/config';

import { useTestSuiteStore } from '@/stores/modules/test-suite';

import JunitGanttChart from '@/components/common/chart/junit-gantt-chart.vue';

export default {
  components: { JunitGanttChart },

  props: {
    alarm: {
      type: Object,
      required: true,
    },
  },

  data() {
    return {
      pending: false,
      meta: {},
      ganttIntervals: [],
      query: {
        itemsPerPage: PAGINATION_LIMIT,
        page: 1,
      },
    };
  },

  watch: {
    query: {
      deep: true,
      immediate: true,
      handler() {
        this.fetchList();
      },
    },
  },

  methods: {
    ...mapActions(useTestSuiteStore, ['fetchEntityGanttIntervalsWithoutStore']),

    async fetchList() {
      try {
        this.pending = true;

        const { page, itemsPerPage: limit } = this.query;
        const { data, meta } = await this.fetchEntityGanttIntervalsWithoutStore({
          id: this.alarm.entity._id,
          params: { page, limit },
        });

        this.ganttIntervals = data;
        this.meta = meta;
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
