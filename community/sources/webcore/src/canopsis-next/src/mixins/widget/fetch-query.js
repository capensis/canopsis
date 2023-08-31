import { isEqual, isEmpty } from 'lodash';

import { prepareWidgetQuery } from '@/helpers/entities/widget/query';

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
      if (!isEqual(value, oldValue) && !isEmpty(value)) {
        this.fetchList();
      }
    },

    tabQueryNonce(value, oldValue) {
      if (!this.editing && value > oldValue) {
        this.fetchList();
      }
    },

    widget: 'setQuery',
  },
  async mounted() {
    if (!this.editing) {
      await this.fetchUserPreference({ id: this.widget._id });

      this.setQuery();
    }
  },
  methods: {
    setQuery() {
      const { search = '' } = this.query;

      this.query = {
        ...prepareWidgetQuery(this.widget, this.userPreference),

        search,
      };
    },
  },
};
