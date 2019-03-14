import moment from 'moment';

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
 * @returns {moment}
 */
export function parseStringToDateInterval(dateString, type) {
  let operator = null;
  let deltaValue = 0;
  let deltaUnit = null;
  let roundUnit = null;
  let preparedDateString = dateString.toLowerCase().trim();

  if (preparedDateString.match(/^now/)) {
    preparedDateString = preparedDateString.substring(3);
  } else {
    throw new Error('Date string pattern not recognized');
  }

  if (!preparedDateString) {
    return moment().utc();
  } else if (preparedDateString.match(/\/[hdwmy]{1}$/) && preparedDateString.match(/^[+--]\d+[hdwmy]{1}/)) {
    [operator] = preparedDateString;
    roundUnit = preparedDateString.slice(-1);
    deltaUnit = preparedDateString.slice(-3, -2);
    deltaValue = preparedDateString.slice(1, -3);

    if (deltaUnit === 'm') {
      deltaUnit = deltaUnit.toUpperCase();
    }

    if (type === 'start') {
      return operator === '+' ?
        moment().utc().startOf(roundUnit).add(deltaValue, deltaUnit) :
        moment().utc().startOf(roundUnit).subtract(deltaValue, deltaUnit);
    }

    if (type === 'stop') {
      return operator === '+' ?
        moment().utc().endOf(roundUnit).add(deltaValue, deltaUnit) :
        moment().utc().endOf(roundUnit).subtract(deltaValue, deltaUnit);
    }
  } else if (preparedDateString.match(/^[+--]\d+[hdwmy]{1}$/)) {
    [operator] = preparedDateString;
    deltaUnit = preparedDateString.slice(-1);
    deltaValue = preparedDateString.slice(1, -1);

    if (deltaUnit === 'm') {
      deltaUnit = deltaUnit.toUpperCase();
    }

    if (type === 'start') {
      return operator === '+' ?
        moment().utc().add(deltaValue, deltaUnit) :
        moment().utc().subtract(deltaValue, deltaUnit);
    }

    if (type === 'stop') {
      return operator === '+' ?
        moment().utc().add(deltaValue, deltaUnit) :
        moment().utc().subtract(deltaValue, deltaUnit);
    }
  } else if (preparedDateString.match(/^\/[hdwmy]{1}$/)) {
    roundUnit = preparedDateString.slice(-1);

    if (roundUnit === 'm') {
      roundUnit = roundUnit.toUpperCase();
    }

    if (type === 'start') {
      return moment().utc().startOf(roundUnit);
    }

    if (type === 'stop') {
      return moment().utc().endOf(roundUnit);
    }
  }

  throw new Error('Date string pattern not recognized');
}

export function dateParse(date, type, format) {
  const momentDate = moment(date, format);

  if (!momentDate.isValid()) {
    return parseStringToDateInterval(date, type);
  }

  return momentDate;
}
