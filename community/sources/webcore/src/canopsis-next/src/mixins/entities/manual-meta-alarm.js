import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('metaAlarm');

export const entitiesManualMetaAlarmMixin = {
  methods: {
    ...mapActions({
    }),
  },
};
