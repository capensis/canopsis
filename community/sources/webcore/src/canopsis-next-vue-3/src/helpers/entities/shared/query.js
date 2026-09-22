import { isArray, isEqual, omit, pick } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';
import { DATETIME_FORMATS, SORT_ORDERS } from '@/constants';

import {
  convertStartDateIntervalToTimestamp,
  convertStopDateIntervalToTimestamp,
} from '@/helpers/date/date-intervals';

/**
 * @typedef {Object} DataTableSortOptions
 * @property {string[]} sortBy
 * @property {boolean[]} sortDesc
 */

/**
 * Convert widget.parameters.sort to query object.
 *
 * @param {Object} widget
 * @returns {{ sortBy: string[], sortDesc: boolean[] }}
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
 * Convert Vuetify data table sort options to query options with comparison.
 *
 * @param {Object} [newOptions={}]
 * @param {Object} [oldOptions={}]
 * @param {string[]} [optionsKeys=['page', 'itemsPerPage', 'sortBy', 'sortDesc']]
 * @returns {Object}
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
 * Convert sortBy and sortDesc query parameters to request.
 *
 * @param {string[]} sortBy
 * @param {boolean[]} sortDesc
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
    multi_sort: sortBy.map(
      (property, index) =>
        `${property},${(sortDesc?.[index] ? SORT_ORDERS.desc : SORT_ORDERS.asc).toLowerCase()}`,
    ),
  };
};

/**
 * Convert canonical query sort to Vuetify 4 data table sort options.
 *
 * @param {string[]} [sortBy=[]]
 * @param {boolean[]} [sortDesc=[]]
 * @returns {{ key: string, order: string }[]}
 */
export const convertQuerySortToDataTableSort = (sortBy = [], sortDesc = []) =>
  sortBy.map((key, index) => ({
    key,
    order: sortDesc?.[index] ? SORT_ORDERS.desc.toLowerCase() : SORT_ORDERS.asc.toLowerCase(),
  }));

/**
 * Convert Vuetify 4 data table sort options to canonical query sort.
 *
 * @param {{ sortBy?: string[], sortDesc?: boolean[] } | { key: string, order: string }[]} [options=[]]
 * @returns {{ sortBy: string[], sortDesc: boolean[] }}
 */
export const convertDataTableSortToQuery = (options = []) => {
  const sortBy = isArray(options) ? options : options.sortBy ?? [];

  return {
    sortBy: sortBy.map(sortItem => sortItem.key ?? sortItem),
    sortDesc: sortBy.map((sortItem, index) =>
      sortItem.order
        ? sortItem.order === SORT_ORDERS.desc.toLowerCase()
        : (options.sortDesc?.[index] ?? false),
    ),
  };
};

/**
 * Convert filter to query filters.
 *
 * @param {string | string[]} filter
 * @returns {string[]}
 */
export const convertFilterToQuery = filter => (isArray(filter) ? filter : [filter]).filter(Boolean);

/**
 * Convert locked filter and main filter to query filters.
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
 * Convert widget query to request parameters.
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
 * @param {Object} [options={}] - The options for constructing the query.
 * @param {number} [options.page=1] - The current page number.
 * @param {string} [options.search] - The search term to filter results.
 * @param {number} options.itemsPerPage - The number of items per page.
 * @param {string[]} [options.sortBy=[]] - The fields to sort by.
 * @param {boolean[]} [options.sortDesc=[]] - The sort order for each field in `sortBy`.
 * @returns {Object} The constructed query object.
 */
export const getQueryForList = ({
  page = 1,
  category,
  search,
  itemsPerPage,
  sortBy = [],
  sortDesc = [],
  search_pattern: searchPattern,
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
    query.search_pattern = searchPattern;
  }

  if (category) {
    query.category = category;
  }

  return query;
};
