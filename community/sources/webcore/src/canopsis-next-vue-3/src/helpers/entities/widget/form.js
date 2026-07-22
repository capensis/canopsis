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
 * Build minimal alarm-list widget config used by pattern preview modals.
 *
 * @returns {Object}
 */
export const generatePreparedDefaultAlarmListWidget = () => ({})

/**
 * Build minimal context widget config used by pattern preview modals.
 *
 * @returns {Object}
 */
export const generatePreparedDefaultContextWidget = () => ({})
