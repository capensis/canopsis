import moment from 'moment';

import { STATS_DURATION_UNITS } from '@/constants';

/**
 * Helper to calculate time intervals
 */
export default {
  today() {
    return ({
      tstart: moment()
        .startOf('day')
        .unix(),
      tstop: moment()
        .unix(),
    });
  },
  yesterday() {
    return ({
      tstart: moment()
        .subtract(1, 'day')
        .startOf('day')
        .unix(),
      tstop: moment()
        .subtract(1, 'day')
        .endOf('day')
        .unix(),
    });
  },
  last7Days() {
    return ({
      tstart: moment()
        .subtract(7, 'day')
        .unix(),
      tstop: moment()
        .unix(),
    });
  },
  last30Days() {
    return ({
      tstart: moment()
        .subtract(30, 'day')
        .unix(),
      tstop: moment()
        .unix(),
    });
  },
  thisMonth() {
    return ({
      tstart: moment()
        .startOf('month')
        .unix(),
      tstop: moment()
        .unix(),
    });
  },
  lastMonth() {
    return ({
      tstart: moment()
        .subtract(1, 'month')
        .startOf('month')
        .unix(),
      tstop: moment()
        .startOf('month')
        .unix(),
    });
  },
};

/**
 * Convert a date interval string to moment date object
 *
 * @param {String} dateString
 * @param {String} type
 * @returns {Moment}
 */
export function parseStringToDateInterval(dateString, type) {
  const preparedDateString = dateString.toLowerCase().trim();
  const matches = preparedDateString.match(/^now(([+--])(\d+)([hdwmy]{1}))?(\/([hdwmy]{1}))?$/);

  if (matches) {
    const result = moment().utc();
    const operator = matches[2];
    const deltaValue = matches[3];
    let roundUnit = matches[6];
    let deltaUnit = matches[4];

    const roundMethod = type === 'start' ? 'startOf' : 'endOf';
    const deltaMethod = operator === '+' ? 'add' : 'subtract';

    if (roundUnit) {
      if (roundUnit === STATS_DURATION_UNITS.month) {
        roundUnit = roundUnit.toUpperCase();
      }

      result[roundMethod](roundUnit);
    }


    if (deltaValue && deltaUnit) {
      if (deltaUnit === STATS_DURATION_UNITS.month) {
        deltaUnit = deltaUnit.toUpperCase();
      }

      result[deltaMethod](deltaValue, deltaUnit);
    }

    return result;
  }

  throw new Error('Date string pattern not recognized');
}

export function dateParse(date, type, format) {
  const momentDate = moment(date, format, true);

  if (!momentDate.isValid()) {
    return parseStringToDateInterval(date, type);
  }

  return momentDate;
}

