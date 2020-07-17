import { omit } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';

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

    getQuery() {
      const query = {};

      const { page, rowsPerPage = PAGINATION_LIMIT, search } = this.query;

      query.limit = rowsPerPage;
      query.page = page;

      if (search) {
        query.search = search;
      }

      return query;
    },
  },
};
