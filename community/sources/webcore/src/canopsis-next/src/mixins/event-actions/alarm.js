import { isArray } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import { prepareEventByAlarm, prepareEventsByAlarms } from '@/helpers/forms/event';

import eventActionsMixin from './index';

const { mapActions: alarmMapActions } = createNamespacedHelpers('alarm');

/**
 * @mixin
 */
export default {
  mixins: [eventActionsMixin],
  methods: {
    ...alarmMapActions({
      fetchAlarmsListWithPreviousParams: 'fetchListWithPreviousParams',
    }),

    /**
     * Function calls dataPreparation and createEvent action and reload list of the entities
     *
     * @param {string} type - type of the event
     * @param {Alarm|Alarm[]} alarmOrAlarms - item of the entity
     * @param {Object} data - data for the event
     * @returns {Promise.<*>}
     */
    async createEvent(type, alarmOrAlarms, data = {}) {
      const eventData = isArray(alarmOrAlarms)
        ? prepareEventsByAlarms(type, alarmOrAlarms, data)
        : prepareEventByAlarm(type, alarmOrAlarms, data);

      await this.createEventAction({ data: eventData });

      if (this.config && this.config.afterSubmit) {
        return this.config.afterSubmit();
      }

      return Promise.resolve();
    },
  },
};
