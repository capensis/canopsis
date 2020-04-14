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
        total: this.alarms.length,
      };
    },
    displayedAlarms() {
      const { page, limit } = this.query;

      return this.alarms.slice((page - 1) * limit, page * limit);
    },
  },
  async mounted() {
    await this.fetchItemWithGroups(this.alarm);

    this.query = {
      ...convertWidgetToQuery(this.widget),
    };
  },
};
