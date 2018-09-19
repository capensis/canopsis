import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('alarm');

/**
 * @mixin
 */
export default {
  computed: {
    ...mapGetters({
      getAlarmItem: 'getItem',
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
      fetchAlarmsListWithPreviousParams: 'fetchListWithPreviousParams',
    }),
  },
};
