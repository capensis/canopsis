import { createNamespacedHelpers } from 'vuex';

import { SORT_ORDERS } from '@/constants';

import { queryMixin } from '@/mixins/query';

const { mapActions } = createNamespacedHelpers('alarm');

export const widgetExpandPanelAlarmMixin = {
  mixins: [queryMixin],
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
      const { correlation = false } = this.getQueryById(this.widget._id);
      const params = {
        correlation,

        with_steps: true,
      };

      if (!this.hideGroups && correlation) {
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
        with_instructions: true,
      };

      if (alarm.v.resolved) {
        defaultParams.resolved = true;
      }

      return this.fetchAlarmItem({
        id: alarm._id,
        params: { ...defaultParams, ...params },
        dataPreparer: ({ data: fetchedAlarms = [] }) => {
          const [firstFetchedAlarm] = fetchedAlarms;

          if (alarm.filtered_children && firstFetchedAlarm) {
            return [{
              ...firstFetchedAlarm,

              filtered_children: alarm.filtered_children,
            }];
          }

          return fetchedAlarms;
        },
      });
    },
  },
};
