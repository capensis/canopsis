import { isEqual, isEmpty } from 'lodash';

import { prepareWidgetQuery } from '@/helpers/entities/widget/query';

import { queryWidgetMixin } from '@/mixins/widget/query';

/**
 * @mixin Add query logic with fetch
 */
export const widgetFetchQueryMixin = {
  mixins: [queryWidgetMixin],
  props: {
    visible: {
      type: Boolean,
      default: true,
    },
  },
  watch: {
    query(value, oldValue) {
      if (!isEqual(value, oldValue) && !isEmpty(value) && this.visible) {
        this.fetchList();
      }
    },

    tabQueryNonce(value, oldValue) {
      if (value > oldValue && this.visible) {
        this.fetchList();
      }
    },

    widget: 'setQuery',
  },
  async mounted() {
    await this.fetchUserPreference({ id: this.widget._id });

    this.setQuery();
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
