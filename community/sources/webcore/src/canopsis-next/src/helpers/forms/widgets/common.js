import { cloneDeep } from 'lodash';

import {
  WIDGET_TYPES,
  WIDGET_GRID_SIZES_KEYS,
  DEFAULT_WIDGET_GRID_PARAMETERS,
} from '@/constants';

import {
  alarmListWidgetParametersToForm,
  formToAlarmListWidgetParameters,
} from './alarm';
import {
  contextWidgetParametersToForm,
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
export const widgetParametersToForm = ({ type, parameters } = {}) => {
  switch (type) {
    case WIDGET_TYPES.alarmList:
      return alarmListWidgetParametersToForm(parameters);
    case WIDGET_TYPES.context:
      return contextWidgetParametersToForm(parameters);
    case WIDGET_TYPES.serviceWeather:
      return serviceWeatherWidgetParametersToForm(parameters);
    case WIDGET_TYPES.statsCalendar:
      return statsCalendarWidgetParametersToForm(parameters);
    case WIDGET_TYPES.counter:
      return counterWidgetParametersToForm(parameters);
    case WIDGET_TYPES.testingWeather:
      return testingWeatherWidgetParametersToForm(parameters);
    case WIDGET_TYPES.text:
      return textWidgetParametersToForm(parameters);
    case WIDGET_TYPES.map:
      return mapWidgetParametersToForm(parameters);
    default:
      return parameters ? cloneDeep(parameters) : {};
  }
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
  switch (type) {
    case WIDGET_TYPES.alarmList:
      return formToAlarmListWidgetParameters(parameters);
    case WIDGET_TYPES.serviceWeather:
      return formToServiceWeatherWidgetParameters(parameters);
    case WIDGET_TYPES.statsCalendar:
      return formToStatsCalendarWidgetParameters(parameters);
    case WIDGET_TYPES.counter:
      return formToCounterWidgetParameters(parameters);
    case WIDGET_TYPES.testingWeather:
      return formToTestingWeatherWidgetParameters(parameters);
    case WIDGET_TYPES.map:
      return formToMapWidgetParameters(parameters);
    default:
      return parameters;
  }
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
