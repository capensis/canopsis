import { isEqual, isEmpty } from 'lodash';

import { prepareQuery } from '@/helpers/query';

import { queryWidgetMixin } from '@/mixins/widget/query';

/**
 * @mixin Add query logic with fetch
 */
export const widgetFetchQueryMixin = {
  mixins: [queryWidgetMixin],
  props: {
    editing: {
      type: Boolean,
      default: false,
    },
  },
  watch: {
    query(value, oldValue) {
      if (!this.editing && !isEqual(value, oldValue) && !isEmpty(value)) {
        this.fetchList();
      }
    },
    tabQueryNonce(value, oldValue) {
      if (!this.editing && value > oldValue) {
        this.fetchList({ isQueryNonceUpdate: true });
      }
    },
    widget() {
      this.setQuery();
    },
  },
  async mounted() {
    if (!this.editing) {
      await this.fetchUserPreference({ id: this.widget._id });

      this.setQuery();
    }
  },
  methods: {
    setQuery() {
      this.query = prepareQuery(this.widget, this.userPreference);
    },
  },
};
