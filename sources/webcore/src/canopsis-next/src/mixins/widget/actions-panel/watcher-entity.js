import {
  MODALS,
  EVENT_ENTITY_TYPES,
  WEATHER_ACK_EVENT_OUTPUT,
  BUSINESS_USER_RIGHTS_ACTIONS_MAP,
} from '@/constants';

import authMixin from '@/mixins/auth';
import modalMixin from '@/mixins/modal';
import eventActionsWatcherEntityMixin from '@/mixins/event-actions/watcher-entity';

export default {
  mixins: [authMixin, modalMixin, eventActionsWatcherEntityMixin],
  methods: {
    /**
     * Filter for available entity actions
     *
     * @param {string} type
     * @return boolean
     */
    actionsAccessFilterHandler({ type }) {
      const right = BUSINESS_USER_RIGHTS_ACTIONS_MAP.weather[type];

      if (!right) {
        return true;
      }

      return this.checkAccess(right);
    },

    prepareAckAction() {
      this.addAckActionToQueue({ entity: this.entity, output: WEATHER_ACK_EVENT_OUTPUT.ack });
      this.actionsClicked.push(EVENT_ENTITY_TYPES.ack);
    },

    prepareDeclareTicketAction() {
      this.showModal({
        name: MODALS.createWatcherDeclareTicketEvent,
        config: {
          action: (ticket) => {
            this.addAckActionToQueue({ entity: this.entity, output: WEATHER_ACK_EVENT_OUTPUT.ack });
            this.addDeclareTicketActionToQueue({ entity: this.entity, ticket });
            this.actionsClicked.push(EVENT_ENTITY_TYPES.declareTicket);
          },
        },
      });
    },

    prepareValidateAction() {
      this.addAckActionToQueue({ entity: this.entity, output: WEATHER_ACK_EVENT_OUTPUT.validateOk });
      this.addValidateActionToQueue({ entity: this.entity });
      this.actionsClicked.push(EVENT_ENTITY_TYPES.validate);
    },

    prepareInvalidateAction() {
      this.addAckActionToQueue({ entity: this.entity, output: WEATHER_ACK_EVENT_OUTPUT.ack });
      this.addInvalidateActionToQueue({ entity: this.entity });
      this.actionsClicked.push(EVENT_ENTITY_TYPES.invalidate);
    },

    preparePauseAction() {
      this.showModal({
        name: MODALS.createWatcherPauseEvent,
        config: {
          action: (pause) => {
            this.addPauseActionToQueue({
              entity: this.entity,
              comment: pause.comment,
              reason: pause.reason,
            });
            this.actionsClicked.push(EVENT_ENTITY_TYPES.pause);
          },
        },
      });
    },

    preparePlayAction() {
      this.addPlayActionToQueue({ entity: this.entity });
      this.actionsClicked.push(EVENT_ENTITY_TYPES.play);
    },
  },
};
