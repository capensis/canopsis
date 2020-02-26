import moment from 'moment';

export default function convertTimestampToMoment(timestamp) {
  let dateObject;

  // If it's unix timestamp in seconds
  if (typeof timestamp === 'number' && timestamp < 100000000000) {
    dateObject = moment.unix(timestamp);
  } else {
    dateObject = moment(timestamp);
  }

  return dateObject;
}
