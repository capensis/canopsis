import { cloneDeep } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';
import {
  DEFAULT_ALARMS_WIDGET_COLUMNS,
  DEFAULT_CONTEXT_WIDGET_COLUMNS,
  DEFAULT_SERVICE_DEPENDENCIES_COLUMNS,
  EXPORT_CSV_DATETIME_FORMATS,
  EXPORT_CSV_SEPARATORS,
} from '@/constants';

import { defaultColumnsToColumns } from '@/helpers/entities';

import { widgetSortToForm, formSortToWidgetSort } from './alarm';

export const widgetFiltersParametersToFormParameters = (parameters = {}) => ({
  filters: parameters.filters
    ? cloneDeep(parameters.filters)
    : [],
  main_filter: parameters.main_filter ?? null,
});

export const contextWidgetParametersToFormParameters = (parameters = {}) => ({
  ...widgetFiltersParametersToFormParameters(parameters),

  itemsPerPage: parameters.itemsPerPage ?? PAGINATION_LIMIT,
  widgetColumns: parameters.widgetColumns
    ? cloneDeep(parameters.widgetColumns)
    : defaultColumnsToColumns(DEFAULT_CONTEXT_WIDGET_COLUMNS),
  serviceDependenciesColumns: parameters.serviceDependenciesColumns
    ? cloneDeep(parameters.serviceDependenciesColumns)
    : defaultColumnsToColumns(DEFAULT_SERVICE_DEPENDENCIES_COLUMNS),
  selectedTypes: parameters.selectedTypes
    ? cloneDeep(parameters.selectedTypes)
    : [],
  sort: widgetSortToForm(parameters.sort), // TODO: move to common
  exportCsvSeparator: parameters.exportCsvSeparator ?? EXPORT_CSV_SEPARATORS.comma,
  exportCsvDatetimeFormat: parameters.exportCsvDatetimeFormat ?? EXPORT_CSV_DATETIME_FORMATS.datetimeSeconds.value,
  widgetExportColumns: parameters.widgetExportColumns
    ? cloneDeep(parameters.widgetExportColumns)
    : defaultColumnsToColumns(DEFAULT_ALARMS_WIDGET_COLUMNS),
});

export const formParametersToContextWidgetParameters = (parameters = {}) => ({
  ...parameters,

  sort: formSortToWidgetSort(parameters.sort), // TODO: move to common
});
