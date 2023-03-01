import { isString, cloneDeep } from 'lodash';

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
  ALARM_DENSE_TYPES,
} from '@/constants';
import { DEFAULT_CATEGORIES_LIMIT, PAGINATION_LIMIT } from '@/config';

import { defaultColumnsToColumns } from '@/helpers/entities';
import { durationWithEnabledToForm, isValidUnit } from '@/helpers/date/duration';

/**
 * @typedef {Object} AlarmsListDataTableColumn
 * @property {string} value
 * @property {string} text
 * @property {boolean} [isHtml]
 * @property {boolean} [colorIndicator]
 */

/**
 * @typedef {Object} WidgetFastAckOutput
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
 * @typedef {Object} AlarmListBaseParameters
 * @property {WidgetSort} sort
 * @property {number} itemsPerPage
 * @property {string} moreInfoTemplate
 * @property {WidgetInfoPopup[]} infoPopups
 * @property {WidgetColumn[]} widgetColumns
 */

/**
 * @typedef {Object} AlarmListWidgetDefaultParameters
 * @property {WidgetFastAckOutput} fastAckOutput
 * @property {WidgetLinksCategoriesAsList} linksCategoriesAsList
 * @property {number} itemsPerPage
 * @property {WidgetInfoPopup[]} infoPopups
 * @property {string} moreInfoTemplate
 * @property {WidgetColumn[]} widgetColumns
 * @property {WidgetColumn[]} widgetGroupColumns
 * @property {WidgetColumn[]} widgetExportColumns
 * @property {WidgetColumn[]} serviceDependenciesColumns
 * @property {boolean} isAckNoteRequired
 * @property {boolean} isSnoozeNoteRequired
 * @property {boolean} isMultiAckEnabled
 * @property {boolean} isMultiDeclareTicketEnabled
 * @property {boolean} isHtmlEnabledOnTimeLine
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
 * @param {WidgetInfoPopup[]} [infoPopups = []]
 * @return {WidgetInfoPopup[]}
 */
export const infoPopupsToForm = (infoPopups = []) => infoPopups.map(infoPopup => ({
  ...infoPopup,
  column: columnValueToForm(infoPopup.column),
}));

/**
 * Convert widget sort parameters to form
 *
 * @param {WidgetSort} [sort = {}]
 * @return {WidgetSort}
 */
export const widgetSortToForm = (sort = {}) => ({
  order: sort.order || SORT_ORDERS.asc,
  column: sort.column ? columnValueToForm(sort.column) : '',
});

/**
 * Convert alarm list base parameters (we are using it inside another widgets with alarmList) to form
 *
 * @param {AlarmListBaseParameters} [alarmListParameters = {}]
 * @return {AlarmListBaseParameters}
 */
export const alarmListBaseParametersToForm = (alarmListParameters = {}) => ({
  sort: widgetSortToForm(alarmListParameters.sort),
  itemsPerPage: alarmListParameters.itemsPerPage ?? PAGINATION_LIMIT,
  moreInfoTemplate: alarmListParameters.moreInfoTemplate ?? '',
  infoPopups: infoPopupsToForm(alarmListParameters.infoPopups),
  widgetColumns: alarmListParameters.widgetColumns
    ? widgetColumnsToForm(alarmListParameters.widgetColumns)
    : defaultColumnsToColumns(DEFAULT_ALARMS_WIDGET_COLUMNS),
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
 * @return {AlarmListWidgetDefaultParameters}
 */
export const alarmListWidgetDefaultParametersToForm = (parameters = {}) => ({
  itemsPerPage: parameters.itemsPerPage ?? PAGINATION_LIMIT,
  infoPopups: infoPopupsToForm(parameters.infoPopups),
  moreInfoTemplate: parameters.moreInfoTemplate ?? '',
  isAckNoteRequired: !!parameters.isAckNoteRequired,
  isSnoozeNoteRequired: !!parameters.isSnoozeNoteRequired,
  isMultiAckEnabled: !!parameters.isMultiAckEnabled,
  isMultiDeclareTicketEnabled: !!parameters.isMultiDeclareTicketEnabled,
  isHtmlEnabledOnTimeLine: !!parameters.isHtmlEnabledOnTimeLine,
  sticky_header: !!parameters.sticky_header,
  dense: parameters.dense ?? ALARM_DENSE_TYPES.large,
  fastAckOutput: parameters.fastAckOutput
    ? { ...parameters.fastAckOutput }
    : {
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
  linksCategoriesAsList: parameters.linksCategoriesAsList
    ? { ...parameters.linksCategoriesAsList }
    : {
      enabled: false,
      limit: DEFAULT_CATEGORIES_LIMIT,
    },
});

/**
 * Convert alarm list widget parameters to form
 *
 * @param {AlarmListWidgetParameters} [parameters = {}]
 * @return {AlarmListWidgetParameters}
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
  opened: parameters.opened ?? true,
  expandGridRangeSize: parameters.expandGridRangeSize
    ? [...parameters.expandGridRangeSize]
    : [GRID_SIZES.min, GRID_SIZES.max],
  exportCsvSeparator: parameters.exportCsvSeparator ?? EXPORT_CSV_SEPARATORS.comma,
  exportCsvDatetimeFormat: parameters.exportCsvDatetimeFormat ?? EXPORT_CSV_DATETIME_FORMATS.datetimeSeconds.value,
  widgetExportColumns: parameters.widgetExportColumns
    ? widgetColumnsToForm(parameters.widgetExportColumns)
    : defaultColumnsToColumns(DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS),
});

/**
 * Convert form sort parameters to widget sort
 *
 * @param {WidgetSort} sort
 * @return {WidgetSort}
 */
export const formSortToWidgetSort = (sort = {}) => ({
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
 * @param {WidgetInfoPopup[]} infoPopups
 * @return {WidgetInfoPopup[]}
 */
export const formInfoPopupsToInfoPopups = infoPopups => infoPopups.map(infoPopup => ({
  ...infoPopup,
  column: columnValueToForm(infoPopup.column),
}));

/**
 * Convert form to alarm list base parameters (we are using it inside another widgets with alarmList)
 *
 * @param {AlarmListBaseParameters} [form = {}]
 * @return {AlarmListBaseParameters}
 */
export const formToAlarmListBaseParameters = (form = {}) => ({
  sort: formSortToWidgetSort(form.sort),
  itemsPerPage: form.itemsPerPage,
  moreInfoTemplate: form.moreInfoTemplate,
  infoPopups: formInfoPopupsToInfoPopups(form.infoPopups),
  widgetColumns: formWidgetColumnsToColumns(form.widgetColumns),
});

/**
 * Convert form parameters to alarm list widget parameters
 *
 * @param {AlarmListWidgetParameters} form
 * @return {AlarmListWidgetParameters}
 */
export const formToAlarmListWidgetParameters = form => ({
  ...form,

  widgetColumns: formWidgetColumnsToColumns(form.widgetColumns),
  widgetGroupColumns: formWidgetColumnsToColumns(form.widgetGroupColumns),
  widgetExportColumns: formWidgetColumnsToColumns(form.widgetExportColumns),
  serviceDependenciesColumns: formWidgetColumnsToColumns(form.serviceDependenciesColumns),
  infoPopups: formInfoPopupsToInfoPopups(form.infoPopups),
  sort: formSortToWidgetSort(form.sort),
});

/**
 * Convert alarms list columns to data table columns
 *
 * @param {WidgetColumn[]} [columns = []]
 * @returns {AlarmsListDataTableColumn[]}
 */
export const alarmsListColumnsToTableColumns = (columns = []) => columns.map(({ label, ...column }) => ({
  ...column,

  text: label,
}));
