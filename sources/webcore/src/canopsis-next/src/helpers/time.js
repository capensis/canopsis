import moment from 'moment';

import { DEFAULT_PERIODIC_REFRESH } from '@/constants';

export const getSecondByUnit = (value, unit = DEFAULT_PERIODIC_REFRESH.unit) => moment(0).add(value, unit) / 1000;
