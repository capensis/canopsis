import { omit } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';
import { DATETIME_FORMATS } from '@/constants';

import { dateParse } from '@/helpers/date-intervals';

import queryMixin from '@/mixins/query';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';
import vuetifyPaginationMixinCreator from '@/mixins/vuetify/pagination-creator';

/**
 * @mixin Add query logic
 */
export default {
  mixins: [
    queryMixin,
    entitiesUserPreferenceMixin,
    vuetifyPaginationMixinCreator({
      field: 'vDataTablePagination',
      mutating: true,
    }),
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
  },
  destroyed() {
    this.removeQuery({
      id: this.queryId,
    });
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
