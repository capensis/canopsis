import { isEmpty } from 'lodash';
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

    fetchItemWithGroups(alarm) {
      const params = {
        with_steps: true,
      };

      if (!isEmpty(alarm.causes)) {
        params.with_causes = true;
      }

      if (!isEmpty(alarm.consequences)) {
        params.with_consequences = true;
      }

      return this.fetchAlarmItemWithParams(alarm, params);
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
