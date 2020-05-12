import moment from 'moment';
import { TIME_UNITS } from '@/constants';

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
/**
 * Convert duration to interval object
 *
 * @param duration
 * @return {{unit: string, interval: number}}
 */
export const convertDurationToIntervalObject = (duration) => {
  const durationUnits = [
    TIME_UNITS.year,
    TIME_UNITS.month,
    TIME_UNITS.week,
    TIME_UNITS.week,
    TIME_UNITS.day,
    TIME_UNITS.hour,
    TIME_UNITS.minute,
    TIME_UNITS.second,
  ];

  const durationType = durationUnits.find(unit => moment.duration(duration, 'seconds').as(unit) % 1 === 0);

  return {
    interval: moment.duration(duration, 'seconds').as(durationType),
    unit: durationType,
  };
};
