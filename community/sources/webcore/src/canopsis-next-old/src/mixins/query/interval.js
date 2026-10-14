import { convertMetricIntervalToTimestamp } from '@/helpers/date/date-intervals';

export const queryIntervalFilterMixin = {
  inject: ['$system'],
  methods: {
    getIntervalQuery() {
      const { interval, with_history: withHistory } = this.query;

      if (!interval) {
        return {};
      }

      const result = convertMetricIntervalToTimestamp({
        interval,
        timezone: this.$system.timezone,
      });

      if (withHistory) {
        const { from: prevFrom, to: prevTo } = convertMetricIntervalToTimestamp({
          interval: result,
          timezone: this.$system.timezone,
        });

        result.prev_from = prevFrom;
        result.prev_to = prevTo;
      }

      return result;
    },

    updateInterval(interval) {
      this.query = {
        ...this.query,
        interval,
      };
    },
  },
};
