import { createNamespacedHelpers } from 'vuex';
import moment from 'moment';

import { EVENT_ENTITY_TYPES } from '@/constants';

import authMixin from './auth';

const { mapActions: eventMapActions } = createNamespacedHelpers('event');
const { mapActions: alarmMapActions } = createNamespacedHelpers('alarm');

/**
 * @mixin
 */
export default {
  mixins: [authMixin],
  methods: {
    ...eventMapActions({
      createEventAction: 'create',
    }),

    ...alarmMapActions({
      fetchAlarmsListWithPreviousParams: 'fetchListWithPreviousParams',
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

      if (this.config && this.config.afterSubmit) {
        return this.config.afterSubmit();
      }

      if (this.widget) {
        return this.fetchAlarmsListWithPreviousParams({ widgetId: this.widget._id });
      }

      return undefined;
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
        author: this.currentUser.crecord_name,
        id: item.id,
        connector: item.v.connector,
        connector_name: item.v.connector_name,
        source_type: item.entity ? item.entity.type : null,
        component: item.v.component,
        state: item.v.state.val,
        event_type: type,
        crecord_type: type,
        timestamp: moment().unix(),
        resource: item.v.resource,
        ref_rk: `${item.v.resource}/${item.v.component}`,
      };

      if (type !== EVENT_ENTITY_TYPES.snooze) {
        preparedData.state_type = item.v.status.val;
      }

      return [{ ...preparedData, ...data }];
    },
  },
};
