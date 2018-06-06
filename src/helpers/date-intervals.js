import moment from 'moment';

export default {
  today() {
    return ({
      tstart: moment().startOf('day').toDate(),
      tstop: moment().toDate(),
    });
  },
  yesterday() {
    return ({
      tstart: moment().subtract(1, 'day').startOf('day').toDate(),
      tstop: moment().subtract(1, 'day').endOf('day').toDate(),
    });
  },
  last7Days() {
    return ({
      tstart: moment().subtract(7, 'day').toDate(),
      tstop: moment().toDate(),
    });
  },
  last30Days() {
    return ({
      tstart: moment().subtract(30, 'day').toDate(),
      tstop: moment().toDate(),
    });
  },
  thisMonth() {
    return ({
      tstart: moment().startOf('month').toDate(),
      tstop: moment().toDate(),
    });
  },
  lastMonth() {
    return ({
      tstart: moment().subtract(1, 'month').startOf('month').toDate(),
      tstop: moment().startOf('month').toDate(),
    });
  },
};
