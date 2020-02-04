import moment from 'moment';

import { TIME_UNITS } from '@/constants';

export const getSecondByUnit = (value, unit = TIME_UNITS.second) => moment(0).add(value, unit) / 1000;
