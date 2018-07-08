import moment from 'moment';

export default {
  today() {
    return ({
      tstart: moment()
        .startOf('day')
        .toDate()
        .getTime() / 1000,
      tstop: moment()
        .toDate()
        .getTime() / 1000,
    });
  },
  yesterday() {
    return ({
      tstart: moment()
        .subtract(1, 'day')
        .startOf('day')
        .toDate()
        .getTime() / 1000,
      tstop: moment()
        .subtract(1, 'day')
        .endOf('day')
        .toDate()
        .getTime() / 1000,
    });
  },
  last7Days() {
    return ({
      tstart: moment()
        .subtract(7, 'day')
        .toDate()
        .getTime() / 1000,
      tstop: moment()
        .toDate()
        .getTime() / 1000,
    });
  },
  last30Days() {
    return ({
      tstart: moment()
        .subtract(30, 'day')
        .toDate()
        .getTime() / 1000,
      tstop: moment()
        .toDate()
        .getTime() / 1000,
    });
  },
  thisMonth() {
    return ({
      tstart: moment()
        .startOf('month')
        .toDate()
        .getTime() / 1000,
      tstop: moment()
        .toDate()
        .getTime() / 1000,
    });
  },
  lastMonth() {
    return ({
      tstart: moment()
        .subtract(1, 'month')
        .startOf('month')
        .toDate()
        .getTime() / 1000,
      tstop: moment()
        .startOf('month')
        .toDate()
        .getTime() / 1000,
    });
  },
};
