import { SORT_ORDERS } from '@/constants';
import queryMixin from '@/mixins/query';
import { convertQueryToWidgetParams } from '@/helpers/query';

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
      return convertQueryToWidgetParams(this.query);
    },
    updateRecordsPerPage(limit) {
      this.updateLockedQuery({
        id: this.queryId,
        query: { limit },
      });
    },
  },
};
