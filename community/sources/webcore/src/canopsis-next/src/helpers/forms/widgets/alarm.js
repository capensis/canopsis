import { isString } from 'lodash';
import {
  DEFAULT_ALARMS_WIDGET_COLUMNS,
  DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS,
  DEFAULT_PERIODIC_REFRESH,
  DEFAULT_SERVICE_DEPENDENCIES_COLUMNS,
  EXPORT_CSV_DATETIME_FORMATS,
  EXPORT_CSV_SEPARATORS,
  GRID_SIZES,
  SORT_ORDERS,
  WIDGET_TYPES,
} from '@/constants';
import { DEFAULT_CATEGORIES_LIMIT, PAGINATION_LIMIT } from '@/config';

import { defaultColumnsToColumns } from '@/helpers/entities';
import { widgetToForm } from '@/helpers/forms/widgets/common';
import { durationWithEnabledToForm, formToDurationWithEnabled } from '@/helpers/date/duration';

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
 * @typedef {Object} AlarmsStateFilter
 * @property {boolean} opened
 * @property {boolean} resolved
 */

/**
 * @typedef {Object} Sort
 * @property {string} order
 * @property {string} column
 */

/**
 * @typedef {Object} AlarmListWidgetParameters
 * @property {FastAckOutput} fastAckOutput
 * @property {LinksCategoriesAsList} linksCategoriesAsList
 * @property {number} itemsPerPage
 * @property {Array} infoPopups
 * @property {string} moreInfoTemplate
 * @property {WidgetColumn[]} widgetColumns
 * @property {WidgetColumn[]} widgetGroupColumns
 * @property {WidgetColumn[]} serviceDependenciesColumns
 * @property {boolean} isAckNoteRequired
 * @property {boolean} isSnoozeNoteRequired
 * @property {boolean} isMultiAckEnabled
 * @property {boolean} isHtmlEnabledOnTimeLine
 * @property {DurationWithEnabled} periodic_refresh
 * @property {Array} viewFilters
 * @property {Object|null} mainFilter
 * @property {number} mainFilterUpdatedAt
 * @property {LiveReporting} liveReporting
 * @property {Sort} sort
 * @property {AlarmsStateFilter} alarmsStateFilter
 * @property {number} expandGridRangeSize
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
 * @typedef {AlarmListWidget} AlarmListWidgetForm
 * @property {AlarmListWidgetParametersForm} parameters
 */

/**
 * Prefix formatter for column value
 *
 * @param {string} [value]
 * @param {boolean} [isInitialization=false]
 * @returns {string}
 */
const columnValuePrefixFormatter = (value, isInitialization = false) => {
  if (isString(value) && value !== '') {
    if (isInitialization) {
      return value.replace(/^v\./, 'alarm.v.');
    }

    return value.replace(/^alarm\./, '');
  }

  return value;
};

/**
 * Convert columns parameters to form
 *
 * @param {WidgetColumn[]} widgetColumns
 * @param {boolean} [isInitialization]
 * @return {WidgetColumn[]}
 */
const widgetColumnsToForm = (widgetColumns, isInitialization) => widgetColumns.map(column => ({
  ...column,
  value: columnValuePrefixFormatter(column.value, isInitialization),
}));

/**
 * Convert alarm list infoPopups parameters to form
 *
 * @param {Array} infoPopups
 * @param {boolean} [isInitialization]
 * @return {Array}
 */
const infoPopupsToForm = (infoPopups, isInitialization) => infoPopups.map(infoPopup => ({
  ...infoPopup,
  column: columnValuePrefixFormatter(infoPopup.column, isInitialization),
}));

/**
 * Convert alarm list widget parameters to form
 *
 * @param {AlarmListWidgetParameters} parameters
 * @return {AlarmListWidgetParametersForm}
 */
const alarmListWidgetParametersToForm = (parameters = {}) => ({
  fastAckOutput: parameters.fastAckOutput || {
    enabled: false,
    value: 'auto ack',
  },
  linksCategoriesAsList: parameters.linksCategoriesAsList || {
    enabled: false,
    limit: DEFAULT_CATEGORIES_LIMIT,
  },
  itemsPerPage: parameters.itemsPerPage || PAGINATION_LIMIT,
  infoPopups: parameters.infoPopups ? infoPopupsToForm(parameters.infoPopups, true) : [],
  moreInfoTemplate: parameters.moreInfoTemplate || '',
  widgetColumns: parameters.widgetColumns
    ? widgetColumnsToForm(parameters.widgetColumns, true)
    : defaultColumnsToColumns(DEFAULT_ALARMS_WIDGET_COLUMNS),
  widgetGroupColumns: parameters.widgetGroupColumns
    ? widgetColumnsToForm(parameters.widgetGroupColumns, true)
    : defaultColumnsToColumns(DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS),
  serviceDependenciesColumns: parameters.serviceDependenciesColumns
    ? widgetColumnsToForm(parameters.serviceDependenciesColumns, true)
    : defaultColumnsToColumns(DEFAULT_SERVICE_DEPENDENCIES_COLUMNS),
  isAckNoteRequired: !!parameters.isAckNoteRequired,
  isSnoozeNoteRequired: !!parameters.isSnoozeNoteRequired,
  isMultiAckEnabled: !!parameters.isMultiAckEnabled,
  isHtmlEnabledOnTimeLine: !!parameters.isHtmlEnabledOnTimeLine,
  periodic_refresh: durationWithEnabledToForm(parameters.periodic_refresh || DEFAULT_PERIODIC_REFRESH),
  viewFilters: parameters.viewFilters || [],
  mainFilter: parameters.mainFilter || null,
  mainFilterUpdatedAt: parameters.mainFilterUpdatedAt || 0,
  liveReporting: parameters.liveReporting || {},
  sort: parameters.sort || {
    order: SORT_ORDERS.asc,
  },
  alarmsStateFilter: parameters.alarmsStateFilter || {
    opened: true,
  },
  expandGridRangeSize: parameters.expandGridRangeSize || [GRID_SIZES.min, GRID_SIZES.max],
  exportCsvSeparator: parameters.exportCsvSeparator || EXPORT_CSV_SEPARATORS.comma,
  exportCsvDatetimeFormat: parameters.exportCsvDatetimeFormat || EXPORT_CSV_DATETIME_FORMATS.datetimeSeconds,
});

/**
 * Convert alarm list widget to form object
 *
 * @param {AlarmListWidget} [alarmListWidget = {}]
 * @returns {AlarmListWidgetForm}
 */
export const alarmListWidgetToForm = (alarmListWidget = {}) => {
  const widget = widgetToForm(alarmListWidget);

  return {
    ...widget,
    type: WIDGET_TYPES.alarmList,
    parameters: alarmListWidgetParametersToForm(alarmListWidget.parameters),
  };
};

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
      widgetColumns: widgetColumnsToForm(parameters.widgetColumns),
      widgetGroupColumns: widgetColumnsToForm(parameters.widgetGroupColumns),
      serviceDependenciesColumns: widgetColumnsToForm(parameters.serviceDependenciesColumns),
      infoPopups: infoPopupsToForm(parameters.infoPopups),
      periodic_refresh: formToDurationWithEnabled(parameters.periodic_refresh),
    },
  };
};
