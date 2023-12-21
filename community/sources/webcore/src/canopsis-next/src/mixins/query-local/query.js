import { isEqual, omit } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';

import { convertSortToRequest } from '@/helpers/entities/shared/query';

export const localQueryMixin = {
  data() {
    return {
      query: {
        page: 1,
        itemsPerPage: PAGINATION_LIMIT,
        search: '',
        sortBy: [],
        sortDesc: [],
      },
    };
  },

  watch: {
    query(query, oldQuery) {
      if (this.customQueryCondition(query, oldQuery)) {
        this.fetchList();
      }
    },
  },

  computed: {
    options: {
      get() {
        const { page = 1, itemsPerPage = PAGINATION_LIMIT, sortBy = [], sortDesc = [], ...restQuery } = this.query;

        return { page, itemsPerPage, sortBy, sortDesc, ...restQuery };
      },
      set(value) {
        this.query = {
          ...this.query,

          search: value.search || '',
          page: value.itemsPerPage <= this.query.itemsPerPage ? value.page : 1,
          itemsPerPage: value.itemsPerPage || PAGINATION_LIMIT,
          sortBy: value.sortBy || [],
          sortDesc: value.sortDesc || [],

          ...omit(value, [
            'search',
            'page',
            'itemsPerPage',
            'totalItems',
            'groupBy',
            'groupDesc',
            'multiSort',
            'mustSort',
          ]),
        };
      },
    },
  },

  methods: {
    updateQueryField(field, value) {
      this.query = {
        ...this.query,
        [field]: value,
      };
    },

    customQueryCondition(query, oldQuery) {
      return !isEqual(query, oldQuery);
    },

    getQuery({
      page,
      search,
      itemsPerPage,
      sortBy = [],
      sortDesc = [],
    } = this.query) {
      const query = {
        page,
        limit: itemsPerPage,

        ...convertSortToRequest(sortBy, sortDesc),
      };

      if (search) {
        query.search = search;
      }

      return query;
    },
  },

  fetchList() {},
};
