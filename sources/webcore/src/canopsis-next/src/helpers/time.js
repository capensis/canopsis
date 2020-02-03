import moment from 'moment';

import { TIME_UNITS } from '@/constants';

export const getSecondByUnit = (value, unit = TIME_UNITS.seconds) => moment(0).add(value, unit) / 1000;

export const formatSecondByUnit = (
  second,
  unit = TIME_UNITS.seconds,
) => moment.duration(second, TIME_UNITS.seconds).as(unit);
