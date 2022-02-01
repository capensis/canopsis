import { cloneDeep } from 'lodash';

import { QUICK_RANGES, TIME_UNITS } from '@/constants';

/**
 * @typedef {Object} StatsWidgetBaseParameters
 * @property {Object} mfilter
 * @property {WidgetDateInterval} dateInterval
 */

/**
 * @typedef {StatsWidgetBaseParameters} StatsHistogramWidgetParameters
 * @property {Object} stats
 * @property {Object} statsColors
 * @property {Object} annotationLine
 */

/**
 * Convert stats histogram widget parameters to form
 *
 * @param {StatsHistogramWidgetParameters | {}} [parameters = {}]
 * @return {StatsHistogramWidgetParameters}
 */
export const statsHistogramWidgetParametersToForm = (parameters = {}) => ({
  mfilter: parameters.mfilter
    ? cloneDeep(parameters.mfilter)
    : {},
  dateInterval: parameters.dateInterval
    ? { ...parameters.dateInterval }
    : {
      periodValue: 1,
      periodUnit: TIME_UNITS.day,
      tstart: QUICK_RANGES.thisMonthSoFar.start,
      tstop: QUICK_RANGES.thisMonthSoFar.stop,
    },
  stats: parameters.stats
    ? cloneDeep(parameters.stats)
    : {},
  statsColors: parameters.statsColors
    ? cloneDeep(parameters.statsColors)
    : {},
  annotationLine: parameters.annotationLine
    ? cloneDeep(parameters.annotationLine)
    : {},
});
