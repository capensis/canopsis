import { cloneDeep, omit } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';
import {
  DEFAULT_ALARMS_WIDGET_COLUMNS,
  DEFAULT_PERIODIC_REFRESH,
  EXPORT_CSV_DATETIME_FORMATS,
  EXPORT_CSV_SEPARATORS,
  GRID_SIZES,
  SORT_ORDERS,
  TIME_UNITS,
  DEFAULT_SERVICE_DEPENDENCIES_COLUMNS,
  DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS,
  ALARM_DENSE_TYPES,
  DEFAULT_LINKS_INLINE_COUNT,
  WIDGET_TYPES,
} from '@/constants';

import { durationWithEnabledToForm, isValidUnit } from '@/helpers/date/duration';
import { addKeyInEntities, removeKeyFromEntities } from '@/helpers/entities';

import { widgetColumnsToForm, formToWidgetColumns } from '../shared/widget-column';
import { kioskParametersToForm } from '../shared/kiosk';
import { widgetTemplateValueToForm, formToWidgetTemplateValue } from '../widget-template';

import { openedToForm } from './common';
import { barChartWidgetParametersToForm, formToBarChartWidgetParameters } from './bar-chart';
import { lineChartWidgetParametersToForm, formToLineChartWidgetParameters } from './line-chart';
import { numbersWidgetParametersToForm, formToNumbersWidgetParameters } from './numbers';

/**
 * @typedef {'BarChart', 'LineChart', 'Numbers'} AlarmChartType
 */

/**
 * @typedef {Object} AlarmsListDataTableColumn
 * @property {string} value
 * @property {string} text
 * @property {boolean} [isHtml]
 * @property {boolean} [colorIndicator]
 */

/**
 * @typedef {Object} WidgetFastActionOutput
 * @property {boolean} enabled
 * @property {string} value
 */

/**
 * @typedef {Object} WidgetLinksCategoriesAsList
 * @property {boolean} enabled
 * @property {number} limit
 */

/**
 * @typedef {Object} WidgetLiveReporting
 * @property {string} [tstart]
 * @property {string} [tstop]
 */

/**
 * @typedef {Object} WidgetInfoPopup
 * @property {string} column
 * @property {string} template
 */

/**
 * @typedef {Object} AlarmChartBarChartWidgetParameters
 * @property {MetricPreset[]} metrics
 * @property {boolean} stacked
 * @property {string} chart_title
 * @property {string} default_time_range
 * @property {Sampling} default_sampling
 * @property {boolean} comparison
 */

/**
 * @typedef {Object} AlarmChartLineChartWidgetParameters
 * @property {MetricPreset[]} metrics
 * @property {string} chart_title
 * @property {string} default_time_range
 * @property {Sampling} default_sampling
 * @property {boolean} comparison
 */

/**
 * @typedef {Object} AlarmChartNumbersWidgetParameters
 * @property {MetricPreset[]} metrics
 * @property {string} chart_title
 * @property {string} default_time_range
 * @property {Sampling} default_sampling
 * @property {string} show_trend
 * @property {number} [font_size]
 */

/**
 * @typedef {Object} AlarmChart
 * @property {string} title
 * @property {AlarmChartType} type
 * @property {
 *   AlarmChartBarChartWidgetParameters
 *   | AlarmChartLineChartWidgetParameters
 *   | AlarmChartNumbersWidgetParameters
 * } parameters
 */

/**
 * @typedef {Object} AlarmListBaseParameters
 * @property {number} itemsPerPage
 * @property {WidgetSort} sort
 * @property {string} moreInfoTemplate
 * @property {string} moreInfoTemplateTemplate
 * @property {WidgetInfoPopup[]} infoPopups
 * @property {string} widgetColumnsTemplate
 * @property {WidgetColumn[]} widgetColumns
 */

/**
 * @typedef {Object} AlarmListWidgetDefaultParameters
 * @property {WidgetFastActionOutput} fastAckOutput
 * @property {WidgetFastActionOutput} fastCancelOutput
 * @property {number} inlineLinksCount
 * @property {number} itemsPerPage
 * @property {WidgetInfoPopup[]} infoPopups
 * @property {string} moreInfoTemplate
 * @property {string} moreInfoTemplateTemplate
 * @property {string} widgetColumnsTemplate
 * @property {string} widgetGroupColumnsTemplate
 * @property {string} widgetExportColumnsTemplate
 * @property {string} serviceDependenciesColumnsTemplate
 * @property {WidgetColumn[]} widgetColumns
 * @property {WidgetColumn[]} widgetGroupColumns
 * @property {WidgetColumn[]} widgetExportColumns
 * @property {WidgetColumn[]} serviceDependenciesColumns
 * @property {boolean} isAckNoteRequired
 * @property {boolean} isSnoozeNoteRequired
 * @property {boolean} isMultiAckEnabled
 * @property {boolean} isMultiDeclareTicketEnabled
 * @property {boolean} isHtmlEnabledOnTimeLine
 * @property {boolean} isActionsAllowWithOkState
 * @property {boolean} sticky_header
 * @property {boolean} dense
 */

/**
 * @typedef {AlarmListWidgetDefaultParameters} AlarmListWidgetParameters
 * @property {DurationWithEnabled} periodic_refresh
 * @property {string | null} mainFilter
 * @property {WidgetLiveReporting} liveReporting
 * @property {WidgetSort} sort
 * @property {boolean | null} opened
 * @property {number[]} expandGridRangeSize
 * @property {WidgetCsvSeparator} exportCsvSeparator
 * @property {string} exportCsvDatetimeFormat
 * @property {boolean} clearFilterDisabled
 * @property {WidgetKioskParameters} kiosk
 * @property {AlarmChart[]} charts
 */

/**
 * @typedef {AlarmListBaseParameters} AlarmListBaseParametersForm
 * @property {string | Symbol} widgetColumnsTemplate
 * @property {WidgetColumnForm[]} widgetColumns
 */

/**
 * @typedef {AlarmChart} AlarmChartForm
 * @property {string} key
 */

/**
 * @typedef {AlarmListWidgetDefaultParameters} AlarmListWidgetDefaultParametersForm
 * @property {string | Symbol} widgetColumnsTemplate
 * @property {string | Symbol} widgetGroupColumnsTemplate
 * @property {string | Symbol} widgetExportColumnsTemplate
 * @property {string | Symbol} serviceDependenciesColumnsTemplate
 * @property {WidgetColumnForm[]} widgetColumns
 * @property {WidgetColumnForm[]} widgetGroupColumns
 * @property {WidgetColumnForm[]} widgetExportColumns
 * @property {WidgetColumnForm[]} serviceDependenciesColumns
 */

/**
 * @typedef {AlarmListWidgetDefaultParametersForm & AlarmListWidgetParameters} AlarmListWidgetParametersForm
 * @property {AlarmChartForm[]} charts
 */

/**
 * Convert alarm list infoPopups parameters to form
 *
 * @param {WidgetInfoPopup[]} [infoPopups = []]
 * @return {WidgetInfoPopup[]}
 */
export const infoPopupsToForm = (infoPopups = []) => infoPopups.map(infoPopup => ({ ...infoPopup }));

/**
 * Convert widget sort parameters to form
 *
 * @param {WidgetSort} [sort = {}]
 * @return {WidgetSort}
 */
export const widgetSortToForm = (sort = {}) => ({
  order: sort.order ?? SORT_ORDERS.asc,
  column: sort.column ?? '',
});

/**
 * Convert alarm list base parameters (we are using it inside another widgets with alarmList) to form
 *
 * @param {AlarmListBaseParameters} [alarmListParameters = {}]
 * @return {AlarmListBaseParametersForm}
 */
export const alarmListBaseParametersToForm = (alarmListParameters = {}) => ({
  sort: widgetSortToForm(alarmListParameters.sort),
  itemsPerPage: alarmListParameters.itemsPerPage ?? PAGINATION_LIMIT,
  moreInfoTemplate: alarmListParameters.moreInfoTemplate ?? '',
  moreInfoTemplateTemplate: widgetTemplateValueToForm(alarmListParameters.moreInfoTemplateTemplate),
  infoPopups: infoPopupsToForm(alarmListParameters.infoPopups),
  widgetColumnsTemplate: widgetTemplateValueToForm(alarmListParameters.widgetColumnsTemplate),
  widgetColumns: widgetColumnsToForm(alarmListParameters.widgetColumns ?? DEFAULT_ALARMS_WIDGET_COLUMNS),
});

/**
 * Convert widget periodic refresh to form duration
 *
 * @param {DurationWithEnabled} periodicRefresh
 * @returns {DurationWithEnabled}
 */
export const periodicRefreshToDurationForm = (periodicRefresh = DEFAULT_PERIODIC_REFRESH) => {
  /**
   * @link https://git.canopsis.net/canopsis/canopsis-pro/-/issues/4390
   */
  const unit = isValidUnit(periodicRefresh.unit)
    ? periodicRefresh.unit
    : TIME_UNITS.second;

  return durationWithEnabledToForm({ ...periodicRefresh, unit });
};

/**
 * Convert alarm list chart to form
 *
 * @param {AlarmChart} [chart = {}]
 * @returns {AlarmChart}
 */
export const alarmListChartToForm = (chart = {}) => {
  const convertersMap = {
    [WIDGET_TYPES.barChart]: barChartWidgetParametersToForm,
    [WIDGET_TYPES.lineChart]: lineChartWidgetParametersToForm,
    [WIDGET_TYPES.numbers]: numbersWidgetParametersToForm,
  };

  const type = chart.type ?? WIDGET_TYPES.barChart;
  const converter = convertersMap[type];

  return {
    type,
    title: chart.title ?? '',
    parameters: omit(converter ? converter(chart.parameters) : chart.parameters, ['periodic_refresh']),
  };
};

/**
 * Convert alarm list widget parameters to form
 *
 * @param {AlarmListWidgetDefaultParameters} [parameters = {}]
 * @return {AlarmListWidgetDefaultParametersForm}
 */
export const alarmListWidgetDefaultParametersToForm = (parameters = {}) => ({
  itemsPerPage: parameters.itemsPerPage ?? PAGINATION_LIMIT,
  infoPopups: infoPopupsToForm(parameters.infoPopups),
  moreInfoTemplate: parameters.moreInfoTemplate ?? '',
  moreInfoTemplateTemplate: widgetTemplateValueToForm(parameters.moreInfoTemplateTemplate),
  isAckNoteRequired: !!parameters.isAckNoteRequired,
  isSnoozeNoteRequired: !!parameters.isSnoozeNoteRequired,
  isMultiAckEnabled: !!parameters.isMultiAckEnabled,
  isMultiDeclareTicketEnabled: !!parameters.isMultiDeclareTicketEnabled,
  isHtmlEnabledOnTimeLine: !!parameters.isHtmlEnabledOnTimeLine,
  isActionsAllowWithOkState: !!parameters.isActionsAllowWithOkState,
  sticky_header: !!parameters.sticky_header,
  dense: parameters.dense ?? ALARM_DENSE_TYPES.large,
  fastAckOutput: parameters.fastAckOutput
    ? { ...parameters.fastAckOutput }
    : {
      enabled: false,
      value: 'auto ack',
    },
  fastCancelOutput: parameters.fastCancelOutput
    ? { ...parameters.fastCancelOutput }
    : {
      enabled: false,
      value: 'auto cancel',
    },
  widgetColumnsTemplate: widgetTemplateValueToForm(parameters.widgetColumnsTemplate),
  widgetGroupColumnsTemplate: widgetTemplateValueToForm(parameters.widgetGroupColumnsTemplate),
  serviceDependenciesColumnsTemplate: widgetTemplateValueToForm(parameters.serviceDependenciesColumnsTemplate),
  widgetExportColumnsTemplate: widgetTemplateValueToForm(parameters.widgetExportColumnsTemplate),
  widgetColumns:
    widgetColumnsToForm(parameters.widgetColumns ?? DEFAULT_ALARMS_WIDGET_COLUMNS),
  widgetGroupColumns:
    widgetColumnsToForm(parameters.widgetGroupColumns ?? DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS),
  serviceDependenciesColumns:
    widgetColumnsToForm(parameters.serviceDependenciesColumns ?? DEFAULT_SERVICE_DEPENDENCIES_COLUMNS),
  widgetExportColumns:
    widgetColumnsToForm(parameters.widgetExportColumns ?? DEFAULT_ALARMS_WIDGET_COLUMNS),
  inlineLinksCount: parameters.inlineLinksCount ?? DEFAULT_LINKS_INLINE_COUNT,
});

/**
 * Convert alarm list widget parameters to form
 *
 * @param {AlarmListWidgetParameters} [parameters = {}]
 * @return {AlarmListWidgetParametersForm}
 */
export const alarmListWidgetParametersToForm = (parameters = {}) => ({
  ...parameters,
  ...alarmListWidgetDefaultParametersToForm(parameters),

  periodic_refresh: periodicRefreshToDurationForm(parameters.periodic_refresh),
  mainFilter: parameters.mainFilter ?? null,
  clearFilterDisabled: parameters.clearFilterDisabled ?? false,
  liveReporting: parameters.liveReporting
    ? cloneDeep(parameters.liveReporting)
    : {},
  sort: widgetSortToForm(parameters.sort),
  opened: openedToForm(parameters.opened),
  expandGridRangeSize: parameters.expandGridRangeSize
    ? [...parameters.expandGridRangeSize]
    : [GRID_SIZES.min, GRID_SIZES.max],
  exportCsvSeparator: parameters.exportCsvSeparator ?? EXPORT_CSV_SEPARATORS.comma,
  exportCsvDatetimeFormat: parameters.exportCsvDatetimeFormat ?? EXPORT_CSV_DATETIME_FORMATS.datetimeSeconds.value,
  kiosk: kioskParametersToForm(parameters.kiosk),
  charts: addKeyInEntities(parameters.charts),
});

/**
 * Convert form to alarm list base parameters (we are using it inside another widgets with alarmList)
 *
 * @param {AlarmListBaseParametersForm} [form = {}]
 * @return {AlarmListBaseParameters}
 */
export const formToAlarmListBaseParameters = (form = {}) => ({
  ...form,

  moreInfoTemplateTemplate: formToWidgetTemplateValue(form.moreInfoTemplateTemplate),
  widgetColumnsTemplate: formToWidgetTemplateValue(form.widgetColumnsTemplate),
  widgetColumns: formToWidgetColumns(form.widgetColumns),
});

/**
 * Convert form to alarm list chart
 *
 * @param {AlarmChartForm} chart
 * @returns {AlarmChart}
 */
export const formToAlarmListChart = ({ type, title, parameters }) => {
  const convertersMap = {
    [WIDGET_TYPES.barChart]: formToBarChartWidgetParameters,
    [WIDGET_TYPES.lineChart]: formToLineChartWidgetParameters,
    [WIDGET_TYPES.numbers]: formToNumbersWidgetParameters,
  };

  const converter = convertersMap[type];

  return {
    type,
    title,
    parameters: omit(converter ? converter(parameters) : parameters, ['periodic_refresh']),
  };
};

/**
 * Convert form parameters to alarm list widget parameters
 *
 * @param {AlarmListWidgetParametersForm} form
 * @return {AlarmListWidgetParameters}
 */
export const formToAlarmListWidgetParameters = form => ({
  ...form,

  moreInfoTemplateTemplate: formToWidgetTemplateValue(form.moreInfoTemplateTemplate),
  widgetColumnsTemplate: formToWidgetTemplateValue(form.widgetColumnsTemplate),
  widgetGroupColumnsTemplate: formToWidgetTemplateValue(form.widgetGroupColumnsTemplate),
  serviceDependenciesColumnsTemplate: formToWidgetTemplateValue(form.serviceDependenciesColumnsTemplate),
  widgetExportColumnsTemplate: formToWidgetTemplateValue(form.widgetExportColumnsTemplate),
  widgetColumns: formToWidgetColumns(form.widgetColumns),
  widgetGroupColumns: formToWidgetColumns(form.widgetGroupColumns),
  widgetExportColumns: formToWidgetColumns(form.widgetExportColumns),
  serviceDependenciesColumns: formToWidgetColumns(form.serviceDependenciesColumns),
  charts: removeKeyFromEntities(form.charts),
});
