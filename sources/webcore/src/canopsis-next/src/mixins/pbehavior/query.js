import { pick } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';
import { SORT_ORDERS } from '@/constants';

export default {
  data() {
    return {
      query: {
        page: 1,
        rowsPerPage: PAGINATION_LIMIT,
        search: '',
        sortKey: '',
        sortDir: SORT_ORDERS.asc,
      },
    };
  },

  computed: {
    pagination: {
      set(value) {
        this.query = {
          ...this.query,
          page: value.page,
          search: value.search || '',
          rowsPerPage: value.rowsPerPage || PAGINATION_LIMIT,
          sortKey: value.sortBy,
          sortDir: value.descending ? SORT_ORDERS.desc : SORT_ORDERS.asc,
        };
      },
      get() {
        return {
          ...pick(this.query, ['page', 'rowsPerPage']),
          sortBy: this.query.sortKey,
          descending: this.query.sortDir === SORT_ORDERS.desc,
        };
      },
    },
  },

  methods: {
    getQuery({
      page, search, rowsPerPage, sortKey, sortDir,
    } = this.query) {
      const query = {};

      query.limit = rowsPerPage;
      query.page = page;

      if (sortKey) {
        query.sort_key = sortKey;
        query.sort_dir = sortDir;
      }

      if (search) {
        query.search = search;
      }

      return query;
    },
  },
};
