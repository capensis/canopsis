import { entitiesAlarmDetailsMixin } from '@/mixins/entities/alarm/details';

export const widgetExpandPanelAlarmDetails = {
  mixins: [entitiesAlarmDetailsMixin],
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
    pending() {
      return this.getAlarmDetailsPending(this.widget._id, this.alarm._id);
    },

    alarmDetails() {
      return this.getAlarmDetailsItem(this.widget._id, this.alarm._id)?.data ?? {};
    },

    steps() {
      return this.alarmDetails?.steps ?? {};
    },

    filteredPerfData() {
      return this.alarmDetails?.filtered_perf_data ?? [];
    },

    children() {
      return this.alarmDetails?.children ?? {};
    },

    query: {
      get() {
        return this.getAlarmDetailsQuery(this.widget._id, this.alarm._id);
      },
      set(query) {
        return this.updateAlarmDetailsQuery({ widgetId: this.widget._id, id: this.alarm._id, query });
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
  methods: {
    updateStepsQuery(query) {
      this.stepsQuery = query;
    },
  },
};
