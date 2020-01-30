import moment from 'moment';

export const getSecondByUnit = (value, unit = 's') => moment(0).add(value, unit) / 1000;
