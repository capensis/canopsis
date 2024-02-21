import { convertWidgetQueryToRequest } from '@/helpers/entities/shared/query';
import { getPageForNewItemsPerPage } from '@/helpers/pagination';

import { queryMixin } from '@/mixins/query';
import { entitiesUserPreferenceMixin } from '@/mixins/entities/user-preference';

import { widgetOptionsMixin } from './options';

export const queryWidgetMixin = {
  mixins: [
    queryMixin,
    widgetOptionsMixin,
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
  },
  methods: {
    getQuery() {
      return convertWidgetQueryToRequest(this.query);
    },

    updateItemsPerPage(itemsPerPage) {
      this.updateLockedQueryById({
        id: this.queryId,
        query: {
          itemsPerPage,
          page: getPageForNewItemsPerPage(itemsPerPage, this.query.itemsPerPage, this.query.page),
        },
      });
    },

    updatePage(page) {
      this.query = { ...this.query, page };
    },
  },
};
