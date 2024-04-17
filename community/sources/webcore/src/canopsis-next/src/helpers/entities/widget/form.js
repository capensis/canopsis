import { cloneDeep } from 'lodash';

import { WIDGET_TYPES, WIDGET_GRID_SIZES_KEYS, DEFAULT_WIDGET_GRID_PARAMETERS } from '@/constants';

import featuresService from '@/services/features';

import { uuid } from '@/helpers/uuid';

import {
  alarmListWidgetParametersToForm,
  formToAlarmListWidgetParameters,
  prepareAlarmListWidget,
} from './forms/alarm';
import { counterWidgetParametersToForm, formToCounterWidgetParameters } from './forms/counter';
import { testingWeatherWidgetParametersToForm, formToTestingWeatherWidgetParameters } from './forms/junit';
import { statsCalendarWidgetParametersToForm, formToStatsCalendarWidgetParameters } from './forms/stats-calendar';
import { serviceWeatherWidgetParametersToForm, formToServiceWeatherWidgetParameters } from './forms/service-weather';
import { mapWidgetParametersToForm, formToMapWidgetParameters } from './forms/map';
import { barChartWidgetParametersToForm, formToBarChartWidgetParameters } from './forms/bar-chart';
import { lineChartWidgetParametersToForm, formToLineChartWidgetParameters } from './forms/line-chart';
import { pieChartWidgetParametersToForm, formToPieChartWidgetParameters } from './forms/pie-chart';
import { numbersWidgetParametersToForm, formToNumbersWidgetParameters } from './forms/numbers-chart';
import { contextWidgetParametersToForm, formToContextWidgetParameters, prepareContextWidget } from './forms/context';
import { textWidgetParametersToForm } from './forms/text';
import { statisticsWidgetParametersToForm, formToStatisticsWidgetParameters } from './forms/statistics';
import { availabilityWidgetParametersToForm, formToAvailabilityWidgetParameters } from './forms/availability';

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
 *   'Availability'
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
 *   MapWidgetParameters |
 *   BarChartWidgetParameters |
 *   LineChartWidgetParameters |
 *   PieChartWidgetParameters |
 *   NumbersWidgetParameters |
 *   StatisticsWidgetParameters
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
 *   StatisticsWidgetParametersForm
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
    [WIDGET_TYPES.barChart]: barChartWidgetParametersToForm,
    [WIDGET_TYPES.lineChart]: lineChartWidgetParametersToForm,
    [WIDGET_TYPES.pieChart]: pieChartWidgetParametersToForm,
    [WIDGET_TYPES.numbers]: numbersWidgetParametersToForm,
    [WIDGET_TYPES.userStatistics]: statisticsWidgetParametersToForm,
    [WIDGET_TYPES.alarmStatistics]: statisticsWidgetParametersToForm,
    [WIDGET_TYPES.availability]: availabilityWidgetParametersToForm,

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
    [WIDGET_TYPES.barChart]: formToBarChartWidgetParameters,
    [WIDGET_TYPES.lineChart]: formToLineChartWidgetParameters,
    [WIDGET_TYPES.pieChart]: formToPieChartWidgetParameters,
    [WIDGET_TYPES.numbers]: formToNumbersWidgetParameters,
    [WIDGET_TYPES.userStatistics]: formToStatisticsWidgetParameters,
    [WIDGET_TYPES.alarmStatistics]: formToStatisticsWidgetParameters,
    [WIDGET_TYPES.availability]: formToAvailabilityWidgetParameters,

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
 * Generate alarm list widget form with default parameters.
 *
 * @return {WidgetForm}
 */
export const generateDefaultAlarmListWidgetForm = () => widgetToForm({ type: WIDGET_TYPES.alarmList });

/**
 * Generate alarm list widget with default parameters.
 *
 * @return {Widget}
 */
export const generateDefaultAlarmListWidget = () => ({
  ...formToWidget(generateDefaultAlarmListWidgetForm()),

  _id: uuid(),
});

/**
 * Generate prepared default alarm list
 *
 * @returns {Widget}
 */
export const generatePreparedDefaultAlarmListWidget = () => prepareAlarmListWidget(generateDefaultAlarmListWidget());

/**
 * Generate context widget with default parameters.
 *
 * @return {Widget}
 */
export const generateDefaultContextWidget = () => ({
  ...formToWidget(widgetToForm({ type: WIDGET_TYPES.context })),

  _id: uuid(),
});

/**
 * Generate prepared default context
 *
 * @returns {Widget}
 */
export const generatePreparedDefaultContextWidget = () => prepareContextWidget(generateDefaultContextWidget());

/**
 * Generate service weather widget with default parameters.
 *
 * @return {Widget}
 */
export const generateDefaultServiceWeatherWidget = () => ({
  ...formToWidget(widgetToForm({ type: WIDGET_TYPES.serviceWeather })),

  _id: uuid(),
});
