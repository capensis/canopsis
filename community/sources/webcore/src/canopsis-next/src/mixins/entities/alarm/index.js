import { createNamespacedHelpers } from 'vuex';

import { EXPORT_STATUSES } from '@/constants';

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
      getAlarmsExportByWidgetId: 'getExportByWidgetId',
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

    alarmsExportPending() {
      const exportData = this.getAlarmsExportByWidgetId(this.widget._id);

      return exportData && exportData.status === EXPORT_STATUSES.running;
    },
  },
  methods: {
    ...mapActions({
      fetchAlarmItem: 'fetchItem',
      fetchAlarmsList: 'fetchList',
      createAlarmsListExport: 'createAlarmsListExport',
      fetchAlarmsListExport: 'fetchAlarmsListExport',
      fetchAlarmsListCsvFile: 'fetchAlarmsListCsvFile',
    }),
  },
};
