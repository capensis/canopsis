import { convertMetricIntervalToTimestamp } from '@/helpers/date/date-intervals';

export const widgetIntervalFilterMixin = {
  inject: ['$system'],
  methods: {
    getIntervalQuery() {
      const { interval, sampling } = this.query;

      if (!interval) {
        return {};
      }

      return convertMetricIntervalToTimestamp({
        interval,
        sampling,
        timezone: this.$system.timezone,
      });
    },

    updateInterval(interval) {
      this.query = {
        ...this.query,
        interval,
      };
    },
  },
};
