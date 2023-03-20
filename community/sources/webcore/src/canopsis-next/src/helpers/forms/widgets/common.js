import { cloneDeep, isBoolean, isNull } from 'lodash';

import {
  WIDGET_TYPES,
  WIDGET_GRID_SIZES_KEYS,
  DEFAULT_WIDGET_GRID_PARAMETERS,
} from '@/constants';

import featuresService from '@/services/features';

import {
  alarmListWidgetParametersToForm,
  formToAlarmListWidgetParameters,
} from './alarm';
import {
  contextWidgetParametersToForm,
  formToContextWidgetParameters,
} from './context';
import {
  serviceWeatherWidgetParametersToForm,
  formToServiceWeatherWidgetParameters,
} from './service-weather';
import {
  statsCalendarWidgetParametersToForm,
  formToStatsCalendarWidgetParameters,
} from './stats-calendar';
import {
  counterWidgetParametersToForm,
  formToCounterWidgetParameters,
} from './counter';
import {
  testingWeatherWidgetParametersToForm,
  formToTestingWeatherWidgetParameters,
} from './testing-weather';
import {
  textWidgetParametersToForm,
} from './text';
import { formToMapWidgetParameters, mapWidgetParametersToForm } from './map';

/**
 * @typedef {
 *   'AlarmsList' |
 *   'Context' |
 *   'ServiceWeather' |
 *   'StatsCalendar' |
 *   'Text' |
 *   'Counter' |
 *   'Junit' |
 *   'Map'
 * } WidgetType
 */

/**
 * @typedef { 'comma' | 'semicolon' | 'tab' | 'space' } WidgetCsvSeparator
 */

/**
 * @typedef {Object} WidgetFilter
 * @property {string} [_id]
 * @property {string} title
 * @property {string} query
 */

/**
 * @typedef {Object} WidgetColumn
 * @property {string} value
 * @property {string} label
 * @property {boolean} [isHtml]
 * @property {boolean} [colorIndicator]
 */

/**
 * @typedef {Object} WidgetSort
 * @property {string} order
 * @property {string} [column]
 */

/**
 * @typedef {Object} WidgetMargin
 * @property {number} top
 * @property {number} right
 * @property {number} bottom
 * @property {number} left
 */

/**
 * @typedef {Object} WidgetCriticityLevels
 * @property {number} minor
 * @property {number} major
 * @property {number} critical
 */

/**
 * @typedef {Object} WidgetCriticityLevelsColors
 * @property {string} ok
 * @property {string} minor
 * @property {string} major
 * @property {string} critical
 */

/**
 * @typedef {Object} WidgetDateInterval
 * @property {number} periodValue
 * @property {DurationUnit} periodUnit
 * @property {string} tstart
 * @property {string} tstop
 */

/**
 * @typedef {Object} WidgetGridParameter
 * @property {number} y
 * @property {number} x
 * @property {number} w
 * @property {number} h
 * @property {number} h
 * @property {boolean} autoHeight
 */

/**
 * @typedef {Object} WidgetGridParameters
 * @property {WidgetGridParameter} mobile
 * @property {WidgetGridParameter} tablet
 * @property {WidgetGridParameter} desktop
 */

/**
 * @typedef {
 *   AlarmListWidgetParameters |
 *   ContextWidgetParameters |
 *   ServiceWeatherWidgetParameters |
 *   StatsCalendarWidgetParameters |
 *   CounterWidgetParameters |
 *   TestingWeatherWidgetParameters |
 *   TextWidgetParameters |
 *   MapWidgetParameters
 * } WidgetParameters
 */

/**
 * @typedef {
 *   AlarmListWidgetParameters |
 *   ContextWidgetParameters |
 *   ServiceWeatherWidgetParameters |
 *   StatsCalendarWidgetParameters |
 *   CounterWidgetParameters |
 *   TestingWeatherWidgetParametersForm |
 *   TextWidgetParameters |
 *   MapWidgetParameters
 * } WidgetParametersForm
 */

/**
 * @typedef {Object} EmptyWidget
 * @property {WidgetType} type
 * @property {string} title
 * @property {Object} parameters
 * @property {WidgetFilter[]} [filters]
 * @property {WidgetGridParameters} grid_parameters
 */

/**
 * @typedef {EmptyWidget} Widget
 * @property {string} [_id]
 * @property {WidgetParameters} parameters
 */

/**
 * @typedef {Widget} WidgetForm
 * @typedef {WidgetParametersForm} parameters
 */

/**
 * Convert widget parameters to form
 *
 * @param {WidgetType} type
 * @param {WidgetParameters} [parameters = {}]
 * @return {WidgetParametersForm}
 */
export const widgetParametersToForm = ({ type, parameters = {} } = {}) => {
  const widgetsMap = {
    [WIDGET_TYPES.alarmList]: alarmListWidgetParametersToForm,
    [WIDGET_TYPES.context]: contextWidgetParametersToForm,
    [WIDGET_TYPES.serviceWeather]: serviceWeatherWidgetParametersToForm,
    [WIDGET_TYPES.statsCalendar]: statsCalendarWidgetParametersToForm,
    [WIDGET_TYPES.counter]: counterWidgetParametersToForm,
    [WIDGET_TYPES.testingWeather]: testingWeatherWidgetParametersToForm,
    [WIDGET_TYPES.text]: textWidgetParametersToForm,
    [WIDGET_TYPES.map]: mapWidgetParametersToForm,

    ...featuresService.get('helpers.forms.widgets.widgetParametersToForm.widgetsMap'),
  };

  const preparer = widgetsMap[type];

  return preparer ? preparer(parameters) : cloneDeep(parameters);
};

/**
 * Get default widget grid parameters
 *
 * @return {WidgetGridParameters}
 */
export const getDefaultGridParameters = () => Object.values(WIDGET_GRID_SIZES_KEYS).reduce((acc, size) => {
  acc[size] = { ...DEFAULT_WIDGET_GRID_PARAMETERS };

  return acc;
}, {});

/**
 * Get empty widget by widget type
 *
 * @param {WidgetType} type
 * @return {EmptyWidget}
 */
export const getEmptyWidgetByType = type => ({
  type,
  title: '',
  parameters: {},
  grid_parameters: getDefaultGridParameters(),
});

/**
 * Convert widget to form object
 *
 * @param {Widget} [widget = {}]
 * @returns {WidgetForm}
 */
export const widgetToForm = (widget = { type: WIDGET_TYPES.alarmList }) => ({
  type: widget.type,
  title: widget.title ?? '',
  filters: widget.filters?.length ? [...widget.filters] : [],
  parameters: widgetParametersToForm(widget),
  grid_parameters: widget.grid_parameters
    ? cloneDeep(widget.grid_parameters)
    : getDefaultGridParameters(),
});

/**
 * Convert widget parameters to form
 *
 * @param {WidgetType} type
 * @param {WidgetParametersForm} parameters
 * @return {WidgetParameters}
 */
export const formToWidgetParameters = ({ type, parameters }) => {
  const widgetsMap = {
    [WIDGET_TYPES.alarmList]: formToAlarmListWidgetParameters,
    [WIDGET_TYPES.context]: formToContextWidgetParameters,
    [WIDGET_TYPES.serviceWeather]: formToServiceWeatherWidgetParameters,
    [WIDGET_TYPES.statsCalendar]: formToStatsCalendarWidgetParameters,
    [WIDGET_TYPES.counter]: formToCounterWidgetParameters,
    [WIDGET_TYPES.testingWeather]: formToTestingWeatherWidgetParameters,
    [WIDGET_TYPES.map]: formToMapWidgetParameters,

    ...featuresService.get('helpers.forms.widgets.formToWidgetParameters.widgetsMap'),
  };

  const preparer = widgetsMap[type];

  return preparer ? preparer(parameters) : cloneDeep(parameters);
};

/**
 * Convert form object to widget
 *
 * @param {WidgetForm} form
 * @returns {Widget}
 */
export const formToWidget = form => ({
  ...form,

  parameters: formToWidgetParameters(form),
});

/**
 * Convert opened field widget
 *
 * @param  {boolean | null} [opened]
 * @returns {boolean | null}
 */
export const openedToForm = (opened) => {
  if (isBoolean(opened) || isNull(opened)) {
    return opened;
  }

  return true;
};
