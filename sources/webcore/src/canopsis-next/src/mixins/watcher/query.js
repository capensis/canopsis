import { pick } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';

export default {
  data() {
    return {
      query: {
        page: 1,
        rowsPerPage: PAGINATION_LIMIT,
      },
    };
  },

  computed: {
    pagination: {
      set(value) {
        this.query = {
          ...this.query,
          page: value.page,
          rowsPerPage: value.rowsPerPage || PAGINATION_LIMIT,
        };
      },
      get() {
        return pick(this.query, ['page', 'rowsPerPage']);
      },
    },
  },

  methods: {
    getQuery({ page, rowsPerPage } = this.query) {
      const query = {};

      query.limit = rowsPerPage;
      query.page = page;

      return query;
    },
  },
};
