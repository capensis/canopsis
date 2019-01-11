import moment from 'moment';
import { createNamespacedHelpers } from 'vuex';

import { EVENT_ENTITY_TYPES, WEATHER_EVENT_DEFAULT_ENTITY, WEATHER_ACK_EVENT_OUTPUT } from '@/constants';

import authMixin from './auth';

const { mapActions } = createNamespacedHelpers('event');

export default {
  mixins: [authMixin],
  methods: {
    ...mapActions({
      createEventAction: 'create',
    }),

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

    addAckActionToQueue({ entity, output }) {
      const data = {
        ...this.prepareData(EVENT_ENTITY_TYPES.ack, entity),
        output,
      };

      this.$emit('addEvent', { type: EVENT_ENTITY_TYPES.ack, data });
    },

    addDeclareTicketActionToQueue({ entity, ticket }) {
      const data = {
        ...this.prepareData(EVENT_ENTITY_TYPES.declareTicket, entity),
        ticket,
      };

      this.$emit('addEvent', { type: EVENT_ENTITY_TYPES.declareTicket, data });
    },

    addValidateActionToQueue({ entity }) {
      const data = {
        ...this.prepareData(EVENT_ENTITY_TYPES.changeState, entity),
        state: 3,
        output: WEATHER_ACK_EVENT_OUTPUT.validateOk,
        keep_state: true,
      };

      this.$emit('addEvent', { type: EVENT_ENTITY_TYPES.changeState, data });
    },

    addInvalidateActionToQueue({ entity }) {
      const data = {
        ...this.prepareData(EVENT_ENTITY_TYPES.cancel, entity),
        state: 2,
        output: WEATHER_ACK_EVENT_OUTPUT.validateCancel,
        keep_state: true,
      };

      this.$emit('addEvent', { type: EVENT_ENTITY_TYPES.cancel, data });
    },

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
        tstop: 2147483647,
        type_: 'pause',
      };

      this.$emit('addEvent', { type: EVENT_ENTITY_TYPES.pause, data, entity });
    },

    addPlayActionToQueue({ entity }) {
      this.$emit('addEvent', { type: EVENT_ENTITY_TYPES.play, data: entity });
    },

  },
};
