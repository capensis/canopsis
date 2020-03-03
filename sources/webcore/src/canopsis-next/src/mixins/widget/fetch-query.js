import { isEqual, isEmpty } from 'lodash';

import queryWidgetMixin from '@/mixins/widget/query';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';
import { convertWidgetToQuery, convertUserPreferenceToQuery } from '@/helpers/query';

/**
 * @mixin Add query logic with fetch
 */
export default {
  mixins: [queryWidgetMixin, entitiesUserPreferenceMixin],
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
    await this.fetchUserPreferenceByWidgetId({ widgetId: this.widget._id });

    this.query = {
      ...this.query,
      ...convertWidgetToQuery(this.widget),
      ...convertUserPreferenceToQuery(this.userPreference),
    };
  },
  updateRecordsPerPage(limit) {
    this.updateLockedQuery({
      id: this.queryId,
      query: { limit },
    });
  },
};
