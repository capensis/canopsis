import {
  AVAILABILITY_DISPLAY_PARAMETERS,
  AVAILABILITY_SHOW_TYPE,
  CSV_SEPARATORS,
  DEFAULT_ALARMS_WIDGET_COLUMNS,
  DEFAULT_CONTEXT_WIDGET_COLUMNS,
  DEFAULT_CONTEXT_WIDGET_RESOLVED_ALARMS_COLUMNS,
  DEFAULT_PERIODIC_REFRESH,
  ENTITY_FIELDS,
  ENTITY_FIELDS_TO_LABELS_KEYS,
  ENTITY_UNSORTABLE_FIELDS,
  ALARM_FIELDS_TO_LABELS_KEYS,
  ALARM_UNSORTABLE_FIELDS,
  QUICK_RANGES,
} from '@/constants';

import { durationWithEnabledToForm } from '@/helpers/date/duration';
import { setSeveralFields } from '@/helpers/immutable';

import { formToWidgetColumns, widgetColumnsToForm } from '../column/form';
import { getWidgetColumnLabel, getWidgetColumnSortable } from '../list';
import { formToWidgetTemplateValue } from '../template/form';

/**
 * @typedef {Object} AvailabilityField
 * @property {boolean} enabled
 * @property {number} default_show_type
 * @property {string} default_time_range
 */

/**
 * Converts the availability field to form values.
 *
 * @param {AvailabilityField} [availability={}]
 * @returns {AvailabilityField}
 */
export const availabilityFieldToForm = (availability = {}) => ({
  enabled: availability.enabled ?? false,
  default_show_type: availability.default_show_type ?? AVAILABILITY_SHOW_TYPE.percent,
  default_time_range: availability.default_time_range ?? QUICK_RANGES.today.value,
});

/**
 * Converts availability widget parameters to form values.
 *
 * @param {Object} [parameters={}]
 * @returns {Object}
 */
export const availabilityWidgetParametersToForm = (parameters = {}) => ({
  mainFilter: parameters.mainFilter ?? null,
  periodic_refresh: durationWithEnabledToForm(
    parameters.periodic_refresh ?? DEFAULT_PERIODIC_REFRESH,
  ),
  widget_columns: widgetColumnsToForm(parameters.widget_columns ?? DEFAULT_CONTEXT_WIDGET_COLUMNS),
  widget_columns_template: formToWidgetTemplateValue(parameters.widget_columns_template),
  active_alarms_columns: widgetColumnsToForm(
    parameters.active_alarms_columns ?? DEFAULT_ALARMS_WIDGET_COLUMNS,
  ),
  active_alarms_columns_template: formToWidgetTemplateValue(
    parameters.active_alarms_columns_template,
  ),
  resolved_alarms_columns: widgetColumnsToForm(
    parameters.resolved_alarms_columns ?? DEFAULT_CONTEXT_WIDGET_RESOLVED_ALARMS_COLUMNS,
  ),
  resolved_alarms_columns_template: formToWidgetTemplateValue(
    parameters.resolved_alarms_columns_template,
  ),
  default_time_range: parameters.default_time_range ?? QUICK_RANGES.today.value,
  default_display_parameter:
    parameters.default_display_parameter ?? AVAILABILITY_DISPLAY_PARAMETERS.uptime,
  default_show_type: parameters.default_show_type ?? AVAILABILITY_SHOW_TYPE.percent,
  export_settings: {
    widgetExportColumns: widgetColumnsToForm(
      parameters.export_settings?.widget_export_columns ?? DEFAULT_CONTEXT_WIDGET_COLUMNS,
    ),
    widgetExportColumnsTemplate: formToWidgetTemplateValue(
      parameters.export_settings?.widget_export_columns_template,
    ),
    exportCsvSeparator: parameters.export_settings?.export_csv_separator ?? CSV_SEPARATORS.comma,
  },
});

/**
 * Converts availability form values to widget parameters.
 *
 * @param {Object} form
 * @returns {Object}
 */
export const formToAvailabilityWidgetParameters = (form = {}) => ({
  ...form,
  widget_columns: formToWidgetColumns(form.widget_columns),
  widget_columns_template: formToWidgetTemplateValue(form.widget_columns_template),
  active_alarms_columns: formToWidgetColumns(form.active_alarms_columns),
  active_alarms_columns_template: formToWidgetTemplateValue(form.active_alarms_columns_template),
  resolved_alarms_columns: formToWidgetColumns(form.resolved_alarms_columns),
  resolved_alarms_columns_template: formToWidgetTemplateValue(
    form.resolved_alarms_columns_template,
  ),
  export_settings: {
    ...form.export_settings,
    widget_export_columns: formToWidgetColumns(form.export_settings?.widgetExportColumns),
    widget_export_columns_template: formToWidgetTemplateValue(
      form.export_settings?.widgetExportColumnsTemplate,
    ),
    export_csv_separator: form.export_settings?.exportCsvSeparator,
  },
});

/**
 * Adds display metadata to availability widget columns.
 *
 * @param {Object} widget
 * @returns {Object}
 */
export const prepareAvailabilityWidget = (widget = {}) => setSeveralFields(widget, {
    'parameters.widget_columns': (columns = []) => columns.map(column => ({
        ...column,
        value: column.value === ENTITY_FIELDS.categoryName ? ENTITY_FIELDS.category : column.value,
        sortable: getWidgetColumnSortable(column, ENTITY_UNSORTABLE_FIELDS),
        text: getWidgetColumnLabel(column, ENTITY_FIELDS_TO_LABELS_KEYS),
      })),
    'parameters.active_alarms_columns': (columns = DEFAULT_ALARMS_WIDGET_COLUMNS) => columns.map(column => ({
        ...column,
        sortable: getWidgetColumnSortable(column, ALARM_UNSORTABLE_FIELDS),
        text: getWidgetColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
      })),
    'parameters.resolved_alarms_columns': (columns = DEFAULT_ALARMS_WIDGET_COLUMNS) => columns.map(column => ({
        ...column,
        sortable: getWidgetColumnSortable(column, ALARM_UNSORTABLE_FIELDS),
        text: getWidgetColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
      })),
    'parameters.export_settings.widget_export_columns': (columns = []) => columns.map(column => ({
        ...column,
        text: getWidgetColumnLabel(column, ENTITY_FIELDS_TO_LABELS_KEYS),
      })),
  });
