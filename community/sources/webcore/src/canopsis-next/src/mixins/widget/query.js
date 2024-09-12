import { pick, isEqual } from 'lodash';

import { SORT_ORDERS } from '@/constants';
import { PAGINATION_LIMIT } from '@/config';

import { convertWidgetQueryToRequest } from '@/helpers/entities/shared/query';
import { getPageForNewRecordsPerPage } from '@/helpers/pagination';

import { queryMixin } from '@/mixins/query';
import { entitiesUserPreferenceMixin } from '@/mixins/entities/user-preference';

export const queryWidgetMixin = {
  mixins: [
    queryMixin,
    entitiesUserPreferenceMixin,
  ],
  props: {
    tabId: {
      type: String,
      required: true,
    },
    defaultQueryId: {
      type: [Number, String],
      required: false,
    },
  },
  computed: {
    query: {
      get() {
        return this.getQueryById(this.queryId);
      },
      set(query) {
        return this.updateQuery({ id: this.queryId, query });
      },
    },

    queryId() {
      return this.defaultQueryId || this.widget._id;
    },

    tabQueryNonce() {
      return this.getQueryNonceById(this.tabId);
    },

    pagination: {
      get() {
        const {
          sortDir,
          limit = PAGINATION_LIMIT,
          sortKey: sortBy = null,
          multiSortBy = [],
          page = 1,
        } = this.query;
        const descending = sortDir === SORT_ORDERS.desc;

        return {
          page,
          sortBy,
          descending,
          multiSortBy,
          rowsPerPage: limit,
        };
      },

      set(value) {
        const paginationKeys = ['sortBy', 'descending', 'multiSortBy', 'page', 'rowsPerPage'];
        const newPagination = pick(value, paginationKeys);
        const oldPagination = pick(this.pagination, paginationKeys);

        if (isEqual(newPagination, oldPagination)) {
          return;
        }

        const {
          page = 1,
          sortBy = null,
          descending = false,
          multiSortBy = [],
          rowsPerPage = PAGINATION_LIMIT,
        } = newPagination;

        const newQuery = {
          multiSortBy,
          sortKey: sortBy,
          sortDir: descending ? SORT_ORDERS.desc : SORT_ORDERS.asc,
          page: rowsPerPage <= this.query.rowsPerPage ? page : 1,
          limit: rowsPerPage || PAGINATION_LIMIT,
        };

        this.query = {
          ...this.query,
          ...newQuery,
        };
      },
    },
  },
  methods: {
    getQuery() {
      return convertWidgetQueryToRequest(this.query);
    },

    updateRecordsPerPage(limit) {
      this.updateLockedQuery({
        id: this.queryId,
        query: {
          limit,
          page: getPageForNewRecordsPerPage(limit, this.query.limit, this.query.page),
        },
      });
    },

    updateQueryPage(page) {
      this.query = { ...this.query, page };
    },
  },
};
