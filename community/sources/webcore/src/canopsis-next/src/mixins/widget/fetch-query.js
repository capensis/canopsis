import { isEqual, isEmpty } from 'lodash';

import queryWidgetMixin from '@/mixins/widget/query';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';
import { prepareQuery } from '@/helpers/query';

/**
 * @mixin Add query logic with fetch
 */
export const widgetFetchQueryMixin = {
  mixins: [queryWidgetMixin, entitiesUserPreferenceMixin],
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
        this.fetchList();
      }
    },
  },
  async mounted() {
    if (!this.isEditingMode) {
      await this.fetchUserPreferenceByWidgetId({ widgetId: this.widget._id });

      this.query = prepareQuery(this.widget, this.userPreference);
    }
  },
};
