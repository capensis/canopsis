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
} from '@/constants';

import {
  widgetColumnsToForm,
  formToWidgetColumns,
  widgetColumnTemplateToForm,
  formToWidgetColumnTemplate,
} from '../shared/widget-column';

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
 * @property {WidgetSort} sort
 * @property {WidgetCsvSeparator} exportCsvSeparator
 * @property {string} exportCsvDatetimeFormat
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
  widgetColumnsTemplate: widgetColumnTemplateToForm(parameters.widgetColumnsTemplate),
  serviceDependenciesColumnsTemplate: widgetColumnTemplateToForm(parameters.serviceDependenciesColumnsTemplate),
  widgetExportColumnsTemplate: widgetColumnTemplateToForm(parameters.widgetExportColumnsTemplate),
  resolvedAlarmsColumnsTemplate: widgetColumnTemplateToForm(parameters.resolvedAlarmsColumnsTemplate),
  activeAlarmsColumnsTemplate: widgetColumnTemplateToForm(parameters.activeAlarmsColumnsTemplate),
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
  sort: parameters.sort ? { ...parameters.sort } : { order: SORT_ORDERS.asc },
  exportCsvSeparator: parameters.exportCsvSeparator ?? EXPORT_CSV_SEPARATORS.comma,
  exportCsvDatetimeFormat: parameters.exportCsvDatetimeFormat ?? EXPORT_CSV_DATETIME_FORMATS.datetimeSeconds.value,
});

/**
 * Convert form to context widget parameters
 *
 * @param {ContextWidgetParametersForm} form
 * @returns {ContextWidgetParameters}
 */
export const formToContextWidgetParameters = form => ({
  ...form,

  widgetColumnsTemplate: formToWidgetColumnTemplate(form.widgetColumnsTemplate),
  serviceDependenciesColumnsTemplate: formToWidgetColumnTemplate(form.serviceDependenciesColumnsTemplate),
  widgetExportColumnsTemplate: formToWidgetColumnTemplate(form.widgetExportColumnsTemplate),
  resolvedAlarmsColumnsTemplate: formToWidgetColumnTemplate(form.resolvedAlarmsColumnsTemplate),
  activeAlarmsColumnsTemplate: formToWidgetColumnTemplate(form.activeAlarmsColumnsTemplate),
  widgetColumns: formToWidgetColumns(form.widgetColumns),
  serviceDependenciesColumns: formToWidgetColumns(form.serviceDependenciesColumns),
  resolvedAlarmsColumns: formToWidgetColumns(form.resolvedAlarmsColumns),
  activeAlarmsColumns: formToWidgetColumns(form.activeAlarmsColumns),
  widgetExportColumns: formToWidgetColumns(form.widgetExportColumns),
});
