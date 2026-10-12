<template>
  <c-layout class="position-relative" justify="center">
    <c-progress-overlay :pending="pending" />
    <junit-gantt-chart
      :items="ganttIntervals"
      :total-items="totalItems"
      :query="query"
      :width="840"
      @update:query="updateQuery"
    />
  </c-layout>
</template>

<script>
import { PAGINATION_LIMIT } from '@/config';

import { useTestSuite } from '@/hooks/store/modules/test-suite';

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

  computed: {
    entityId() {
      return this.alarm.entity?._id;
    },

    totalItems() {
      return this.meta.total_count ?? 0;
    },
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
    fetchEntityGanttIntervalsWithoutStore(params) {
      const { fetchTestSuiteEntityGanttIntervalsWithoutStore } = useTestSuite();

      return fetchTestSuiteEntityGanttIntervalsWithoutStore(params);
    },

    updateQuery(query) {
      this.query = query;
    },

    async fetchList() {
      if (!this.entityId) {
        this.ganttIntervals = [];
        this.meta = {};

        return;
      }

      try {
        this.pending = true;

        const { page, itemsPerPage: limit } = this.query;
        const { data, meta } = await this.fetchEntityGanttIntervalsWithoutStore({
          id: this.entityId,
          params: { page, limit },
        });

        this.ganttIntervals = data ?? [];
        this.meta = meta ?? {};
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
