import { createNamespacedHelpers } from 'vuex';

import { EVENT_TYPES } from '@/config';

const { mapActions: eventMapActions } = createNamespacedHelpers('event');
const { mapActions: alarmMapActions } = createNamespacedHelpers('entities/alarm');

export default {
  methods: {
    ...eventMapActions({
      createEventAction: 'create',
    }),

    ...alarmMapActions({
      fetchAlarmListWithPreviousParams: 'fetchListWithPreviousParams',
    }),

    async createEvent(type, item, data) {
      await this.createEventAction({ data: this.prepareData(type, item, data) });

      return this.fetchAlarmListWithPreviousParams(); // TODO: check items type for correct request
    },

    prepareData(type, item, data = {}) {
      /* if (Array.isArray(data)) {
        return data.reduce(
          (acc, dataPortion) => acc.concat(this.prepareData(eventType, dataPortion)),
          [],
        );
      } */

      const preparedData = {
        author: 'root',
        id: item.id,
        connector: item.v.connector,
        connector_name: item.v.connector_name,
        source_type: item.entity.type,
        component: item.v.component,
        state: item.v.state.val,
        event_type: type,
        crecord_type: type,
        timestamp: Date.now(),
        resource: item.v.resource,
        ref_rk: `${item.v.resource}/${item.v.component}`,
      };

      if (type !== EVENT_TYPES.snooze) {
        preparedData.state_type = item.v.status.val;
      }

      return [{ ...preparedData, ...data }];
    },
  },
};
