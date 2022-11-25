import { get } from 'lodash';

import queryWidgetMixin from '@/mixins/widget/query';
import widgetExpandPanelAlarm from '@/mixins/widget/expand-panel/alarm/expand-panel';

import { convertWidgetToQuery } from '@/helpers/query';

/**
 * @mixin Add query logic with group fetch
 */
export default {
  mixins: [queryWidgetMixin, widgetExpandPanelAlarm],
  props: {
    alarm: {
      type: Object,
      required: true,
    },
  },
  computed: {
    alarms() {
      return get(this.alarm, 'consequences.data') || get(this.alarm, 'causes.data', []);
    },

    alarmsMeta() {
      return {
        total_count: this.alarms.length,
      };
    },
  },
  mounted() {
    this.query = convertWidgetToQuery(this.widget);
  },
};
