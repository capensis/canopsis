import { isArray, omit } from 'lodash';

import { WIDGET_TYPES } from '@/constants';

import { featuresService } from '@/services/features';

import { mapIds } from '@/helpers/array';
import {
  convertAlarmUserPreferenceToQuery,
  convertAlarmWidgetToQuery,
} from '@/helpers/entities/alarm/query';
import { convertCounterWidgetToQuery } from '@/helpers/entities/counter/query';
import {
  convertContextUserPreferenceToQuery,
  convertContextWidgetToQuery,
} from '@/helpers/entities/entity/query';
import {
  convertExternalDataTableUserPreferenceToQuery,
  convertExternalDataTableWidgetToQuery,
} from '@/helpers/entities/external-data-table/query';
import {
  convertChartUserPreferenceToQuery,
  convertChartWidgetToQuery,
  convertNumbersWidgetToQuery,
} from '@/helpers/entities/metric/query';

/**
 * This function converts userPreference to query Object
 *
 * @param {Object} userPreference
 * @param {WidgetType} widgetType
 * @returns {Object}
 */
export function convertUserPreferenceToQuery(userPreference, widgetType) {
  const convertersMap = {
    [WIDGET_TYPES.alarmList]: convertAlarmUserPreferenceToQuery,
    [WIDGET_TYPES.context]: convertContextUserPreferenceToQuery,
    [WIDGET_TYPES.barChart]: convertChartUserPreferenceToQuery,
    [WIDGET_TYPES.lineChart]: convertChartUserPreferenceToQuery,
    [WIDGET_TYPES.numbers]: convertChartUserPreferenceToQuery,
    [WIDGET_TYPES.externalDataTable]: convertExternalDataTableUserPreferenceToQuery,

    ...featuresService.get('helpers.query.convertUserPreferenceToQuery.convertersMap'),
  };

  const converter = convertersMap[widgetType];

  return converter ? converter(userPreference) : {};
}

/**
 * This function converts widget to query Object
 *
 * @param {Object} widget
 * @returns {{}}
 */
export function convertWidgetToQuery(widget) {
  const convertersMap = {
    [WIDGET_TYPES.alarmList]: convertAlarmWidgetToQuery,
    [WIDGET_TYPES.context]: convertContextWidgetToQuery,
    [WIDGET_TYPES.counter]: convertCounterWidgetToQuery,
    [WIDGET_TYPES.barChart]: convertChartWidgetToQuery,
    [WIDGET_TYPES.lineChart]: convertChartWidgetToQuery,
    [WIDGET_TYPES.numbers]: convertNumbersWidgetToQuery,
    [WIDGET_TYPES.externalDataTable]: convertExternalDataTableWidgetToQuery,

    ...featuresService.get('helpers.query.convertWidgetToQuery.convertersMap'),
  };

  const converter = convertersMap[widget.type];

  return converter ? converter(widget) : {};
}

/**
 * Prepare query by widget and userPreference objects
 *
 * @param {Object} widget
 * @param {Object} userPreference
 * @returns {Object}
 */
export function prepareWidgetQuery(widget, userPreference) {
  const { filters: widgetFilters = [] } = widget;
  const { filters: userPreferenceFilters = [] } = userPreference;
  const widgetQuery = convertWidgetToQuery(widget);
  const userPreferenceQuery = convertUserPreferenceToQuery(userPreference, widget.type);

  let query = {
    ...widgetQuery,
    ...userPreferenceQuery,
  };

  if (query.filter) {
    const allFiltersIds = mapIds([...widgetFilters, ...userPreferenceFilters]);

    if (isArray(query.filter)) {
      query.filter = query.filter.filter(id => allFiltersIds.includes(id));

      if (!query.filter.length) {
        query = omit(query, ['filter']);
      }
    } else if (!allFiltersIds.includes(query.filter)) {
      query = omit(query, ['filter']);
    }
  }

  return query;
}

/**
 * Transforms widget parameters into a format suitable for export query columns.
 *
 * @param {Object} [parameters={}] - The parameters object containing widget column information.
 * @param {WidgetColumn[]} [parameters.widgetExportColumns] - An array of export column objects, each containing
 *                                                            `value`, `text`, and `template`.
 * @param {WidgetColumn[]} [parameters.widgetColumns] - An array of widget column objects, each containing `value`,
 *                                                      `text`, and `template`.
 * @returns {Object[]} An array of objects, each representing a column with `name`, `label`, and optionally `template`.
 */
export const widgetToExportQueryColumns = ({ parameters } = {}) => {
  const { widgetExportColumns, widgetColumns } = parameters;

  const hasExportColumns = !!widgetExportColumns?.length;
  const columns = hasExportColumns ? widgetExportColumns : widgetColumns;

  return columns.map(({ value, text, template }) => ({
    name: value,
    label: text,
    template: hasExportColumns ? template : undefined,
  }));
};

/**
 * Converts a list of selected items into a query string or returns a search string if no items are selected.
 *
 * @param {Array} [selected=[]] - An array of selected items, each containing an `_id` property.
 * @param {string} [search=''] - A search string to return if no items are selected.
 * @returns {string} A query string constructed from selected items or the provided search string.
 */
export const selectedToQuery = (selected = [], search = '') =>
  selected?.length ? selected.map(alarm => `_id="${alarm._id}"`).join(' OR ') : search;
