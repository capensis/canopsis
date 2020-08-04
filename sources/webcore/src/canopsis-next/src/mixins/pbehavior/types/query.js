import { omit } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';
import { SORT_ORDERS } from '@/constants';

export default {
  data() {
    return {
      query: {
        page: 1,
        rowsPerPage: PAGINATION_LIMIT,
        search: '',
      },
    };
  },

  methods: {
    handleSearch(search) {
      this.query = {
        ...this.query,

        search,
      };
    },

    handleSearchClear() {
      this.query = omit(this.query, ['search']);
    },

    getQuery({
      page, search, rowsPerPage = PAGINATION_LIMIT, sortBy, descending,
    } = this.query) {
      const query = {};

      query.limit = rowsPerPage;
      query.page = page;

      if (sortBy) {
        query.sort_key = sortBy;
        query.sort_dir = descending ? SORT_ORDERS.desc : SORT_ORDERS.asc;
      }

      if (search) {
        query.search = search;
      }

      return query;
    },
  },
};
