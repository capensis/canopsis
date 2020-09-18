import { createNamespacedHelpers } from 'vuex';
import { SORT_ORDERS } from '@/constants';

const { mapActions } = createNamespacedHelpers('alarm');

export default {
  props: {
    hideGroups: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    ...mapActions({
      fetchAlarmItem: 'fetchItem',
    }),

    fetchAlarmItemWithGroupsAndSteps(alarm) {
      const params = {
        with_steps: true,
        correlation: this.widget.parameters.isCorrelationEnabled || false,
      };

      if (!this.hideGroups) {
        if (alarm.causes) {
          params.with_causes = true;
        }

        if (alarm.consequences) {
          params.with_consequences = true;
        }
      }

      return this.fetchAlarmItemWithParams(alarm, params);
    },

    fetchAlarmItemWithParams(alarm, params) {
      const defaultParams = {
        sort_key: 't',
        sort_dir: SORT_ORDERS.desc.toLowerCase(),
        limit: 1,
      };

      if (alarm.v.resolved) {
        defaultParams.resolved = true;
      }

      return this.fetchAlarmItem({
        id: alarm._id,
        params: { ...defaultParams, ...params },
        dataPreparer: (d) => {
          const { alarms: fetchedAlarms = [] } = d.data[0];
          const [firstFetchedAlarm] = fetchedAlarms;

          if (alarm.filtered && firstFetchedAlarm) {
            return [{
              ...firstFetchedAlarm,

              filtered: alarm.filtered,
            }];
          }

          return fetchedAlarms;
        },
      });
    },
  },
};
