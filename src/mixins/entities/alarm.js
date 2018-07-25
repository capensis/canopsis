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
      return this.getAlarmsListByWidgetId(this.widget.id);
    },
    alarmsMeta() {
      return this.getAlarmsMetaByWidgetId(this.widget.id);
    },
    alarmsPending() {
      return this.getAlarmsPendingByWidgetId(this.widget.id);
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
