import { ALARM_METRIC_PARAMETERS, DATETIME_FORMATS, SAMPLINGS, TIME_UNITS, USER_METRIC_PARAMETERS } from '@/constants';
import { addUnitToDate, convertDateToString } from '@/helpers/date/date';

/**
 * @typedef { 'hour' | 'day' | 'week' | 'month' } Sampling
 */

/**
 * Check metric is time
 *
 * @param {string} metric
 * @returns {boolean}
 */
export const isTimeMetric = metric => [
  USER_METRIC_PARAMETERS.averageSession,
  ALARM_METRIC_PARAMETERS.averageAck,
  ALARM_METRIC_PARAMETERS.averageResolve,
].includes(metric);

/**
 * Check metric is ratio
 *
 * @param {string} metric
 * @returns {boolean}
 */
export const isRatioMetric = metric => [
  ALARM_METRIC_PARAMETERS.ratioCorrelation,
  ALARM_METRIC_PARAMETERS.ratioInstructions,
  ALARM_METRIC_PARAMETERS.ratioTickets,
  ALARM_METRIC_PARAMETERS.ratioNonDisplayed,
].includes(metric);

/**
 * Return label by sampling
 *
 * @param {number | string} value
 * @param {Sampling} sampling
 * @returns {string}
 *
 * @example
 * getDateLabelBySampling(1636523087405, 'hour') // 10/11/2021 12:44
 * getDateLabelBySampling(1636523087405, 'day') // 10/11/2021
 * getDateLabelBySampling(1636523087405, 'week') // 10/11/2021 - \n17/11/2021
 * getDateLabelBySampling(1636523087405, 'month') // November 2021
 */
export const getDateLabelBySampling = (value, sampling) => {
  switch (sampling) {
    case SAMPLINGS.hour:
      return convertDateToString(value, DATETIME_FORMATS.dateTimePicker);
    case SAMPLINGS.day:
      return convertDateToString(value, DATETIME_FORMATS.short);
    case SAMPLINGS.week:
      return [
        convertDateToString(value, DATETIME_FORMATS.short),
        convertDateToString(addUnitToDate(value, 1, TIME_UNITS.week), DATETIME_FORMATS.short),
      ].join(' - \n');
  }

  return convertDateToString(value, DATETIME_FORMATS.yearWithMonth);
};
