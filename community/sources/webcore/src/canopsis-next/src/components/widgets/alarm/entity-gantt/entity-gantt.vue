<template>
  <v-layout justify-center>
    <c-progress-overlay :pending="pending" />
    <junit-gantt-chart
      :items="ganttIntervals"
      :total-items="meta.total_count"
      :query.sync="query"
      :width="840"
    />
  </v-layout>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { PAGINATION_LIMIT } from '@/config';

import JunitGanttChart from '@/components/common/chart/junit-gantt-chart.vue';

const { mapActions } = createNamespacedHelpers('testSuite/entityGantt');

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
    ...mapActions({
      fetchEntityGanttIntervalsWithoutStore: 'fetchItemGanttIntervalsWithoutStore',
    }),

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
