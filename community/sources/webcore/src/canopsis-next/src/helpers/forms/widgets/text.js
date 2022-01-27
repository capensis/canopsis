import { cloneDeep } from 'lodash';

import { QUICK_RANGES, TIME_UNITS, WIDGET_TYPES } from '@/constants';

import { widgetToForm } from '@/helpers/forms/widgets/common';

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
 * @typedef {Object} TextWidgetDateInterval
 * @property {number} periodValue
 * @property {DurationUnit} periodUnit
 * @property {string | number} tstart
 * @property {string | number} tstop
 */

/**
 * @typedef {Object} TextWidgetParameters
 * @property {TextWidgetDateInterval} dateInterval
 * @property {Filter} mfilter
 * @property {Object.<string, TextWidgetStats>} stats
 * @property {string} template
 */

/**
 * @typedef {Widget} TextWidget
 * @property {TextWidgetParameters} parameters
 */

/**
 * @typedef {TextWidgetParameters} TextWidgetParametersForm
 */

/**
 * @typedef {Widget} TextWidgetForm
 * @property {TextWidgetParametersForm} parameters
 */

/**
 * Convert text widget parameters to form
 *
 * @param {TextWidgetParameters} parameters
 * @return {TextWidgetParametersForm}
 */
const textWidgetParametersToForm = (parameters = {}) => ({
  dateInterval: parameters.dateInterval
    ? cloneDeep(parameters.dateInterval)
    : {
      periodValue: 1,
      periodUnit: TIME_UNITS.day,
      tstart: QUICK_RANGES.thisMonthSoFar.start,
      tstop: QUICK_RANGES.thisMonthSoFar.stop,
    },
  mfilter: parameters.mfilter ? cloneDeep(parameters.mfilter) : {},
  stats: parameters.stats ? cloneDeep(parameters.stats) : {},
  template: parameters.template ?? '',
});

/**
 * Convert text widget to form
 *
 * @param {TextWidget} [textWidget = {}]
 * @returns {TextWidgetForm}
 */
export const textWidgetToForm = (textWidget = {}) => {
  const widget = widgetToForm(textWidget);

  return {
    ...widget,
    type: WIDGET_TYPES.text,
    parameters: textWidgetParametersToForm(textWidget.parameters),
  };
};
