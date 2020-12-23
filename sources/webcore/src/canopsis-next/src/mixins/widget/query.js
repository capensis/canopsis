import { omit } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';
import { DATETIME_FORMATS, SORT_ORDERS } from '@/constants';

import { dateParse } from '@/helpers/date/date-intervals';

import queryMixin from '@/mixins/query';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';

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
        const { sortDir, sortKey: sortBy } = this.query;
        const descending = sortDir === SORT_ORDERS.desc;

        return { sortBy, descending };
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
  },
  methods: {
    getQuery() {
      const query = omit(this.query, [
        'sortKey',
        'sortDir',
        'tstart',
        'tstop',
      ]);

      const {
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
        query.sort_dir = this.query.sortDir.toLowerCase();
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
