import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('alarm');

/**
 * @mixin
 */
export const entitiesAlarmMixin = {
  computed: {
    ...mapGetters({
      getAlarmItem: 'getItem',
      getAlarmsList: 'getList',
      getAlarmsListByWidgetId: 'getListByWidgetId',
      getAlarmsMetaByWidgetId: 'getMetaByWidgetId',
      getAlarmsPendingByWidgetId: 'getPendingByWidgetId',
    }),

    alarms() {
      return this.getAlarmsListByWidgetId(this.widget._id);
    },

    alarmsMeta() {
      return this.getAlarmsMetaByWidgetId(this.widget._id);
    },

    alarmsPending() {
      return this.getAlarmsPendingByWidgetId(this.widget._id);
    },
  },
  methods: {
    ...mapActions({
      fetchAlarmItem: 'fetchItem',
      fetchAlarmsList: 'fetchList',
      createAlarmsListExport: 'createAlarmsListExport',
      fetchAlarmsListExport: 'fetchAlarmsListExport',
    }),
  },
};
