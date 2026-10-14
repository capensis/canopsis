import { isEmpty, isEqual } from 'lodash';

import { prepareWidgetQuery } from '@/helpers/entities/widget/query';

import { queryWidgetMixin } from './query';

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
      if (this.customQueryCondition(value, oldValue) && this.visible) {
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
    customQueryCondition(value, oldValue) {
      return !isEqual(value, oldValue) && !isEmpty(value);
    },

    setQuery() {
      const { search = '' } = this.query;

      this.query = {
        ...prepareWidgetQuery(this.widget, this.userPreference),

        search,
      };
    },
  },
};
