import { cloneDeep } from 'lodash';

import { QUICK_RANGES, TIME_UNITS } from '@/constants';

/**
 * @typedef {StatsWidgetBaseParameters} StatsCurvesWidgetParameters
 * @property {Object} stats
 * @property {Object} statsColors
 * @property {Object} statsPointsStyles
 * @property {Object} annotationLine
 */

/**
 * Convert stats curves widget parameters to form
 *
 * @param {StatsCurvesWidgetParameters | {}} [parameters = {}]
 * @return {StatsCurvesWidgetParameters}
 */
export const statsCurvesWidgetParametersToForm = (parameters = {}) => ({
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
  statsPointsStyles: parameters.statsPointsStyles
    ? cloneDeep(parameters.statsPointsStyles)
    : {},
  annotationLine: parameters.annotationLine
    ? cloneDeep(parameters.annotationLine)
    : {},
});
