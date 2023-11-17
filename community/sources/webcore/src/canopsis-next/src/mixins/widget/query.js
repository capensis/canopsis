import { PAGINATION_LIMIT } from '@/config';

import { convertDataTableOptionsToQuery, convertWidgetQueryToRequest } from '@/helpers/entities/shared/query';

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

    options: {
      get() {
        const { page = 1, itemsPerPage = PAGINATION_LIMIT, sortBy = [], sortDesc = [] } = this.query;

        return { page, itemsPerPage, sortBy, sortDesc };
      },

      set(newOptions) {
        this.query = {
          ...this.query,
          ...convertDataTableOptionsToQuery(newOptions, this.options),
        };
      },
    },
  },
  methods: {
    getQuery() {
      return convertWidgetQueryToRequest(this.query);
    },

    updateItemsPerPage(itemsPerPage) {
      this.updateLockedQuery({
        id: this.queryId,
        query: { itemsPerPage },
      });
    },

    updateQueryPage(page) {
      this.query = { ...this.query, page };
    },
  },
};
