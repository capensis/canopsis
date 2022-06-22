import { isString, isUndefined } from 'lodash';

import {
  DEFAULT_ALARMS_WIDGET_COLUMNS,
  DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS,
  DEFAULT_PERIODIC_REFRESH,
  DEFAULT_SERVICE_DEPENDENCIES_COLUMNS,
  EXPORT_CSV_DATETIME_FORMATS,
  EXPORT_CSV_SEPARATORS,
  GRID_SIZES,
  SORT_ORDERS,
  TIME_UNITS,
  WIDGET_TYPES,
} from '@/constants';
import { DEFAULT_CATEGORIES_LIMIT, PAGINATION_LIMIT } from '@/config';

import { defaultColumnsToColumns } from '@/helpers/entities';
import { widgetToForm } from '@/helpers/forms/widgets/common';
import { durationWithEnabledToForm, isValidUnit } from '@/helpers/date/duration';

/**
 * @typedef {Object} FastAckOutput
 * @property {boolean} enabled
 * @property {string} value
 */

/**
 * @typedef {Object} LinksCategoriesAsList
 * @property {boolean} enabled
 * @property {number} limit
 */

/**
 * @typedef {Object} LiveReporting
 * @property {string} [tstart]
 * @property {string} [tstop]
 */

/**
 * @typedef {Object} Sort
 * @property {string} order
 * @property {string} column
 */

/**
 * @typedef {Object} InfoPopup
 * @property {string} column
 * @property {string} template
 */

/**
 * @typedef {Object} AlarmListWidgetDefaultParameters
 * @property {FastAckOutput} fastAckOutput
 * @property {LinksCategoriesAsList} linksCategoriesAsList
 * @property {number} itemsPerPage
 * @property {InfoPopup[]} infoPopups
 * @property {string} moreInfoTemplate
 * @property {WidgetColumn[]} widgetColumns
 * @property {WidgetColumn[]} widgetGroupColumns
 * @property {WidgetColumn[]} widgetExportColumns
 * @property {WidgetColumn[]} serviceDependenciesColumns
 * @property {boolean} isAckNoteRequired
 * @property {boolean} isSnoozeNoteRequired
 * @property {boolean} isMultiAckEnabled
 * @property {boolean} isHtmlEnabledOnTimeLine
 * @property {boolean} sticky_header
 * @property {boolean} dense
 */

/**
 * @typedef {AlarmListWidgetDefaultParameters} AlarmListWidgetDefaultParametersForm
 */

/**
 * @typedef {AlarmListWidgetDefaultParameters} AlarmListWidgetParameters
 * @property {DurationWithEnabled} periodic_refresh
 * @property {Array} viewFilters
 * @property {Object|null} mainFilter
 * @property {number} mainFilterUpdatedAt
 * @property {LiveReporting} liveReporting
 * @property {Sort} sort
 * @property {boolean | null} opened
 * @property {number[]} expandGridRangeSize
 * @property {CsvSeparators} exportCsvSeparator
 * @property {string} exportCsvDatetimeFormat
 */

/**
 * @typedef {AlarmListWidgetParameters} AlarmListWidgetParametersForm
 * @property {DurationWithEnabledForm} periodic_refresh
 */

/**
 * @typedef {Widget} AlarmListWidget
 * @property {AlarmListWidgetParameters} parameters
 */

/**
 * @typedef {Widget & AlarmListWidget} AlarmListWidgetForm
 * @property {AlarmListWidgetParametersForm} parameters
 */

/**
 * Add 'alarm.' prefix in column value
 *
 * @param {string} [value]
 * @returns {string}
 */
export const columnValueToForm = value => (value && isString(value) ? value.replace(/^v\./, 'alarm.v.') : value);

/**
 * Remove 'alarm.' prefix from column value
 *
 * @param {string} [value]
 * @returns {string}
 */
export const formToColumnValue = value => (value && isString(value) ? value.replace(/^alarm\./, '') : value);

/**
 * Convert columns parameters to form
 *
 * @param {WidgetColumn[]} [widgetColumns = []]
 * @return {WidgetColumn[]}
 */
export const widgetColumnsToForm = (widgetColumns = []) => widgetColumns.map(column => ({
  ...column,
  value: columnValueToForm(column.value),
}));

/**
 * Convert alarm list infoPopups parameters to form
 *
 * @param {InfoPopup[]} [infoPopups = []]
 * @return {InfoPopup[]}
 */
const infoPopupsToForm = (infoPopups = []) => infoPopups.map(infoPopup => ({
  ...infoPopup,
  column: columnValueToForm(infoPopup.column),
}));

/**
 * Convert widget sort parameters to form
 *
 * @param {Sort} [sort = {}]
 * @return {Sort}
 */
const widgetSortToForm = (sort = {}) => ({
  order: sort.order || SORT_ORDERS.asc,
  column: sort.column ? columnValueToForm(sort.column) : '',
});

/**
 * Convert widget periodic refresh to form duration
 *
 * @param {DurationWithEnabled} periodicRefresh
 * @returns {DurationWithEnabled}
 */
export const periodicRefreshToDurationForm = (periodicRefresh = DEFAULT_PERIODIC_REFRESH) => {
  /*
  * @link https://git.canopsis.net/canopsis/canopsis-pro/-/issues/4390
  */
  const unit = isValidUnit(periodicRefresh.unit)
    ? periodicRefresh.unit
    : TIME_UNITS.second;

  return durationWithEnabledToForm({ ...periodicRefresh, unit });
};

/**
 * Convert alarm list widget parameters to form
 *
 * @param {AlarmListWidgetDefaultParameters} [parameters = {}]
 * @return {AlarmListWidgetDefaultParametersForm}
 */
export const alarmListWidgetDefaultParametersToForm = (parameters = {}) => ({
  itemsPerPage: parameters.itemsPerPage || PAGINATION_LIMIT,
  infoPopups: infoPopupsToForm(parameters.infoPopups),
  moreInfoTemplate: parameters.moreInfoTemplate || '',
  isAckNoteRequired: !!parameters.isAckNoteRequired,
  isSnoozeNoteRequired: !!parameters.isSnoozeNoteRequired,
  isMultiAckEnabled: !!parameters.isMultiAckEnabled,
  isHtmlEnabledOnTimeLine: !!parameters.isHtmlEnabledOnTimeLine,
  sticky_header: !!parameters.sticky_header,
  dense: !!parameters.dense,
  fastAckOutput: parameters.fastAckOutput || {
    enabled: false,
    value: 'auto ack',
  },
  widgetColumns: parameters.widgetColumns
    ? widgetColumnsToForm(parameters.widgetColumns)
    : defaultColumnsToColumns(DEFAULT_ALARMS_WIDGET_COLUMNS),
  widgetGroupColumns: parameters.widgetGroupColumns
    ? widgetColumnsToForm(parameters.widgetGroupColumns)
    : defaultColumnsToColumns(DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS),
  serviceDependenciesColumns: parameters.serviceDependenciesColumns
    ? widgetColumnsToForm(parameters.serviceDependenciesColumns)
    : defaultColumnsToColumns(DEFAULT_SERVICE_DEPENDENCIES_COLUMNS),
  linksCategoriesAsList: parameters.linksCategoriesAsList || {
    enabled: false,
    limit: DEFAULT_CATEGORIES_LIMIT,
  },
});

/**
 * Convert alarm list widget parameters to form
 *
 * @param {AlarmListWidgetParameters} [parameters = {}]
 * @return {AlarmListWidgetParametersForm}
 */
const alarmListWidgetParametersToForm = (parameters = {}) => ({
  ...parameters,
  ...alarmListWidgetDefaultParametersToForm(parameters),

  periodic_refresh: periodicRefreshToDurationForm(parameters.periodic_refresh),
  viewFilters: parameters.viewFilters || [],
  mainFilter: parameters.mainFilter || null,
  mainFilterUpdatedAt: parameters.mainFilterUpdatedAt || 0,
  liveReporting: parameters.liveReporting || {},
  sort: widgetSortToForm(parameters.sort),
  opened: isUndefined(parameters.opened) ? true : parameters.opened,
  expandGridRangeSize: parameters.expandGridRangeSize || [GRID_SIZES.min, GRID_SIZES.max],
  exportCsvSeparator: parameters.exportCsvSeparator || EXPORT_CSV_SEPARATORS.comma,
  exportCsvDatetimeFormat: parameters.exportCsvDatetimeFormat || EXPORT_CSV_DATETIME_FORMATS.datetimeSeconds.value,
  widgetExportColumns: parameters.widgetExportColumns
    ? widgetColumnsToForm(parameters.widgetExportColumns)
    : defaultColumnsToColumns(DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS),
});

/**
 * Convert alarm list widget to form object
 *
 * @param {AlarmListWidget} [alarmListWidget = {}]
 * @returns {AlarmListWidgetForm}
 */
export const alarmListWidgetToForm = (alarmListWidget = { type: WIDGET_TYPES.alarmList }) => {
  const widget = widgetToForm(alarmListWidget);

  return {
    ...widget,
    type: WIDGET_TYPES.alarmList,
    parameters: alarmListWidgetParametersToForm(alarmListWidget.parameters),
  };
};

/**
 * Generate alarm list widget with default parameters.
 *
 * @return {AlarmListWidgetForm}
 */
export const generateDefaultAlarmListWidget = () => alarmListWidgetToForm();

/**
 * Convert form sort parameters to widget sort
 *
 * @param {Sort} sort
 * @return {Sort}
 */
const formSortToWidgetSort = (sort = {}) => ({
  order: sort.order,
  column: formToColumnValue(sort.column),
});

/**
 * Convert form columns parameters to widget columns
 *
 * @param {WidgetColumn[]} widgetColumns
 * @return {WidgetColumn[]}
 */
export const formWidgetColumnsToColumns = widgetColumns => widgetColumns.map(column => ({
  ...column,
  value: formToColumnValue(column.value),
}));

/**
 * Convert infoPopups parameters to alarm list
 *
 * @param {InfoPopup[]} infoPopups
 * @return {InfoPopup[]}
 */
const formInfoPopupsToInfoPopups = infoPopups => infoPopups.map(infoPopup => ({
  ...infoPopup,
  column: columnValueToForm(infoPopup.column),
}));

/**
 * Convert alarm list settings form to alarm list object
 *
 * @param {AlarmListWidgetForm} [form = {}]
 * @returns {AlarmListWidget}
 */
export const formToAlarmListWidget = (form = {}) => {
  const { parameters } = form;

  return {
    ...form,
    parameters: {
      ...parameters,
      widgetColumns: formWidgetColumnsToColumns(parameters.widgetColumns),
      widgetGroupColumns: formWidgetColumnsToColumns(parameters.widgetGroupColumns),
      widgetExportColumns: formWidgetColumnsToColumns(parameters.widgetExportColumns),
      serviceDependenciesColumns: formWidgetColumnsToColumns(parameters.serviceDependenciesColumns),
      infoPopups: formInfoPopupsToInfoPopups(parameters.infoPopups),
      sort: formSortToWidgetSort(parameters.sort),
    },
  };
};
