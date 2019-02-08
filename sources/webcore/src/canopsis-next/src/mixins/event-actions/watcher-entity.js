import moment from 'moment';

import {
  EVENT_ENTITY_TYPES,
  ENTITIES_STATES,
  WEATHER_EVENT_DEFAULT_ENTITY,
  WEATHER_ACK_EVENT_OUTPUT,
} from '@/constants';

export default {
  methods: {
    /**
     * Prepare weather entity data for event creation
     *
     * @param {string} eventType - type of entity
     * @param {Object} item - entity
     */
    prepareData(eventType, item) {
      return {
        author: this.currentUser.crecord_name,
        component: item.component || WEATHER_EVENT_DEFAULT_ENTITY,
        connector: item.connector || WEATHER_EVENT_DEFAULT_ENTITY,
        connector_name: item.connector_name || WEATHER_EVENT_DEFAULT_ENTITY,
        crecord_type: eventType,
        event_type: eventType,
        ref_rk: `${item.resource || WEATHER_EVENT_DEFAULT_ENTITY}/${item.component || WEATHER_EVENT_DEFAULT_ENTITY}`,
        resource: item.resource || WEATHER_EVENT_DEFAULT_ENTITY,
        source_type: item.source_type,
      };
    },

    /**
     * Call emit addEvent for ack entity event
     *
     * @param {Object} entity
     * @param {string} output
     */
    addAckActionToQueue({ entity, output }) {
      const data = {
        ...this.prepareData(EVENT_ENTITY_TYPES.ack, entity),

        output,
      };

      this.$emit('addEvent', { type: EVENT_ENTITY_TYPES.ack, data });
    },

    /**
     * Call emit addEvent for declare entity ticker event
     *
     * @param {Object} entity
     * @param {string} ticket
     */
    addDeclareTicketActionToQueue({ entity, ticket }) {
      const data = {
        ...this.prepareData(EVENT_ENTITY_TYPES.declareTicket, entity),

        ticket,
      };

      this.$emit('addEvent', { type: EVENT_ENTITY_TYPES.declareTicket, data });
    },

    /**
     * Call emit addEvent for validate entity event
     *
     * @param {Object} entity
     */
    addValidateActionToQueue({ entity }) {
      const data = {
        ...this.prepareData(EVENT_ENTITY_TYPES.changeState, entity),

        state: ENTITIES_STATES.critical,
        output: WEATHER_ACK_EVENT_OUTPUT.validateOk,
        keep_state: true,
      };

      this.$emit('addEvent', { type: EVENT_ENTITY_TYPES.changeState, data });
    },

    /**
     * Call emit addEvent for invalidate entity event
     *
     * @param {Object} entity
     */
    addInvalidateActionToQueue({ entity }) {
      const data = {
        ...this.prepareData(EVENT_ENTITY_TYPES.cancel, entity),

        state: ENTITIES_STATES.major,
        output: WEATHER_ACK_EVENT_OUTPUT.validateCancel,
        keep_state: true,
      };

      this.$emit('addEvent', { type: EVENT_ENTITY_TYPES.cancel, data });
    },

    /**
     * Call emit addEvent for pause entity event
     *
     * @param {Object} entity
     * @param {string} comment
     * @param {string} reason
     */
    addPauseActionToQueue({ entity, comment, reason }) {
      const data = {
        author: this.currentUser.crecord_name,
        comments: [{
          author: this.currentUser.crecord_name,
          message: comment,
        }],
        filter: {
          _id: entity.entity_id,
        },
        name: 'downtime',
        reason,
        tstart: moment().unix(),
        tstop: 2147483647, // 01/19/2038 @ 3:14am (UTC)
        type_: 'pause',
      };

      this.$emit('addEvent', { type: EVENT_ENTITY_TYPES.pause, data, entity });
    },

    /**
     * Call emit addEvent for play entity event
     *
     * @param {Object} entity
     */
    addPlayActionToQueue({ entity }) {
      this.$emit('addEvent', { type: EVENT_ENTITY_TYPES.play, data: entity });
    },
  },
};
