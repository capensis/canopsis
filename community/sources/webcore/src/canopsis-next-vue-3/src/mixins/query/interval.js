import { convertMetricIntervalToTimestamp } from '@/helpers/date/date-intervals';

/**
 * Mixin for components that keep metric interval filters in `query`.
 */
export const queryIntervalFilterMixin = {
  inject: ['$system'],

  methods: {
    /**
     * Convert current interval query into request interval params.
     *
     * @returns {Object}
     */
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

    /**
     * Update interval value in local query.
     *
     * @param {Object} interval
     */
    updateInterval(interval) {
      this.query = {
        ...this.query,
        interval,
      };
    },
  },
};
