import { isEqual, isEmpty } from 'lodash';

import { prepareQuery } from '@/helpers/query';

import { queryWidgetMixin } from '@/mixins/widget/query';

/**
 * @mixin Add query logic with fetch
 */
export const widgetFetchQueryMixin = {
  mixins: [queryWidgetMixin],
  props: {
    isEditingMode: {
      type: Boolean,
      default: false,
    },
  },
  watch: {
    query(value, oldValue) {
      if (!this.isEditingMode && !isEqual(value, oldValue) && !isEmpty(value)) {
        this.fetchList();
      }
    },
    tabQueryNonce(value, oldValue) {
      if (!this.isEditingMode && value > oldValue) {
        this.fetchList({ isQueryNonceUpdate: true });
      }
    },
  },
  async mounted() {
    if (!this.isEditingMode) {
      await this.fetchUserPreference({ id: this.widget._id });

      this.query = prepareQuery(this.widget, this.userPreference);
    }
  },
};
