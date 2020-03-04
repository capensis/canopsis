import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('alarm');

export default {
  methods: {
    ...mapActions({
      fetchAlarmItem: 'fetchItem',
    }),

    fetchItemWithSteps(alarm) {
      return this.fetchAlarmItemWithParams(alarm, {
        with_steps: true,
      });
    },

    fetchAlarmItemWithParams(alarm, params) {
      const defaultParams = {
        sort_key: 't',
        sort_dir: 'DESC',
        limit: 1,
      };

      if (alarm.v.resolved) {
        defaultParams.resolved = true;
      }

      return this.fetchAlarmItem({
        id: alarm._id,
        params: { ...defaultParams, ...params },
      });
    },
  },
};
