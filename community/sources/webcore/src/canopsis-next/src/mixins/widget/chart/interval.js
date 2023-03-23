import { DATETIME_FORMATS } from '@/constants';

import { convertIntervalToTimestamp } from '@/helpers/date/date-intervals';

export const widgetIntervalFilterMixin = {
  inject: ['$system'],
  methods: {
    getIntervalQuery() {
      const { interval, sampling } = this.query;

      if (!interval) {
        return {};
      }

      return convertIntervalToTimestamp(interval, DATETIME_FORMATS.datePicker, sampling, this.$system.timezone);
    },

    getPrevIntervalQuery() {
      const { from = 0, to = 0 } = this.getIntervalQuery();
      const diffInSeconds = to - from;
      const newTo = from;
      const newFrom = newTo - diffInSeconds;

      return {
        from: newFrom,
        to: newTo,
      };
    },

    updateInterval(interval) {
      this.query = {
        ...this.query,
        interval,
      };
    },
  },
};
