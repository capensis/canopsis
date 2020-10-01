import { isEqual, pick } from 'lodash';

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

  watch: {
    query(query, oldQuery) {
      if (!isEqual(query, oldQuery)) {
        this.fetchList();
      }
    },
  },

  computed: {
    pagination: {
      set(value) {
        this.query = {
          ...this.query,
          search: value.search || '',
          page: value.page,
          rowsPerPage: value.rowsPerPage || PAGINATION_LIMIT,
          sortKey: value.sortBy,
          sortDir: value.descending ? SORT_ORDERS.desc : SORT_ORDERS.asc,
        };
      },
      get() {
        return {
          ...pick(this.query, ['page', 'rowsPerPage', 'search']),
          sortBy: this.query.sortKey,
          descending: this.query.sortDir === SORT_ORDERS.desc,
        };
      },
    },
  },

  methods: {
    getQuery({
      page,
      search,
      rowsPerPage,
      sortBy,
      sortDir,
    } = this.query) {
      const query = {};

      query.limit = rowsPerPage;
      query.page = page;

      if (sortBy) {
        query.sort_key = sortBy;
        query.sort_dir = sortDir;
      }

      if (search) {
        query.search = search;
      }

      return query;
    },
  },

  fetchList() {},
};
