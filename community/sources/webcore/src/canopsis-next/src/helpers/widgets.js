import {
  ALARM_FIELDS,
  ALARM_FIELDS_TO_LABELS_KEYS,
  ALARM_UNSORTABLE_FIELDS,
  DEFAULT_ALARMS_WIDGET_COLUMNS,
  DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS,
  DEFAULT_CONTEXT_WIDGET_COLUMNS,
  DEFAULT_SERVICE_DEPENDENCIES_COLUMNS,
  ENTITY_FIELDS_TO_LABELS_KEYS,
  ENTITY_UNSORTABLE_FIELDS,
} from '@/constants';

import i18n from '@/i18n';

import { setSeveralFields } from './immutable';
import { getInfosWidgetColumn, isLinksWidgetColumn } from './forms/shared/widget-column';

/**
 * Get translated label for widget column
 *
 * @param {WidgetColumn} column
 * @param {Object<string, string>} labelsMap
 * @returns {string}
 */
export const getColumnLabel = (column, labelsMap) => {
  if (column.label) {
    return column.label;
  }

  const infosColumn = getInfosWidgetColumn(column.value);

  if (infosColumn) {
    return i18n.tc(labelsMap[infosColumn], 2);
  }

  if (isLinksWidgetColumn(column.value)) {
    return i18n.tc(labelsMap[ALARM_FIELDS.links], 2);
  }

  return i18n.tc(labelsMap[column.value], 2);
};

/**
 * Get sortable property for widget column for table
 *
 * @param {WidgetColumn} column
 * @param {string[]} unsortableFields
 * @returns {boolean}
 */
export const getSortable = (column, unsortableFields = []) => (
  unsortableFields.some(field => column.value.startsWith(field))
);

/**
 * Prepared alarms list widget for displaying
 *
 * @param {Object} widget
 * @returns {Object}
 */
export const prepareAlarmListWidget = (widget = {}) => setSeveralFields(widget, {
  'parameters.widgetColumns': (columns = []) => (
    columns.map(column => ({
      ...column,

      sortable: getSortable(column, ALARM_UNSORTABLE_FIELDS),
      text: getColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
    }))
  ),

  'parameters.widgetGroupColumns': (columns = DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS) => (
    columns.map(column => ({
      ...column,

      sortable: getSortable(column, ALARM_UNSORTABLE_FIELDS),
      text: getColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
    }))
  ),

  'parameters.serviceDependenciesColumns': (columns = DEFAULT_SERVICE_DEPENDENCIES_COLUMNS) => (
    columns.map(column => ({
      ...column,

      sortable: false,
      text: getColumnLabel(column, ENTITY_FIELDS_TO_LABELS_KEYS),
      value: column.value.startsWith('entity.') ? column.value : `entity.${column.value}`,
    }))
  ),

  'parameters.widgetExportColumns': (columns = []) => (
    columns.map(column => ({
      ...column,

      text: getColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
    }))
  ),
});

/**
 * Prepared context widget for displaying
 *
 * @param {Object} widget
 * @returns {Object}
 */
export const prepareContextWidget = (widget = {}) => setSeveralFields(widget, {
  'parameters.widgetColumns': (columns = []) => (
    columns.map(column => ({
      ...column,

      sortable: getSortable(column, ENTITY_UNSORTABLE_FIELDS),
      text: getColumnLabel(column, ENTITY_FIELDS_TO_LABELS_KEYS),
    }))
  ),

  'parameters.serviceDependenciesColumns': (columns = DEFAULT_SERVICE_DEPENDENCIES_COLUMNS) => (
    columns.map(column => ({
      ...column,

      sortable: false,
      text: getColumnLabel(column, ENTITY_FIELDS_TO_LABELS_KEYS),
      value: column.value.startsWith('entity.') ? column.value : `entity.${column.value}`,
    }))
  ),

  'parameters.activeAlarmsColumns': (columns = DEFAULT_ALARMS_WIDGET_COLUMNS) => (
    columns.map(column => ({
      ...column,

      sortable: getSortable(column, ALARM_UNSORTABLE_FIELDS),
      text: getColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
    }))
  ),

  'parameters.resolvedAlarmsColumns': (columns = DEFAULT_ALARMS_WIDGET_COLUMNS) => (
    columns.map(column => ({
      ...column,

      sortable: getSortable(column, ALARM_UNSORTABLE_FIELDS),
      text: getColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
    }))
  ),

  'parameters.widgetExportColumns': (columns = []) => (
    columns.map(column => ({
      ...column,

      text: getColumnLabel(column, ENTITY_FIELDS_TO_LABELS_KEYS),
    }))
  ),
});

/**
 * Prepared service weather widget for displaying
 *
 * @param {Object} widget
 * @returns {Object}
 */
export const prepareServiceWeatherWidget = (widget = {}) => setSeveralFields(widget, {
  'parameters.serviceDependenciesColumns': (columns = DEFAULT_SERVICE_DEPENDENCIES_COLUMNS) => (
    columns.map(column => ({
      ...column,

      sortable: false,
      text: getColumnLabel(column, ENTITY_FIELDS_TO_LABELS_KEYS),
      value: column.value.startsWith('entity.') ? column.value : `entity.${column.value}`,
    }))
  ),

  'parameters.alarmsList.widgetColumns': (columns = DEFAULT_ALARMS_WIDGET_COLUMNS) => (
    columns.map(column => ({
      ...column,

      sortable: getSortable(column, ALARM_UNSORTABLE_FIELDS),
      text: getColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
    }))
  ),
});

/**
 * Prepared stats calendar/counter widget for displaying
 *
 * @param {Object} widget
 * @returns {Object}
 */
export const prepareStatsCalendarAndCounterWidget = (widget = {}) => setSeveralFields(widget, {
  'parameters.alarmsList.widgetColumns': (columns = DEFAULT_ALARMS_WIDGET_COLUMNS) => (
    columns.map(column => ({
      ...column,

      sortable: getSortable(column, ALARM_UNSORTABLE_FIELDS),
      text: getColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
    }))
  ),
});

/**
 * Prepared map widget for displaying
 *
 * @param {Object} widget
 * @returns {Object}
 */
export const prepareMapWidget = (widget = {}) => setSeveralFields(widget, {
  'parameters.alarmsColumns': (columns = DEFAULT_ALARMS_WIDGET_COLUMNS) => (
    columns.map(column => ({
      ...column,

      sortable: getSortable(column, ALARM_UNSORTABLE_FIELDS),
      text: getColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
    }))
  ),

  'parameters.entitiesColumns': (columns = DEFAULT_CONTEXT_WIDGET_COLUMNS) => (
    columns.map(column => ({
      ...column,

      sortable: getSortable(column, ENTITY_UNSORTABLE_FIELDS),
      text: getColumnLabel(column, ENTITY_FIELDS_TO_LABELS_KEYS),
    }))
  ),
});
