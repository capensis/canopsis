import { isArray, isEqual, omit, pick } from 'lodash';

import { DATETIME_FORMATS, SORT_ORDERS } from '@/constants';
import { PAGINATION_LIMIT } from '@/config';

import { convertStartDateIntervalToTimestamp, convertStopDateIntervalToTimestamp } from '@/helpers/date/date-intervals';

/**
 * @typedef {Object} DataTableSortOptions
 * @property {string[]} sortBy
 * @property {boolean[]} sortDesc
 */

/**
 * This function converts widget.parameters.sort to query Object
 *
 * @param {Object} widget
 * @returns {DataTableSortOptions | {}}
 */
export const convertSortToQuery = ({ parameters }) => {
  const { sort } = parameters;
  const query = { sortBy: [], sortDesc: [] };

  if (sort && sort.column && sort.order) {
    query.sortBy = [sort.column];
    query.sortDesc = [sort.order === SORT_ORDERS.desc];
  }

  return query;
};

/**
 * Convert vuetify data table sort options to sort options with comparison
 *
 * @param {DataTableSortOptions} [newOptions = {}]
 * @param {DataTableSortOptions} [oldOptions = {}]
 * @param {string[]} [optionsKeys = ['page', 'itemsPerPage', 'sortBy', 'sortDesc']]
 * @return {DataTableSortOptions}
 */
export const convertDataTableOptionsToQuery = (
  newOptions = {},
  oldOptions = {},
  optionsKeys = ['page', 'itemsPerPage', 'sortBy', 'sortDesc'],
) => {
  const newOptionsToCompare = pick(newOptions, optionsKeys);
  const oldOptionsToCompare = pick(oldOptions, optionsKeys);

  if (isEqual(newOptionsToCompare, oldOptionsToCompare)) {
    return oldOptions;
  }

  const { page = 1, itemsPerPage = PAGINATION_LIMIT, sortBy = [], sortDesc = [] } = newOptions;

  return { page, itemsPerPage, sortBy, sortDesc };
};

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
    itemsPerPage = PAGINATION_LIMIT,
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
      'filter',
      'lockedFilter',
      'search',
      'itemsPerPage',
    ]),

    ...convertSortToRequest(sortBy, sortDesc),

    limit: itemsPerPage,
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

  return result;
};
