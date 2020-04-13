import { omit, isEqual, isEmpty } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';
import { SORT_ORDERS, DATETIME_FORMATS } from '@/constants';
import queryMixin from '@/mixins/query';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';
import { dateParse } from '@/helpers/date-intervals';
import { prepareQuery } from '@/helpers/query';

/**
 * @mixin Add query logic
 */
export default {
  mixins: [queryMixin, entitiesUserPreferenceMixin],
  props: {
    tabId: {
      type: String,
      required: true,
    },
  },
  computed: {
    query: {
      get() {
        return this.getQueryById(this.widget._id);
      },
      set(query) {
        return this.updateQuery({ id: this.widget._id, query });
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
  watch: {
    query(value, oldValue) {
      if (!isEqual(value, oldValue) && !isEmpty(value)) {
        this.fetchList();
      }
    },
    tabQueryNonce(value, oldValue) {
      if (value > oldValue) {
        this.fetchList();
      }
    },
  },
  async mounted() {
    await this.fetchUserPreferenceByWidgetId({ widgetId: this.widget._id });

    this.query = prepareQuery(this.widget, this.userPreference);
  },
  destroyed() {
    this.removeQuery({
      id: this.widget._id,
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
  },
};
