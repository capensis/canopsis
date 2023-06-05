import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('manualMetaAlarm');

export const entitiesManualMetaAlarmMixin = {
  methods: {
    ...mapActions({
      fetchManualMetaAlarmsListWithoutStore: 'fetchListWithoutStore',
      createManualMetaAlarm: 'create',
      addAlarmsIntoManualMetaAlarm: 'addAlarms',
      removeAlarmsFromManualMetaAlarm: 'removeAlarms',
    }),
  },
};
