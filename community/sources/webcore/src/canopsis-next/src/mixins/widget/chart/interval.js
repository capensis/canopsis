import { DATETIME_FORMATS } from '@/constants';

import {
  convertStartDateIntervalToTimestampByTimezone,
  convertStopDateIntervalToTimestampByTimezone,
} from '@/helpers/date/date-intervals';

export const widgetIntervalFilterMixin = {
  inject: ['$system'],
  methods: {
    getIntervalQuery() {
      if (!this.query.interval) {
        return {};
      }

      return {
        from: convertStartDateIntervalToTimestampByTimezone(
          this.query.interval.from,
          DATETIME_FORMATS.datePicker,
          this.query.sampling,
          this.$system.timezone,
        ),
        to: convertStopDateIntervalToTimestampByTimezone(
          this.query.interval.to,
          DATETIME_FORMATS.datePicker,
          this.query.sampling,
          this.$system.timezone,
        ),
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
