import { omit, pick, isEqual } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';
import { DATETIME_FORMATS, DATETIME_INTERVAL_TYPES, SORT_ORDERS } from '@/constants';

import { dateParse } from '@/helpers/date/date-intervals';

import queryMixin from '@/mixins/query';
import { entitiesUserPreferenceMixin } from '@/mixins/entities/user-preference';

/**
 * @mixin Add query logic
 */
export default {
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

    vDataTablePagination: {
      get() {
        const { sortDir, sortKey: sortBy = null, multiSortBy = [] } = this.query;
        const descending = sortDir === SORT_ORDERS.desc;

        return { sortBy, descending, multiSortBy };
      },

      set(value) {
        const paginationKeys = ['sortBy', 'descending', 'multiSortBy'];
        const newPagination = pick(value, paginationKeys);
        const oldPagination = pick(this.vDataTablePagination, paginationKeys);

        if (isEqual(newPagination, oldPagination)) {
          return;
        }

        const {
          sortBy = null,
          descending = false,
          multiSortBy = [],
        } = newPagination;

        const newQuery = {
          sortKey: sortBy,
          sortDir: descending ? SORT_ORDERS.desc : SORT_ORDERS.asc,
          multiSortBy,
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
      const query = omit(this.query, [
        'tstart',
        'tstop',
        'sortKey',
        'sortDir',
        'category',
        'multiSortBy',
        'limit',
      ]);

      const {
        tstart,
        tstop,
        sortKey,
        sortDir,
        category,
        multiSortBy = [],
        limit = PAGINATION_LIMIT,
      } = this.query;

      if (tstart) {
        const convertedTstart = dateParse(tstart, DATETIME_INTERVAL_TYPES.start, DATETIME_FORMATS.dateTimePicker);

        query.tstart = convertedTstart.unix();
      }

      if (tstop) {
        const convertedTstop = dateParse(tstop, DATETIME_INTERVAL_TYPES.stop, DATETIME_FORMATS.dateTimePicker);

        query.tstop = convertedTstop.unix();
      }

      if (sortKey) {
        query.sort_key = sortKey;
        query.sort_dir = sortDir.toLowerCase();
      }

      if (category) {
        query.category = category;
      }

      if (multiSortBy.length) {
        query.multi_sort = multiSortBy.map(({ sortBy, descending }) =>
          `${sortBy},${(descending ? SORT_ORDERS.desc : SORT_ORDERS.asc).toLowerCase()}`);
      }

      query.limit = limit;

      return query;
    },

    updateRecordsPerPage(limit) {
      this.updateLockedQuery({
        id: this.queryId,
        query: { limit },
      });
    },

    updateQueryPage(page) {
      this.query = { ...this.query, page };
    },
  },
};
