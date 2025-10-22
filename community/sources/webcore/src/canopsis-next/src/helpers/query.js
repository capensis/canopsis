import { isNil } from 'lodash';

import { convertSortToRequest } from '@/helpers/entities/shared/query';
import { convertMetricIntervalToTimestamp } from '@/helpers/date/date-intervals';

/**
 * Converts query parameters into a request object.
 *
 * @param {Object} params - The query parameters.
 * @param {number} params.page - The current page number.
 * @param {string} [params.search] - The search term to filter results.
 * @param {number} [params.type] - The type to filter results.
 * @param {Object} [params.interval] - The interval term to filter results.
 * @param {number} params.itemsPerPage - The number of items per page.
 * @param {number} [params.type] - The type of the items to filter by.
 * @param {string[]} [params.sortBy=[]] - An array of fields to sort by.
 * @param {boolean[]} [params.sortDesc=[]] - An array indicating the sort direction for each field in `sortBy`.
 * @returns {Object} The request object containing pagination, sorting, and search parameters.
 */
export const convertQueryToRequest = ({
  page,
  search,
  type,
  interval,
  itemsPerPage,
  sortBy = [],
  sortDesc = [],
}) => {
  const query = {
    page,
    limit: itemsPerPage,

    ...convertSortToRequest(sortBy, sortDesc),
  };

  if (search) {
    query.search = search;
  }

  if (!isNil(type)) {
    query.type = type;
  }

  if (interval?.from && interval?.to) {
    const { from, to } = convertMetricIntervalToTimestamp({ interval });

    query.from = from;
    query.to = to;
  }

  if (!isNil(type)) {
    query.type = type;
  }

  return query;
};
