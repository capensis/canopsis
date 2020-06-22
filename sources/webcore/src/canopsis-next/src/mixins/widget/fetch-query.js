import { isEqual, isEmpty } from 'lodash';

import queryWidgetMixin from '@/mixins/widget/query';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';
import { prepareQuery } from '@/helpers/query';

/**
 * @mixin Add query logic with fetch
 */
export default {
  mixins: [queryWidgetMixin, entitiesUserPreferenceMixin],
  props: {
    isEditingMode: {
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
      if (value > oldValue) {
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
