import { get, orderBy } from 'lodash';

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

    displayedAlarms() {
      const {
        page,
        limit,
        sortKey,
        sortDir = '',
      } = this.query;

      let { alarms } = this;

      if (sortKey) {
        alarms = orderBy(alarms, [sortKey], [sortDir.toLowerCase()]);
      }

      return alarms.slice((page - 1) * limit, page * limit);
    },
  },
  mounted() {
    this.query = convertWidgetToQuery(this.widget);
  },
};
