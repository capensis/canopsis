import queryWidgetMixin from '@/mixins/widget/query';
import { convertWidgetToQuery } from '@/helpers/query';

/**
 * @mixin Add query logic with group fetch
 */
export default {
  mixins: [queryWidgetMixin],
  props: {
    alarms: {
      type: Array,
      required: true,
    },
  },
  computed: {
    alarmsMeta() {
      return {
        total: this.alarms.length,
      };
    },
    displayedAlarms() {
      const { page, limit } = this.query;

      return this.alarms.slice((page - 1) * limit, page * limit);
    },
  },
  mounted() {
    this.query = {
      ...convertWidgetToQuery(this.widget),
    };
  },
};
