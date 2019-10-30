import { isEmpty, isEqual } from 'lodash';

import { SORT_ORDERS } from '@/constants';

export default {
  data() {
    return {
      query: {},
      pagination: {},
    };
  },
  watch: {
    query(value, oldValue) {
      if (!isEqual(value, oldValue) && !isEmpty(value)) {
        this.fetchList();
      }
    },
    pagination(value, oldValue) {
      if (!isEqual(value, oldValue)) {
        this.query = {
          ...this.query,

          page: value.page,
          limit: value.rowsPerPage,
          sortKey: value.sortBy,
          sortDir: value.descending ? SORT_ORDERS.desc : SORT_ORDERS.asc,
        };
      }
    },
  },
  methods: {
    getBaseQuery() {
      const {
        limit,
        page,
        sortKey,
        sortDir,
      } = this.query;

      return {
        limit,
        start: (page - 1) * limit,
        sort: [{
          property: sortKey,
          direction: sortDir,
        }],
      };
    },

    getQuery() {
      return this.getBaseQuery();
    },
  },
};
