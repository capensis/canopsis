import { cloneDeep } from 'lodash';

import {
  ALARM_FIELDS_TO_LABELS_KEYS,
  ALARM_UNSORTABLE_FIELDS,
  DEFAULT_ALARMS_WIDGET_COLUMNS,
  ALARM_LEVELS,
  ALARM_LEVELS_COLORS,
} from '@/constants';

import { setSeveralFields } from '@/helpers/immutable';
import { getWidgetColumnLabel, getWidgetColumnSortable } from '@/helpers/entities/widget/table';
import { openedToForm, alarmListBaseParametersToForm, formToAlarmListBaseParameters } from '@/helpers/entities/alarm/widget';

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
 * @typedef {StatsCalendarWidgetParameters} StatsCalendarWidgetParametersForm
 * @property {AlarmListBaseParametersForm} alarmsList
 */

/**
 * Convert stats calendar widget parameters to form
 *
 * @param {StatsCalendarWidgetParameters} [parameters = {}]
 * @return {StatsCalendarWidgetParametersForm}
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
 * @param {StatsCalendarWidgetParametersForm} form
 * @return {StatsCalendarWidgetParameters}
 */
export const formToStatsCalendarWidgetParameters = form => ({
  ...form,

  alarmsList: formToAlarmListBaseParameters(form.alarmsList),
});

/**
 * Prepared stats calendar/counter widget for displaying
 *
 * @param {Object} widget
 * @returns {Object}
 */
export const prepareStatsCalendarAndCounterWidget = (widget = {}) => setSeveralFields(widget, {
  'parameters.alarmsList.widgetColumns': (columns = DEFAULT_ALARMS_WIDGET_COLUMNS) => (
    columns.map(column => ({
      ...column,

      sortable: getWidgetColumnSortable(column, ALARM_UNSORTABLE_FIELDS),
      text: getWidgetColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
    }))
  ),
});
