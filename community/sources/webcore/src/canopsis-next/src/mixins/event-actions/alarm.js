import { isArray } from 'lodash';

import { prepareEventByAlarm, prepareEventsByAlarms } from '@/helpers/forms/event';

import { eventActionsMixin } from './index';

export const eventActionsAlarmMixin = {
  mixins: [eventActionsMixin],
  methods: {
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
