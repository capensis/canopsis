import { DATETIME_FORMATS, SAMPLINGS } from '@/constants';

import {
  convertStartDateIntervalToTimestampByTimezone,
  convertStopDateIntervalToTimestampByTimezone,
} from '@/helpers/date/date-intervals';

export const widgetIntervalFilterMixin = {
  methods: {
    updateInterval(interval) {
      this.query = {
        ...this.query,
        interval,
      };
    },

    getIntervalQuery() {
      if (!this.query.interval) {
        return {};
      }

      return {
        from: convertStartDateIntervalToTimestampByTimezone(
          this.query.interval.from,
          DATETIME_FORMATS.datePicker,
          SAMPLINGS.day,
          this.$system.timezone,
        ),
        to: convertStopDateIntervalToTimestampByTimezone(
          this.query.interval.to,
          DATETIME_FORMATS.datePicker,
          SAMPLINGS.day,
          this.$system.timezone,
        ),
      };
    },
  },
};
