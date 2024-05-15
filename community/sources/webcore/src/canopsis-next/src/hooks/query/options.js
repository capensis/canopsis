import { omit } from 'lodash';
import { computed, unref } from 'vue';

import { PAGINATION_LIMIT } from '@/config';

import { isPickEqual } from '@/helpers/collection';

export const useQueryOptions = (query, updateQuery) => {
  const unwrappedQuery = unref(query);
  const options = computed({
    get: () => {
      const { page = 1, itemsPerPage = PAGINATION_LIMIT, sortBy = [], sortDesc = [], ...restQuery } = query.value;

      return { page, itemsPerPage, sortBy, sortDesc, ...restQuery };
    },

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
