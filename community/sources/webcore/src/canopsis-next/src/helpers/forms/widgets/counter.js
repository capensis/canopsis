import { cloneDeep } from 'lodash';

import {
  ALARM_LEVELS,
  ALARM_LEVELS_COLORS,
  AVAILABLE_COUNTERS,
  DEFAULT_COUNTER_BLOCK_TEMPLATE,
  DEFAULT_WIDGET_MARGIN,
} from '@/constants';

import {
  alarmListBaseParametersToForm,
  formToAlarmListBaseParameters,
} from './alarm';
import { openedToForm } from './common';

/**
 * @typedef {Object} CounterWidgetLevels
 * @property {number} counter
 * @property {WidgetCriticityLevels} values
 * @property {WidgetCriticityLevelsColors} colors
 */

/**
 * @typedef {Object} CounterWidgetParameters
 * @property {boolean} opened
 * @property {string} blockTemplate
 * @property {number} columnSM
 * @property {number} columnMD
 * @property {number} columnLG
 * @property {WidgetMargin} margin
 * @property {number} heightFactor
 * @property {boolean} isCorrelationEnabled
 * @property {CounterWidgetLevels} levels
 * @property {AlarmListBaseParameters} alarmsList
 */

/**
 * Convert counter widget parameters to form
 *
 * @param {CounterWidgetParameters} [parameters = {}]
 * @return {CounterWidgetParameters}
 */
export const counterWidgetParametersToForm = (parameters = {}) => ({
  opened: openedToForm(parameters.opened),
  blockTemplate: parameters.blockTemplate ?? DEFAULT_COUNTER_BLOCK_TEMPLATE,
  columnSM: parameters.columnSM ?? 6,
  columnMD: parameters.columnMD ?? 4,
  columnLG: parameters.columnLG ?? 3,
  heightFactor: parameters.heightFactor ?? 6,
  margin: parameters.margin
    ? { ...parameters.margin }
    : { ...DEFAULT_WIDGET_MARGIN },
  isCorrelationEnabled: parameters.isCorrelationEnabled ?? false,
  levels: parameters.levels
    ? cloneDeep(parameters.levels)
    : {
      counter: AVAILABLE_COUNTERS.total,
      colors: { ...ALARM_LEVELS_COLORS },
      values: { ...ALARM_LEVELS },
    },
  alarmsList: alarmListBaseParametersToForm(parameters.alarmsList),
});

/**
 * Convert form to counter widget parameters
 *
 * @param {CounterWidgetParameters} form
 * @return {CounterWidgetParameters}
 */
export const formToCounterWidgetParameters = form => ({
  ...form,

  alarmsList: formToAlarmListBaseParameters(form.alarmsList),
});
