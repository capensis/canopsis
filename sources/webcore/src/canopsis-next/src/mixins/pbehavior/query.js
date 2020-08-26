import { PAGINATION_LIMIT } from '@/config';
import { SORT_ORDERS } from '@/constants';

import queryMixin from '@/mixins/query';

export default {
  mixins: [queryMixin],

  props: {
    queryId: {
      type: [Number, String],
      required: true,
    },
  },

  created() {
    this.updateQuery({
      id: this.queryId,
      query: {
        page: 1,
        limit: PAGINATION_LIMIT,
      },
    });
  },

  beforeDestroy() {
    this.removeQuery({ id: this.queryId });
  },

  computed: {
    query: {
      set(query) {
        this.updateQuery({ id: this.queryId, query: this.prepareQuery(query) });
      },
      get() {
        return this.getQueryById(this.queryId);
      },
    },

    pagination: {
      set(value) {
        this.query = value;
      },
      get() {
        return {
          page: this.query.page || 1,
          rowsPerPage: this.query.limit || PAGINATION_LIMIT,
          search: this.query.search || '',
          sortBy: this.query.sort_key || '',
          descending: this.query.sort_dir === SORT_ORDERS.desc,
        };
      },
    },
  },

  methods: {
    prepareQuery({
      page,
      search,
      rowsPerPage,
      sortKey,
      sortDir,
    }) {
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
