import { isArray, isEqual, omit, pick } from 'lodash';

import { DATETIME_FORMATS, SORT_ORDERS, ADVANCED_SEARCH_FIELDS } from '@/constants';
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
    instructionsFilter,
    search,
  } = query;

  let result = {
    ...omit(query, [
      'tstart',
      'tstop',
      'sortBy',
      'sortDesc',
      'category',
      'filter',
      'lockedFilter',
      'instructionsFilter',
      'search',
      'itemsPerPage',
    ]),

    ...convertSortToRequest(sortBy, sortDesc),

    limit: itemsPerPage,
  };

  if (lockedFilter || filter) {
    result.filters = convertFiltersToQuery(filter, lockedFilter);
  }

  if (instructionsFilter) {
    result = {
      ...result,
      ...instructionsFilter,
    };
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

/**
 * Constructs a query object for list data retrieval.
 *
 * @param {Object} [options = {}] - The options for constructing the query.
 * @param {number} [options.page = 1] - The current page number.
 * @param {string} [options.search] - The search term to filter results.
 * @param {number} options.itemsPerPage - The number of items per page.
 * @param {string[]} [options.sortBy = []] - The fields to sort by.
 * @param {string[]} [options.sortDesc = []] - The sort order for each field in `sortBy`.
 * @param {string} [options.[ADVANCED_SEARCH_FIELDS.search]] - The search pattern to filter results.
 * @returns {Object} The constructed query object.
 */
export const getQueryForList = ({
  page = 1,
  search,
  itemsPerPage,
  sortBy = [],
  sortDesc = [],
  [ADVANCED_SEARCH_FIELDS.search]: searchPattern,
} = {}) => {
  const query = {
    page,
    limit: itemsPerPage,

    ...convertSortToRequest(sortBy, sortDesc),
  };

  if (search) {
    query.search = search;
  }

  if (searchPattern) {
    query[ADVANCED_SEARCH_FIELDS.search] = searchPattern;
  }

  return query;
};
