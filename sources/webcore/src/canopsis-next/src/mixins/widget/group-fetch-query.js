import queryWidgetMixin from '@/mixins/widget/query';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';
import { convertWidgetToQuery, convertUserPreferenceToQuery } from '@/helpers/query';

/**
 * @mixin Add query logic with group fetch
 */
export default {
  mixins: [queryWidgetMixin, entitiesUserPreferenceMixin],
  async mounted() {
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
