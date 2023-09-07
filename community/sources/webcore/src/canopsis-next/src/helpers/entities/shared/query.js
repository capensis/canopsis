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
 * Convert multiSortBy query parameter to request
 *
 * @param {Object[]} multiSortBy
 * @returns {string[]}
 */
export const convertMultiSortToRequest = (multiSortBy = []) => multiSortBy
  .map(({ sortBy, descending }) => `${sortBy},${(descending ? SORT_ORDERS.desc : SORT_ORDERS.asc).toLowerCase()}`);

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
  const result = omit(query, [
    'tstart',
    'tstop',
    'sortKey',
    'sortDir',
    'category',
    'multiSortBy',
    'limit',
    'filter',
    'lockedFilter',
  ]);

  const {
    tstart,
    tstop,
    sortKey,
    sortDir,
    category,
    filter,
    lockedFilter,
    multiSortBy = [],
    limit = PAGINATION_LIMIT,
  } = query;

  if (lockedFilter || filter) {
    result.filters = convertFiltersToQuery(filter, lockedFilter);
  }

  if (tstart) {
    result.tstart = convertStartDateIntervalToTimestamp(tstart, DATETIME_FORMATS.dateTimePicker);
  }

  if (tstop) {
    result.tstop = convertStopDateIntervalToTimestamp(tstop, DATETIME_FORMATS.dateTimePicker);
  }

  if (sortKey) {
    result.sort_by = sortKey;
    result.sort = sortDir.toLowerCase();
  }

  if (category) {
    result.category = category;
  }

  if (multiSortBy.length) {
    result.multi_sort = convertMultiSortToRequest(multiSortBy);
  }

  result.limit = limit;

  return result;
};
