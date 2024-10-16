import {
  ENTITY_UNSORTABLE_FIELDS,
  AVAILABILITY_DISPLAY_PARAMETERS,
  AVAILABILITY_SHOW_TYPE,
  DEFAULT_ALARMS_WIDGET_COLUMNS,
  DEFAULT_CONTEXT_WIDGET_COLUMNS,
  DEFAULT_PERIODIC_REFRESH,
  ENTITY_FIELDS_TO_LABELS_KEYS,
  QUICK_RANGES,
  ALARM_UNSORTABLE_FIELDS,
  ALARM_FIELDS_TO_LABELS_KEYS,
  DEFAULT_CONTEXT_WIDGET_RESOLVED_ALARMS_COLUMNS,
  EXPORT_CSV_SEPARATORS,
  ENTITY_FIELDS,
} from '@/constants';

import { durationWithEnabledToForm } from '@/helpers/date/duration';
import { formToWidgetTemplateValue } from '@/helpers/entities/widget/template/form';
import { formToWidgetColumns, widgetColumnsToForm } from '@/helpers/entities/widget/column/form';
import { setSeveralFields } from '@/helpers/immutable';
import { getWidgetColumnLabel, getWidgetColumnSortable } from '@/helpers/entities/widget/list';

/**
 * @typedef {Object} AvailabilityWidgetParametersExportSettings
 * @property {WidgetColumn[]} widget_export_columns
 * @property {string} widget_export_columns_template
 * @property {string} export_csv_separator
 */

/**
 * @typedef {Object} AvailabilityWidgetParametersExportSettingsForm
 * @property {WidgetColumnForm[]} widgetExportColumns
 * @property {string | Symbol} widgetExportColumnsTemplate
 * @property {string} exportCsvSeparator
 */

/**
 * @typedef {Object} AvailabilityWidgetParameters
 * @property {DurationWithEnabled} periodic_refresh
 * @property {string | null} mainFilter
 * @property {WidgetColumn[]} widget_columns
 * @property {string} widget_columns_template
 * @property {WidgetColumn[]} active_alarms_columns
 * @property {string} active_alarms_columns_template
 * @property {WidgetColumn[]} resolved_alarms_columns
 * @property {string} resolved_alarms_columns_template
 * @property {number} default_show_type
 * @property {string} default_time_range
 * @property {number} default_display_parameter
 * @property {AvailabilityWidgetParametersExportSettings} export_settings
 */

/**
 * @typedef {Object} AvailabilityField
 * @property {boolean} enabled
 * @property {number} default_show_type
 * @property {string} default_time_range
 */

/**
 * @typedef {AvailabilityWidgetParameters} AvailabilityWidgetParametersForm
 * @property {WidgetColumnForm[]} widget_columns
 * @property {string | Symbol} widget_columns_template
 * @property {WidgetColumnForm[]} active_alarms_columns
 * @property {string | Symbol} active_alarms_columns_template
 * @property {WidgetColumnForm[]} resolved_alarms_columns
 * @property {string | Symbol} resolved_alarms_columns_template
 * @property {AvailabilityWidgetParametersExportSettingsForm} export_settings
 */

/**
 * Convert availability widget field to form
 *
 * @param {AvailabilityField} [availability = {}]
 * @return {AvailabilityField}
 */
export const availabilityFieldToForm = (availability = {}) => ({
  enabled: availability.enabled ?? false,
  default_show_type: availability.default_show_type ?? AVAILABILITY_SHOW_TYPE.percent,
  default_time_range: availability.default_time_range ?? QUICK_RANGES.today.value,
});

/**
 * Convert availability widget field to form
 *
 * @param {AvailabilityWidgetParametersExportSettings} [exportSettings = {}]
 * @return {AvailabilityWidgetParametersExportSettingsForm}
 */
export const availabilityWidgetParametersExportSettingsToForm = (exportSettings = {}) => ({
  widgetExportColumns:
    widgetColumnsToForm(exportSettings.widget_export_columns ?? DEFAULT_CONTEXT_WIDGET_COLUMNS),
  widgetExportColumnsTemplate: formToWidgetTemplateValue(exportSettings.widget_export_columns_template),
  exportCsvSeparator: exportSettings.export_csv_separator ?? EXPORT_CSV_SEPARATORS.comma,
});

/**
 * Convert form to availability widget parameters to form
 *
 * @param {AvailabilityWidgetParameters} [parameters = {}]
 * @returns {AvailabilityWidgetParametersForm}
 */
export const availabilityWidgetParametersToForm = parameters => ({
  mainFilter: parameters.mainFilter ?? null,
  periodic_refresh: durationWithEnabledToForm(parameters.periodic_refresh ?? DEFAULT_PERIODIC_REFRESH),
  widget_columns:
    widgetColumnsToForm(parameters.widget_columns ?? DEFAULT_CONTEXT_WIDGET_COLUMNS),
  widget_columns_template: formToWidgetTemplateValue(parameters.widget_columns_template),
  active_alarms_columns:
    widgetColumnsToForm(parameters.active_alarms_columns ?? DEFAULT_ALARMS_WIDGET_COLUMNS),
  active_alarms_columns_template: formToWidgetTemplateValue(parameters.active_alarms_columns_template),
  resolved_alarms_columns:
    widgetColumnsToForm(parameters.resolved_alarms_columns ?? DEFAULT_CONTEXT_WIDGET_RESOLVED_ALARMS_COLUMNS),
  resolved_alarms_columns_template: formToWidgetTemplateValue(parameters.resolved_alarms_columns_template),
  default_time_range: parameters.default_time_range ?? QUICK_RANGES.today.value,
  default_display_parameter: parameters.default_display_parameter ?? AVAILABILITY_DISPLAY_PARAMETERS.uptime,
  default_show_type: parameters.default_show_type ?? AVAILABILITY_SHOW_TYPE.percent,
  export_settings: availabilityWidgetParametersExportSettingsToForm(parameters.export_settings),
});

/**
 * Convert form export setting availability widget parameters to form
 *
 * @param {AvailabilityWidgetParametersExportSettingsForm} form
 * @returns {AvailabilityWidgetParametersExportSettings}
 */
export const formToAvailabilityWidgetParametersExportSettings = form => ({
  export_csv_separator: form.exportCsvSeparator,
  widget_export_columns: formToWidgetColumns(form.widgetExportColumns),
  widget_export_columns_template: formToWidgetTemplateValue(form.widgetExportColumnsTemplate),
});

/**
 * Convert form to statistics widget parameters to form
 *
 * @param {AvailabilityWidgetParametersForm} form
 * @returns {StatisticsWidgetParameters}
 */
export const formToAvailabilityWidgetParameters = form => ({
  ...form,
  widget_columns: formToWidgetColumns(form.widget_columns),
  widget_columns_template: formToWidgetTemplateValue(form.widget_columns_template),
  active_alarms_columns: formToWidgetColumns(form.active_alarms_columns),
  active_alarms_columns_template: formToWidgetTemplateValue(form.active_alarms_columns_template),
  resolved_alarms_columns: formToWidgetColumns(form.resolved_alarms_columns),
  resolved_alarms_columns_template: formToWidgetTemplateValue(form.resolved_alarms_columns_template),
  export_settings: formToAvailabilityWidgetParametersExportSettings(form.export_settings),
});

/**
 * Prepared availability widget for displaying
 *
 * @param {Object} widget
 * @returns {Object}
 */
export const prepareAvailabilityWidget = (widget = {}) => setSeveralFields(widget, {
  'parameters.widget_columns': (columns = []) => (
    columns.map(column => ({
      ...column,

      value: column.value === ENTITY_FIELDS.categoryName ? ENTITY_FIELDS.category : column.value,
      sortable: getWidgetColumnSortable(column, ENTITY_UNSORTABLE_FIELDS),
      text: getWidgetColumnLabel(column, ENTITY_FIELDS_TO_LABELS_KEYS),
    }))
  ),

  'parameters.active_alarms_columns': (columns = DEFAULT_ALARMS_WIDGET_COLUMNS) => (
    columns.map(column => ({
      ...column,

      sortable: getWidgetColumnSortable(column, ALARM_UNSORTABLE_FIELDS),
      text: getWidgetColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
    }))
  ),

  'parameters.resolved_alarms_columns': (columns = DEFAULT_ALARMS_WIDGET_COLUMNS) => (
    columns.map(column => ({
      ...column,

      sortable: getWidgetColumnSortable(column, ALARM_UNSORTABLE_FIELDS),
      text: getWidgetColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
    }))
  ),

  'parameters.export_settings.widget_export_columns': (columns = []) => (
    columns.map(column => ({
      ...column,

      text: getWidgetColumnLabel(column, ENTITY_FIELDS_TO_LABELS_KEYS),
    }))
  ),
});
