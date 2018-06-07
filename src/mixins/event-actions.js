import { createNamespacedHelpers } from 'vuex';

import { EVENT_TYPES } from '@/config';

const { mapActions: eventMapActions } = createNamespacedHelpers('event');
const { mapActions: alarmMapActions } = createNamespacedHelpers('alarm');

/**
 * @mixin
 */
export default {
  methods: {
    ...eventMapActions({
      createEventAction: 'create',
    }),

    ...alarmMapActions({
      fetchAlarmListWithPreviousParams: 'fetchListWithPreviousParams',
    }),

    /**
     * Function calls dataPreparation and createEvent action and reload list of the entities
     *
     * @param {string} type - type of the event
     * @param {Object} item - item of the entity
     * @param {Object} data - data for the event
     * @returns {Promise.<*>}
     */
    async createEvent(type, item, data) {
      await this.createEventAction({ data: this.prepareData(type, item, data) });

      return this.fetchAlarmListWithPreviousParams(); // TODO: check items type for correct request
    },

    /**
     * Function for data preparation
     *
     * @param {string} type - type of the event
     * @param {Object|Array} item - item of the entity | Array of items of entity
     * @param {Object} data - data for the event
     * @returns {Object[]}
     */
    prepareData(type, item, data = {}) {
      if (Array.isArray(item)) {
        return item.reduce((acc, value) => acc.concat(this.prepareData(type, value, data)), []);
      }

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
