import { isEmpty } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';

import { convertWidgetChartsToPerfDataQuery } from '@/helpers/entities/metric/query';
import { convertSortToQuery } from '@/helpers/entities/shared/query';

/**
 * This function converts widget with type 'Context' to query Object
 *
 * @param {Object} widget
 * @returns {{}}
 */
export function convertContextWidgetToQuery(widget) {
  const {
    itemsPerPage,
    selectedTypes,
    widgetColumns,
    mainFilter,
    charts = [],
  } = widget.parameters;

  const query = {
    page: 1,
    limit: itemsPerPage || PAGINATION_LIMIT,
    lockedFilter: mainFilter,
    perf_data: convertWidgetChartsToPerfDataQuery(charts),
  };

  if (widgetColumns) {
    query.active_columns = widgetColumns.map(v => v.value);
  }

  if (!isEmpty(selectedTypes)) {
    query.type = selectedTypes;
  }

  return { ...query, ...convertSortToQuery(widget) };
}

/**
 * This function converts userPreference with widgetXtype 'Context' to query Object
 *
 * @param {Object} userPreference
 * @returns {{ category: string }}
 */
export function convertContextUserPreferenceToQuery({ content }) {
  const { category, noEvents, mainFilter } = content;

  return {
    category,
    filter: mainFilter,
    no_events: noEvents,
  };
}
