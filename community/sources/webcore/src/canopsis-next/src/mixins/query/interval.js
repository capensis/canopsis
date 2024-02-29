import { convertQueryIntervalToTimestamp } from '@/helpers/date/date-intervals';

export const queryIntervalFilterMixin = {
  inject: ['$system'],
  methods: {
    getIntervalQuery() {
      const { interval } = this.query;

      if (!interval) {
        return {};
      }

      return convertQueryIntervalToTimestamp({
        interval,
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
