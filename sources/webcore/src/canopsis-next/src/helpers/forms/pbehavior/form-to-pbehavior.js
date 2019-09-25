import moment from 'moment';

export default function (form) {
  return {
    ...form,

    comments: [],
    tstart: moment(form.tstart).unix(),
    tstop: moment(form.tstop).unix(),
  };
}
