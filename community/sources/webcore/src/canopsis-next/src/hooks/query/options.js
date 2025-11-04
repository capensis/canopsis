import { omit } from 'lodash';
import { computed, unref } from 'vue';

import { PAGINATION_LIMIT } from '@/config';

import { isSeveralEqual } from '@/helpers/collection';

/**
 * Custom hook to manage query options for pagination and sorting.
 *
 * @param {Ref} query - The reactive reference to the query object.
 * @param {Function} updateQuery - The function to update the query object.
 * @returns {Object} - Object containing computed options for pagination and sorting.
 */
export const useQueryOptions = (query, updateQuery) => {
  const options = computed(() => {
    const { page = 1, itemsPerPage = PAGINATION_LIMIT, sortBy = [], sortDesc = [], ...restQuery } = unref(query);

    return { page, itemsPerPage, sortBy, sortDesc, ...restQuery };
  });

  /**
   * Update options based on the provided value.
   *
   * @param {Object} value - The new options values to update.
   */
  const updateOptions = (value) => {
    const unwrappedQuery = unref(query);
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

    const isEqualOptions = isSeveralEqual(
      options.value,
      value,
      ['search', 'page', 'itemsPerPage', 'sortBy', 'sortDesc', ...Object.keys(otherValues)],
    );

    if (isEqualOptions) {
      return;
    }

    const oldItemsPerPage = unwrappedQuery.itemsPerPage ?? PAGINATION_LIMIT;

    updateQuery({
      ...unwrappedQuery,

      search: value.search || '',
      page: value.itemsPerPage <= oldItemsPerPage ? value.page : 1,
      itemsPerPage: value.itemsPerPage || PAGINATION_LIMIT,
      sortBy: value.sortBy || [],
      sortDesc: value.sortDesc || [],

      ...otherValues,
    });
  };

  const optionsWithSet = computed({
    /**
     * Get function to extract pagination and sorting options from the query object.
     *
     * @returns {Object} - Object containing page, itemsPerPage, sortBy, sortDesc, and other query parameters.
     */
    get: () => options.value,

    /**
     * Set function to update the query object based on the new options value.
     *
     * @param {Object} value - New options object for pagination and sorting.
     */
    set: updateOptions,
  });

  return { options: optionsWithSet, updateOptions };
};
