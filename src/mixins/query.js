import { createNamespacedHelpers } from 'vuex';
import omit from 'lodash/omit';
import isEqual from 'lodash/isEqual';

import Pagination from '@/components/tables/pagination.vue';
import dateIntervals from '@/helpers/date-intervals';
import { convertWidgetToQuery, convertUserPreferenceToQuery } from '@/helpers/query';

const { mapGetters: queryMapGetters, mapActions: queryMapActions } = createNamespacedHelpers('query');

/**
 * @mixin Add query logic
 */
export default {
  components: {
    Pagination,
  },
  computed: {
    ...queryMapGetters(['getQueryById', 'getQueryPendingById']),

    query: {
      get() {
        return this.getQueryById(this.widget.id);
      },
      set(query) {
        return this.updateQuery({ id: this.widget.id, query });
      },
    },

    queryPending() {
      return this.getQueryPendingById(this.widget.id);
    },
  },
  watch: {
    query(value, oldValue) {
      if (!isEqual(value, oldValue)) {
        this.fetchList();
      }
    },
    queryPending() {
      this.fetchList();
    },
    widget(value) {
      const widgetQuery = convertWidgetToQuery(value);

      this.updateQuery({
        id: this.widget.id,
        query: {
          ...this.query,
          ...widgetQuery,
        },
      });
    },
    userPreference(value) {
      const userPreferenceQuery = convertUserPreferenceToQuery(value);

      this.updateQuery({
        id: this.widget.id,
        query: {
          ...this.query,
          ...userPreferenceQuery,
        },
      });
    },
  },
  async mounted() {
    await this.startQueryPending({ id: this.widget.id });
    await this.fetchUserPreferenceByWidgetId({ widgetId: this.widget.id });
    await this.stopQueryPending({ id: this.widget.id });
  },
  methods: {
    ...queryMapActions({
      updateQuery: 'update',
      startQueryPending: 'startPending',
      stopQueryPending: 'stopPending',
    }),

    getQuery() {
      const query = omit(this.query, [
        'page',
        'interval',
        'descending',
        'sortBy',
        'rowsPerPage',
        'totalItems',
      ]);

      const { page, interval } = this.query;

      if (interval && interval !== 'custom') {
        try {
          const { tstart, tstop } = dateIntervals[interval]();

          query.tstart = tstart;
          query.tstop = tstop;
        } catch (err) {
          console.warn(err);
        }
      }

      if (this.query.sortBy) {
        query.sort_key = this.query.sortBy;
        query.sort_dir = this.query.descending ? 'DESC' : 'ASC';
      }

      query.limit = this.query.rowsPerPage;
      query.skip = ((page - 1) * this.query.rowsPerPage) || 0;

      return query;
    },
    fetchList() {
      if (!this.queryPending) {
        let method;

        if (this.fetchAlarmsList) {
          method = 'fetchAlarmsList';
        } else if (this.fetchContextEntitiesList) {
          method = 'fetchContextEntitiesList';
        }

        if (method) {
          this[method]({
            widgetId: this.widget.id,
            params: this.getQuery(),
          });
        }
      }
    },
  },
};
