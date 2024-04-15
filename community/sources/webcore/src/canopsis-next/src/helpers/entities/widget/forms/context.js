import { cloneDeep } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';
import {
  DEFAULT_CONTEXT_WIDGET_COLUMNS,
  DEFAULT_CONTEXT_WIDGET_RESOLVED_ALARMS_COLUMNS,
  DEFAULT_CONTEXT_WIDGET_ACTIVE_ALARM_COLUMNS,
  DEFAULT_SERVICE_DEPENDENCIES_COLUMNS,
  EXPORT_CSV_DATETIME_FORMATS,
  EXPORT_CSV_SEPARATORS,
  SORT_ORDERS,
  ENTITY_UNSORTABLE_FIELDS,
  ENTITY_FIELDS_TO_LABELS_KEYS,
  DEFAULT_ALARMS_WIDGET_COLUMNS,
  ALARM_UNSORTABLE_FIELDS,
  ALARM_FIELDS_TO_LABELS_KEYS,
  TREE_OF_DEPENDENCIES_SHOW_TYPES,
  COLOR_INDICATOR_TYPES,
  GRID_SIZES,
} from '@/constants';

import { addKeyInEntities, removeKeyFromEntities } from '@/helpers/array';
import { setSeveralFields } from '@/helpers/immutable';
import { availabilityFieldToForm } from '@/helpers/entities/widget/forms/availability';

import { widgetColumnsToForm, formToWidgetColumns } from '../column/form';
import { widgetTemplateValueToForm, formToWidgetTemplateValue } from '../template/form';
import { getWidgetColumnLabel, getWidgetColumnSortable } from '../list';

/**
 * @typedef {Object} ContextWidgetParameters
 * @property {string | null} mainFilter
 * @property {number} itemsPerPage
 * @property {string} widgetColumnsTemplate
 * @property {string} serviceDependenciesColumnsTemplate
 * @property {string} widgetExportColumnsTemplate
 * @property {string} resolvedAlarmsColumnsTemplate
 * @property {string} activeAlarmsColumnsTemplate
 * @property {WidgetColumn[]} widgetColumns
 * @property {WidgetColumn[]} serviceDependenciesColumns
 * @property {WidgetColumn[]} widgetExportColumns
 * @property {WidgetColumn[]} resolvedAlarmsColumns
 * @property {WidgetColumn[]} activeAlarmsColumns
 * @property {string[]} selectedTypes
 * @property {number[]} expandGridRangeSize
 * @property {WidgetSort} sort
 * @property {WidgetCsvSeparator} exportCsvSeparator
 * @property {string} exportCsvDatetimeFormat
 * @property {AlarmChart[]} charts
 * @property {boolean} showRootCauseByStateClick
 * @property {ColorIndicator} rootCauseColorIndicator
 * @property {number} treeOfDependenciesShowType
 * @property {AvailabilityField} availability
 */

/**
 * @typedef {ContextWidgetParameters} ContextWidgetParametersForm
 * @property {string | Symbol} widgetColumnsTemplate
 * @property {string | Symbol} serviceDependenciesColumnsTemplate
 * @property {string | Symbol} widgetExportColumnsTemplate
 * @property {string | Symbol} resolvedAlarmsColumnsTemplate
 * @property {string | Symbol} activeAlarmsColumnsTemplate
 * @property {WidgetColumnForm[]} widgetColumns
 * @property {WidgetColumnForm[]} serviceDependenciesColumns
 * @property {WidgetColumnForm[]} widgetExportColumns
 * @property {WidgetColumnForm[]} resolvedAlarmsColumns
 * @property {WidgetColumnForm[]} activeAlarmsColumns
 * @property {AlarmChartForm[]} charts
 */

/**
 * Convert context widget parameters to form
 *
 * @param {ContextWidgetParameters | {}} [parameters = {}]
 * @return {ContextWidgetParametersForm}
 */
export const contextWidgetParametersToForm = (parameters = {}) => ({
  mainFilter: parameters.mainFilter ?? null,
  itemsPerPage: parameters.itemsPerPage ?? PAGINATION_LIMIT,
  widgetColumnsTemplate: widgetTemplateValueToForm(parameters.widgetColumnsTemplate),
  serviceDependenciesColumnsTemplate: widgetTemplateValueToForm(parameters.serviceDependenciesColumnsTemplate),
  widgetExportColumnsTemplate: widgetTemplateValueToForm(parameters.widgetExportColumnsTemplate),
  resolvedAlarmsColumnsTemplate: widgetTemplateValueToForm(parameters.resolvedAlarmsColumnsTemplate),
  activeAlarmsColumnsTemplate: widgetTemplateValueToForm(parameters.activeAlarmsColumnsTemplate),
  widgetColumns:
    widgetColumnsToForm(parameters.widgetColumns ?? DEFAULT_CONTEXT_WIDGET_COLUMNS),
  serviceDependenciesColumns:
    widgetColumnsToForm(parameters.serviceDependenciesColumns ?? DEFAULT_SERVICE_DEPENDENCIES_COLUMNS),
  resolvedAlarmsColumns:
    widgetColumnsToForm(parameters.resolvedAlarmsColumns ?? DEFAULT_CONTEXT_WIDGET_RESOLVED_ALARMS_COLUMNS),
  activeAlarmsColumns:
    widgetColumnsToForm(parameters.activeAlarmsColumns ?? DEFAULT_CONTEXT_WIDGET_ACTIVE_ALARM_COLUMNS),
  widgetExportColumns:
    widgetColumnsToForm(parameters.widgetExportColumns ?? DEFAULT_CONTEXT_WIDGET_COLUMNS),
  selectedTypes: parameters.selectedTypes
    ? cloneDeep(parameters.selectedTypes)
    : [],
  expandGridRangeSize: parameters.expandGridRangeSize
    ? [...parameters.expandGridRangeSize]
    : [GRID_SIZES.min, GRID_SIZES.max],
  sort: parameters.sort ? { ...parameters.sort } : { order: SORT_ORDERS.asc },
  exportCsvSeparator: parameters.exportCsvSeparator ?? EXPORT_CSV_SEPARATORS.comma,
  exportCsvDatetimeFormat: parameters.exportCsvDatetimeFormat ?? EXPORT_CSV_DATETIME_FORMATS.datetimeSeconds.value,
  charts: addKeyInEntities(parameters.charts),
  showRootCauseByStateClick: parameters.showRootCauseByStateClick ?? true,
  rootCauseColorIndicator: parameters.rootCauseColorIndicator ?? COLOR_INDICATOR_TYPES.state,
  treeOfDependenciesShowType: parameters.treeOfDependenciesShowType ?? TREE_OF_DEPENDENCIES_SHOW_TYPES.custom,
  availability: availabilityFieldToForm(parameters.availability),
});

/**
 * Convert form to context widget parameters
 *
 * @param {ContextWidgetParametersForm} form
 * @returns {ContextWidgetParameters}
 */
export const formToContextWidgetParameters = form => ({
  ...form,

  widgetColumnsTemplate: formToWidgetTemplateValue(form.widgetColumnsTemplate),
  serviceDependenciesColumnsTemplate: formToWidgetTemplateValue(form.serviceDependenciesColumnsTemplate),
  widgetExportColumnsTemplate: formToWidgetTemplateValue(form.widgetExportColumnsTemplate),
  resolvedAlarmsColumnsTemplate: formToWidgetTemplateValue(form.resolvedAlarmsColumnsTemplate),
  activeAlarmsColumnsTemplate: formToWidgetTemplateValue(form.activeAlarmsColumnsTemplate),
  widgetColumns: formToWidgetColumns(form.widgetColumns),
  serviceDependenciesColumns: formToWidgetColumns(form.serviceDependenciesColumns),
  resolvedAlarmsColumns: formToWidgetColumns(form.resolvedAlarmsColumns),
  activeAlarmsColumns: formToWidgetColumns(form.activeAlarmsColumns),
  widgetExportColumns: formToWidgetColumns(form.widgetExportColumns),
  charts: removeKeyFromEntities(form.charts),
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

      sortable: getWidgetColumnSortable(column, ENTITY_UNSORTABLE_FIELDS),
      text: getWidgetColumnLabel(column, ENTITY_FIELDS_TO_LABELS_KEYS),
    }))
  ),

  'parameters.serviceDependenciesColumns': (columns = DEFAULT_SERVICE_DEPENDENCIES_COLUMNS) => (
    columns.map(column => ({
      ...column,

      sortable: false,
      text: getWidgetColumnLabel(column, ENTITY_FIELDS_TO_LABELS_KEYS),
    }))
  ),

  'parameters.activeAlarmsColumns': (columns = DEFAULT_ALARMS_WIDGET_COLUMNS) => (
    columns.map(column => ({
      ...column,

      sortable: getWidgetColumnSortable(column, ALARM_UNSORTABLE_FIELDS),
      text: getWidgetColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
    }))
  ),

  'parameters.resolvedAlarmsColumns': (columns = DEFAULT_ALARMS_WIDGET_COLUMNS) => (
    columns.map(column => ({
      ...column,

      sortable: getWidgetColumnSortable(column, ALARM_UNSORTABLE_FIELDS),
      text: getWidgetColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
    }))
  ),

  'parameters.widgetExportColumns': (columns = []) => (
    columns.map(column => ({
      ...column,

      text: getWidgetColumnLabel(column, ENTITY_FIELDS_TO_LABELS_KEYS),
    }))
  ),
});
