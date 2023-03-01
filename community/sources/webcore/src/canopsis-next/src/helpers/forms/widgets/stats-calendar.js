import { cloneDeep } from 'lodash';

import { ALARM_LEVELS, ALARM_LEVELS_COLORS } from '@/constants';

import { openedToForm } from './common';
import { alarmListBaseParametersToForm, formToAlarmListBaseParameters } from './alarm';

/**
 * @typedef {Object} StatsCalendarWidgetParameters
 * @property {WidgetFilter[]} filters
 * @property {boolean} opened
 * @property {boolean} considerPbehaviors
 * @property {WidgetCriticityLevelsColors} criticityLevelsColors
 * @property {WidgetCriticityLevels} criticityLevels
 * @property {AlarmListBaseParameters} alarmsList
 */

/**
 * Convert stats calendar widget parameters to form
 *
 * @param {StatsCalendarWidgetParameters} [parameters = {}]
 * @return {StatsCalendarWidgetParameters}
 */
export const statsCalendarWidgetParametersToForm = (parameters = {}) => ({
  filters: parameters.filters
    ? cloneDeep(parameters.filters)
    : [],
  opened: openedToForm(parameters.opened),
  considerPbehaviors: parameters.considerPbehaviors ?? false,
  criticityLevels: parameters.criticityLevels
    ? { ...parameters.criticityLevels }
    : { ...ALARM_LEVELS },
  criticityLevelsColors: parameters.criticityLevelsColors
    ? { ...parameters.criticityLevelsColors }
    : { ...ALARM_LEVELS_COLORS },
  alarmsList: alarmListBaseParametersToForm(parameters.alarmsList),
});

/**
 * Convert form to stats calendar widget parameters
 *
 * @param {StatsCalendarWidgetParameters} form
 * @return {StatsCalendarWidgetParameters}
 */
export const formToStatsCalendarWidgetParameters = form => ({
  ...form,

  alarmsList: formToAlarmListBaseParameters(form.alarmsList),
});
