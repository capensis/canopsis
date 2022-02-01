import { cloneDeep } from 'lodash';

import { QUICK_RANGES, TIME_UNITS } from '@/constants';

/**
 * @typedef {Object} TextWidgetStats
 * @property {Object} stat
 * @property {string} stat.value
 * @property {string[]} stat.options
 * @property {boolean} trend
 * @property {Object} parameters
 * @property {boolean} parameters.recursive
 * @property {number[]} parameters.states
 *
 */

/**
 * @typedef {Object} TextWidgetParameters
 * @property {WidgetDateInterval} dateInterval
 * @property {Filter} mfilter
 * @property {Object.<string, TextWidgetStats>} stats
 * @property {string} template
 */

/**
 * Convert text widget parameters to form
 *
 * @param {TextWidgetParameters} parameters
 * @return {TextWidgetParameters}
 */
export const textWidgetParametersToForm = (parameters = {}) => ({
  dateInterval: parameters.dateInterval
    ? cloneDeep(parameters.dateInterval)
    : {
      periodValue: 1,
      periodUnit: TIME_UNITS.day,
      tstart: QUICK_RANGES.thisMonthSoFar.start,
      tstop: QUICK_RANGES.thisMonthSoFar.stop,
    },
  mfilter: parameters.mfilter
    ? cloneDeep(parameters.mfilter)
    : {},
  stats: parameters.stats
    ? cloneDeep(parameters.stats)
    : {},
  template: parameters.template ?? '',
});
