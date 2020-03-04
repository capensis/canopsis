import queryWidgetMixin from '@/mixins/widget/query';
import { convertWidgetToQuery } from '@/helpers/query';

/**
 * @mixin Add query logic with group fetch
 */
export default {
  mixins: [queryWidgetMixin],
  mounted() {
    this.query = {
      ...this.query,
      ...convertWidgetToQuery(this.widget),
    };
  },
  updateRecordsPerPage(limit) {
    this.updateLockedQuery({
      id: this.queryId,
      query: { limit },
    });
  },
};
