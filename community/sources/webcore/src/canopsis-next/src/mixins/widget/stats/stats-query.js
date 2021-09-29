import { get, isString } from 'lodash';

import { DATETIME_FORMATS, DATETIME_INTERVAL_TYPES, TIME_UNITS } from '@/constants';

import {
  convertDateIntervalToMoment,
  prepareStatsStopForMonthPeriod,
  prepareStatsStartForMonthPeriod,
} from '@/helpers/date/date-intervals';

export default {
  methods: {
    getStatsQuery() {
      const { dateInterval, stats, mfilter } = this.query;
      const { periodValue } = dateInterval;
      let { periodUnit, tstart, tstop } = dateInterval;
      let filter = get(mfilter, 'filter', {});

      if (isString(filter)) {
        try {
          filter = JSON.parse(filter);
        } catch (err) {
          filter = {};

          console.error(err);
        }
      }

      tstart = convertDateIntervalToMoment(tstart, DATETIME_INTERVAL_TYPES.start, DATETIME_FORMATS.dateTimePicker);
      tstop = convertDateIntervalToMoment(tstop, DATETIME_INTERVAL_TYPES.stop, DATETIME_FORMATS.dateTimePicker);

      if (periodUnit === TIME_UNITS.month) {
        periodUnit = periodUnit.toUpperCase();

        /**
         * If period unit is 'month', we need to put the dates at the first day of the month, at 00:00 UTC
         * And add the difference between the local date, and the UTC one.
         */

        tstart = prepareStatsStartForMonthPeriod(tstart);
        tstop = prepareStatsStopForMonthPeriod(tstop);
      }

      return {
        stats,
        filter,
        tstart,
        tstop,
        periodUnit,
        periodValue,

        mfilter: filter,
        duration: `${periodValue}${periodUnit.toLowerCase()}`,
      };
    },
  },
};
