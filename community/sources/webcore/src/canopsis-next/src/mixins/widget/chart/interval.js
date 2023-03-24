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

    updateInterval(interval) {
      this.query = {
        ...this.query,
        interval,
      };
    },
  },
};
