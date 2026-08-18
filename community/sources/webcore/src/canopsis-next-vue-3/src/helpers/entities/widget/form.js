import { cloneDeep } from 'lodash';

import { DEFAULT_WIDGET_GRID_PARAMETERS, WIDGET_GRID_SIZES_KEYS, WIDGET_TYPES } from '@/constants';

import { uid } from '@/helpers/uid';

import {
  alarmListWidgetParametersToForm,
  formToAlarmListWidgetParameters,
  prepareAlarmListWidget,
} from './forms/alarm';
import {
  contextWidgetParametersToForm,
  formToContextWidgetParameters,
  prepareContextWidget,
} from './forms/context';
import { counterWidgetParametersToForm, formToCounterWidgetParameters } from './forms/counter';
import {
  externalDataTableWidgetParametersToForm,
  formToExternalDataTableWidgetParameters,
} from './forms/external-data-table';
import { formToMapWidgetParameters, mapWidgetParametersToForm } from './forms/map';
import {
  formToServiceWeatherWidgetParameters,
  serviceWeatherWidgetParametersToForm,
} from './forms/service-weather';
import {
  formToStatsCalendarWidgetParameters,
  statsCalendarWidgetParametersToForm,
} from './forms/stats-calendar';
import { textWidgetParametersToForm } from './forms/text';

/**
 * @typedef {
 *   'AlarmsList' |
 *   'Context' |
 *   'ServiceWeather' |
 *   'StatsCalendar' |
 *   'Text' |
 *   'Counter' |
 *   'Junit' |
 *   'Map' |
 *   'BarChart' |
 *   'LineChart' |
 *   'PieChart' |
 *   'Numbers' |
 *   'UserStatistics' |
 *   'AlarmStatistics' |
 *   'Availability' |
 *   'ExternalDataTable'
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
 *   MapWidgetParameters |
 *   BarChartWidgetParameters |
 *   LineChartWidgetParameters |
 *   PieChartWidgetParameters |
 *   NumbersWidgetParameters |
 *   StatisticsWidgetParameters |
 *   AvailabilityWidgetParameters |
 *   ExternalDataTableWidgetParameters
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
 *   MapWidgetParameters |
 *   BarChartWidgetParametersForm |
 *   LineChartWidgetParametersForm |
 *   PieChartWidgetParametersForm |
 *   NumbersWidgetParametersForm |
 *   StatisticsWidgetParametersForm |
 *   AvailabilityWidgetParametersForm |
 *   ExternalDataTableWidgetParametersForm
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
    [WIDGET_TYPES.text]: textWidgetParametersToForm,
    [WIDGET_TYPES.counter]: counterWidgetParametersToForm,
    [WIDGET_TYPES.map]: mapWidgetParametersToForm,
    [WIDGET_TYPES.serviceWeather]: serviceWeatherWidgetParametersToForm,
    [WIDGET_TYPES.statsCalendar]: statsCalendarWidgetParametersToForm,
    [WIDGET_TYPES.externalDataTable]: externalDataTableWidgetParametersToForm,
  };

  const preparer = widgetsMap[type];

  return preparer ? preparer(parameters) : cloneDeep(parameters);
};

/**
 * Generate alarm list widget form with default parameters.
 *
 * @return {WidgetForm}
 */
export const generateDefaultAlarmListWidgetForm = () =>
  widgetToForm({ type: WIDGET_TYPES.alarmList });

/**
 * Generate alarm list widget with default parameters.
 *
 * @return {Widget}
 */
export const generateDefaultAlarmListWidget = () => ({
  ...formToWidget(generateDefaultAlarmListWidgetForm()),

  _id: uid(),
});

/**
 * Generate prepared default alarm list.
 *
 * @returns {Widget}
 */
export const generatePreparedDefaultAlarmListWidget = () =>
  prepareAlarmListWidget(generateDefaultAlarmListWidget());

/**
 * Generate context widget form with default parameters.
 *
 * @return {WidgetForm}
 */
export const generateDefaultContextWidgetForm = () => widgetToForm({ type: WIDGET_TYPES.context });

/**
 * Generate context widget with default parameters.
 *
 * @return {Widget}
 */
export const generateDefaultContextWidget = () => ({
  ...formToWidget(generateDefaultContextWidgetForm()),

  _id: uid(),
});

/**
 * Generate prepared default context widget.
 *
 * @returns {Widget}
 */
export const generatePreparedDefaultContextWidget = () =>
  prepareContextWidget(generateDefaultContextWidget());

/**
 * Get default widget grid parameters.
 *
 * @return {WidgetGridParameters}
 */
export const getDefaultGridParameters = () =>
  Object.values(WIDGET_GRID_SIZES_KEYS).reduce((acc, size) => {
    acc[size] = { ...DEFAULT_WIDGET_GRID_PARAMETERS };

    return acc;
  }, {});

/**
 * Get empty widget by widget type.
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
 * Convert widget parameters to payload
 *
 * @param {WidgetType} type
 * @param {WidgetParametersForm} parameters
 * @return {WidgetParameters}
 */
export const formToWidgetParameters = ({ type, parameters }) => {
  const widgetsMap = {
    [WIDGET_TYPES.alarmList]: formToAlarmListWidgetParameters,
    [WIDGET_TYPES.context]: formToContextWidgetParameters,
    [WIDGET_TYPES.counter]: formToCounterWidgetParameters,
    [WIDGET_TYPES.map]: formToMapWidgetParameters,
    [WIDGET_TYPES.serviceWeather]: formToServiceWeatherWidgetParameters,
    [WIDGET_TYPES.statsCalendar]: formToStatsCalendarWidgetParameters,
    [WIDGET_TYPES.externalDataTable]: formToExternalDataTableWidgetParameters,
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
