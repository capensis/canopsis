import queryWidgetMixin from '@/mixins/widget/query';
import { convertWidgetToQuery } from '@/helpers/query';
import { isEmpty, isEqual } from 'lodash';

/**
 * @mixin Add query logic with group fetch
 */
export default {
  mixins: [queryWidgetMixin],
  watch: {
    query(value, oldValue) {
      if (!isEqual(value, oldValue) && !isEmpty(value)) {
        this.fetchGroupAlarmListData();
      }
    },
    tabQueryNonce(value, oldValue) {
      if (value > oldValue) {
        this.fetchGroupAlarmListData();
      }
    },
  },
  mounted() {
    this.query = {
      ...convertWidgetToQuery(this.widget),
    };
  },
};
