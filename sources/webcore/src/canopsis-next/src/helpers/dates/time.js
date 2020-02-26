import moment from 'moment';

import { TIME_UNITS } from '@/constants';

export const getSecondsByUnit = (value, unit = TIME_UNITS.second) => moment(0).add(value, unit) / 1000;

export const getUnitValueFromSeconds = (
  value,
  unit = TIME_UNITS.second,
) => moment.duration(value, TIME_UNITS.second).as(unit);
