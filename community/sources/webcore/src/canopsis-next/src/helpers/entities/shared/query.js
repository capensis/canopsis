import { isArray, omit } from 'lodash';

import { DATETIME_FORMATS, SORT_ORDERS } from '@/constants';
import { PAGINATION_LIMIT } from '@/config';

import { convertStartDateIntervalToTimestamp, convertStopDateIntervalToTimestamp } from '@/helpers/date/date-intervals';

/**
 * This function converts widget.parameters.sort to query Object
 *
 * @param {Object} widget
 * @returns {{}}
 */
export function convertSortToQuery({ parameters }) {
  const { sort } = parameters;

  if (sort && sort.column && sort.order) {
    return { sortKey: sort.column, sortDir: sort.order };
  }

  return { sortKey: null, sortDir: null };
}

/**
 * Convert sortBy and sordDesc query parameters to request
 *
 * @param {string[]} sortBy
 * @param {string[]} sortDesc
 * @returns {{ sort_by: string, sort: string } | { multi_sort: string[] } | {}}
 */
export const convertSortToRequest = (sortBy = [], sortDesc = []) => {
  if (!sortBy?.length) {
    return {};
  }

  if (sortBy.length === 1) {
    return {
      sort_by: sortBy[0],
      sort: (sortDesc?.[0] ? SORT_ORDERS.desc : SORT_ORDERS.asc).toLowerCase(),
    };
  }

  return {
    multi_sort: sortBy
      .map((property, index) => (
        `${property},${(sortDesc?.[index] ? SORT_ORDERS.desc : SORT_ORDERS.asc).toLowerCase()}`
      )),
  };
};

/**
 * Convert filter to query filters
 *
 * @param {string | string[]} filter
 * @returns {string[]}
 */
export const convertFilterToQuery = filter => (isArray(filter) ? filter : [filter]).filter(Boolean);

/**
 * Convert locked filter and main filter to query filters
 *
 * @param {string | string[]} filter
 * @param {string | string[]} lockedFilter
 * @returns {string[]}
 */
export const convertFiltersToQuery = (filter, lockedFilter) => [
  ...convertFilterToQuery(filter),
  ...convertFilterToQuery(lockedFilter),
];

/**
 * Convert widget query to request parameters
 *
 * @param {Object} query
 * @returns {Object}
 */
export const convertWidgetQueryToRequest = (query) => {
  const {
    sortBy = [],
    sortDesc = [],
    limit = PAGINATION_LIMIT,
    tstart,
    tstop,
    category,
    filter,
    lockedFilter,
    search,
  } = query;

  const result = {
    ...omit(query, [
      'tstart',
      'tstop',
      'sortBy',
      'sortDesc',
      'category',
      'limit',
      'filter',
      'lockedFilter',
      'search',
    ]),
    ...convertSortToRequest(sortBy, sortDesc),
  };

  if (lockedFilter || filter) {
    result.filters = convertFiltersToQuery(filter, lockedFilter);
  }

  if (tstart) {
    result.tstart = convertStartDateIntervalToTimestamp(tstart, DATETIME_FORMATS.dateTimePicker);
  }

  if (tstop) {
    result.tstop = convertStopDateIntervalToTimestamp(tstop, DATETIME_FORMATS.dateTimePicker);
  }

  if (category) {
    result.category = category;
  }

  if (search) {
    result.search = search;
  }

  result.limit = limit;

  return result;
};
