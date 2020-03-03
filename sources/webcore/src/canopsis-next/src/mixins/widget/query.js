import { omit } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';
import { SORT_ORDERS, DATETIME_FORMATS } from '@/constants';
import queryMixin from '@/mixins/query';
import { dateParse } from '@/helpers/date-intervals';

/**
 * @mixin Add widget query logic
 */
export default {
  mixins: [queryMixin],
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
    queryId() {
      return this.defaultQueryId || this.widget._id;
    },
    query: {
      get() {
        return this.getQueryById(this.queryId);
      },
      set(query) {
        return this.updateQuery({ id: this.queryId, query });
      },
    },

    vDataTablePagination: {
      get() {
        const descending = this.query.sortDir !== null ? this.query.sortDir === SORT_ORDERS.desc : null;

        return { sortBy: this.query.sortKey, descending };
      },
      set(value) {
        const isNotEqualSortBy = value.sortBy !== this.vDataTablePagination.sortBy;
        const isNotEqualDescending = value.descending !== this.vDataTablePagination.descending;

        if (isNotEqualSortBy || isNotEqualDescending) {
          this.query = {
            ...this.query,
            sortKey: value.sortBy,
            sortDir: value.descending ? SORT_ORDERS.desc : SORT_ORDERS.asc,
          };
        }
      },
    },

    tabQueryNonce() {
      return this.getQueryNonceById(this.tabId);
    },
  },
  destroyed() {
    this.removeQuery({ id: this.queryId });
  },
  methods: {
    getQuery() {
      const query = omit(this.query, [
        'page',
        'sortKey',
        'sortDir',
        'tstart',
        'tstop',
      ]);

      const {
        page,
        tstart,
        tstop,
        limit = PAGINATION_LIMIT,
      } = this.query;

      if (tstart) {
        const convertedTstart = dateParse(tstart, 'start', DATETIME_FORMATS.dateTimePicker);

        query.tstart = convertedTstart.unix();
      }

      if (tstop) {
        const convertedTstop = dateParse(tstop, 'stop', DATETIME_FORMATS.dateTimePicker);

        query.tstop = convertedTstop.unix();
      }

      if (this.query.sortKey) {
        query.sort_key = this.query.sortKey;
        query.sort_dir = this.query.sortDir;
      }

      query.limit = limit;
      query.skip = ((page - 1) * limit) || 0;

      return query;
    },
    updateRecordsPerPage(limit) {
      this.updateLockedQuery({
        id: this.queryId,
        query: { limit },
      });
    },
  },
};
