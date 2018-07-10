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
