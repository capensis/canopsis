import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('metaAlarm');

export const entitiesMetaAlarmMixin = {
  methods: {
    ...mapActions({
      fetchMetaAlarmsListWithoutStore: 'fetchListWithoutStore',
      createMetaAlarm: 'create',
      addAlarmsIntoMetaAlarm: 'addAlarms',
      removeAlarmsFromMetaAlarm: 'removeAlarms',
    }),
  },
};
