import { generateAlarmDetailsQueryId } from '@/helpers/query';

import { queryMixin } from '@/mixins/query';
import { entitiesAlarmDetailsMixin } from '@/mixins/entities/alarm/details';

export const widgetExpandPanelAlarmDetails = {
  mixins: [queryMixin, entitiesAlarmDetailsMixin],
  props: {
    alarm: {
      type: Object,
      required: true,
    },
    widget: {
      type: Object,
      required: true,
    },
  },
  computed: {
    queryId() {
      return generateAlarmDetailsQueryId(this.alarm, this.widget);
    },

    pending() {
      return this.getAlarmDetailsPending(this.queryId);
    },

    alarmDetails() {
      return this.getAlarmDetailsItem(this.alarm._id)?.data ?? {};
    },

    steps() {
      return this.alarmDetails?.steps ?? {};
    },

    children() {
      return this.alarmDetails?.children ?? {};
    },

    query: {
      get() {
        return this.getQueryById(this.queryId);
      },
      set(query) {
        return this.updateQuery({ id: this.queryId, query });
      },
    },

    stepsQuery: {
      get() {
        return this.query?.steps ?? {};
      },
      set(stepsQuery) {
        this.query = {
          ...this.query,

          steps: stepsQuery,
        };
      },
    },

    childrenQuery: {
      get() {
        return this.query?.children ?? {};
      },
      set(childrenQuery) {
        this.query = {
          ...this.query,

          children: childrenQuery,
        };
      },
    },
  },
};
