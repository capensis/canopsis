import { omit } from 'lodash';
import { computed, unref } from 'vue';

import { PAGINATION_LIMIT } from '@/config';

import { isPickEqual } from '@/helpers/collection';

/**
 * Custom hook to manage query options for pagination and sorting.
 *
 * @param {Ref} query - The reactive reference to the query object.
 * @param {Function} updateQuery - The function to update the query object.
 * @returns {Object} - Object containing computed options for pagination and sorting.
 */
export const useQueryOptions = (query, updateQuery) => {
  const unwrappedQuery = unref(query);
  const options = computed({
    /**
     * Get function to extract pagination and sorting options from the query object.
     *
     * @returns {Object} - Object containing page, itemsPerPage, sortBy, sortDesc, and other query parameters.
     */
    get: () => {
      const { page = 1, itemsPerPage = PAGINATION_LIMIT, sortBy = [], sortDesc = [], ...restQuery } = query.value;

      return { page, itemsPerPage, sortBy, sortDesc, ...restQuery };
    },

    /**
     * Set function to update the query object based on the new options value.
     *
     * @param {Object} value - New options object for pagination and sorting.
     */
    set: (value) => {
      const otherValues = omit(value, [
        'search',
        'page',
        'itemsPerPage',
        'totalItems',
        'groupBy',
        'groupDesc',
        'multiSort',
        'mustSort',
      ]);

      const isEqualOptions = isPickEqual(
        options.value,
        value,
        ['search', 'page', 'itemsPerPage', 'sortBy', 'sortDesc', ...Object.keys(otherValues)],
      );

      if (isEqualOptions) {
        return;
      }

      updateQuery({
        ...unwrappedQuery,

        search: value.search || '',
        page: value.itemsPerPage <= unwrappedQuery.itemsPerPage ? value.page : 1,
        itemsPerPage: value.itemsPerPage || PAGINATION_LIMIT,
        sortBy: value.sortBy || [],
        sortDesc: value.sortDesc || [],

        ...otherValues,
      });
    },
  });

  return { options };
};
