import {
  EVENT_ENTITY_TYPES,
  ENTITIES_STATES,
  WEATHER_EVENT_DEFAULT_ENTITY,
  WEATHER_ACK_EVENT_OUTPUT,
  MAX_PBEHAVIOR_DEFAULT_TSTOP,
  WEATHER_ENTITY_PBEHAVIOR_DEFAULT_TITLE,
  PATTERN_CONDITIONS,
  ENTITY_PATTERN_FIELDS,
} from '@/constants';

/**
 * @typedef {Object} ServiceEvent
 * @property {Object} data
 * @property {string} type
 */

export default {
  methods: {
    /**
     * Add event to queue
     *
     * @param {ServiceEvent} event
     * @param {Object} entity
     */
    addEvent(event, entity) {
      this.$emit('add:event', { ...event, entityId: entity._id });
    },
    /**
     * Prepare weather entity data for event creation
     *
     * @param {string} eventType - type of entity
     * @param {Object} item - entity
     */
    prepareData(eventType, item) {
      return {
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
     * Call emit add:event for ack entity event
     *
     * @param {Object} entity
     * @param {string} output
     */
    addAckActionToQueue({ entity, output }) {
      const data = {
        ...this.prepareData(EVENT_ENTITY_TYPES.ack, entity),

        output,
      };

      this.addEvent({ type: EVENT_ENTITY_TYPES.ack, data }, entity);
    },

    /**
     * Call emit add:event for associate ticket event
     *
     * @param {Object} entity
     * @param {string} ticket
     */
    addAssocTicketActionToQueue({ entity, ticket }) {
      const data = {
        ...this.prepareData(EVENT_ENTITY_TYPES.assocTicket, entity),

        ticket,
      };

      this.addEvent({ type: EVENT_ENTITY_TYPES.assocTicket, data }, entity);
    },

    /**
     * Call emit add:event for validate entity event
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

      this.addEvent({ type: EVENT_ENTITY_TYPES.changeState, data }, entity);
    },

    /**
     * Call emit add:event for invalidate entity event
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

      this.addEvent({ type: EVENT_ENTITY_TYPES.cancel, data }, entity);
    },

    /**
     * Call emit add:event for pause entity event
     *
     * @param {Entity} entity
     * @param {string} comment
     * @param {Object} reason
     * @param {Object} type
     */
    addPauseActionToQueue({
      entity,
      comment,
      reason,
      type,
    }) {
      const data = {
        reason,
        type,
        comments: [{
          message: comment,
        }],
        entityPattern: [[{
          field: ENTITY_PATTERN_FIELDS.id,
          cond: {
            type: PATTERN_CONDITIONS.equal,
            value: entity._id,
          },
        }]],
        name: `${WEATHER_ENTITY_PBEHAVIOR_DEFAULT_TITLE}-${entity.name}-${Date.now()}`,
        tstart: new Date(),
        tstop: new Date(MAX_PBEHAVIOR_DEFAULT_TSTOP * 1000),
      };

      this.addEvent({ type: EVENT_ENTITY_TYPES.pause, data, entity }, entity);
    },

    /**
     * Call emit add:event for play entity event
     *
     * @param {Entity} entity
     */
    addPlayActionToQueue({ entity }) {
      this.addEvent({ type: EVENT_ENTITY_TYPES.play, data: entity }, entity);
    },

    /**
     * Call emit add:event for cancel entity event
     *
     * @param {Entity} entity
     * @param {string} output
     */
    addCancelActionToQueue({ entity, output }) {
      const data = {
        ...this.prepareData(EVENT_ENTITY_TYPES.cancel, entity),

        output,
      };

      this.addEvent({ type: EVENT_ENTITY_TYPES.cancel, data }, entity);
    },

    /**
     * Call emit add:event for comment entity event
     *
     * @param {Entity} entity
     * @param {string} output
     */
    addCommentActionToQueue({ entity, output }) {
      const data = {
        ...this.prepareData(EVENT_ENTITY_TYPES.comment, entity),
        output,
      };

      this.addEvent({ type: EVENT_ENTITY_TYPES.comment, data }, entity);
    },
  },
};
