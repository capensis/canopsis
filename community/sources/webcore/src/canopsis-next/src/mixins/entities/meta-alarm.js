import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('metaAlarm');

export const entitiesMetaAlarmMixin = {
  methods: {
    ...mapActions({
      removeAlarmsFromMetaAlarm: 'removeAlarms',
    }),
  },
};
