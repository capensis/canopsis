import moment from 'moment';

export default function (exdate) {
  return exdate.filter(({ value }) => value).map(({ value }) => moment(value).unix());
}
