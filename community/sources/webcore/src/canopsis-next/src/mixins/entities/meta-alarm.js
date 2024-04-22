import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('metaAlarm');

export const entitiesMetaAlarmMixin = {
  methods: {
    ...mapActions({
      fetchManualMetaAlarmsListWithoutStore: 'fetchListWithoutStore',
      createMetaAlarm: 'create',
      addAlarmsIntoMetaAlarm: 'addAlarms',
      removeAlarmsFromMetaAlarm: 'removeAlarms',
    }),
  },
};
