import moment from 'moment-timezone';
import { get, isString } from 'lodash';

import { DATETIME_FORMATS, STATS_DURATION_UNITS } from '@/constants';

import { dateParse } from '@/helpers/date-intervals';

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

      tstart = dateParse(tstart, 'start', DATETIME_FORMATS.dateTimePicker);
      tstop = dateParse(tstop, 'stop', DATETIME_FORMATS.dateTimePicker);

      if (periodUnit === STATS_DURATION_UNITS.month) {
        periodUnit = periodUnit.toUpperCase();

        /**
         * If period unit is 'month', we need to put the dates at the first day of the month, at 00:00 UTC
         * And add the difference between the local date, and the UTC one.
         */
        tstart = moment.utc(tstart).startOf('month').tz(moment.tz.guess());

        const startOfTstopMonthUtc = moment.utc(tstop).startOf('month');
        const startOfCurrentMonthUtc = moment.utc().startOf('month');

        if (startOfTstopMonthUtc.isSame(startOfCurrentMonthUtc)) {
          tstop = startOfTstopMonthUtc.add(1, 'month').tz(moment.tz.guess());
        } else {
          tstop = startOfTstopMonthUtc.tz(moment.tz.guess());
        }
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
