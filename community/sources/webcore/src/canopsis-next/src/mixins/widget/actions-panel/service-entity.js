import {
  MODALS,
  WEATHER_ACK_EVENT_OUTPUT,
  BUSINESS_USER_PERMISSIONS_ACTIONS_MAP,
  PBEHAVIOR_TYPE_TYPES,
} from '@/constants';

import { authMixin } from '@/mixins/auth';
import eventActionsServiceEntityMixin from '@/mixins/event-actions/service-entity';
import { entitiesPbehaviorMixin } from '@/mixins/entities/pbehavior';
import entitiesPbehaviorTypesMixin from '@/mixins/entities/pbehavior/types';

export default {
  mixins: [
    authMixin,
    entitiesPbehaviorTypesMixin,
    eventActionsServiceEntityMixin,
    entitiesPbehaviorMixin,
  ],
  methods: {
    /**
     * Filter for available entity actions
     *
     * @param {string} type
     * @return boolean
     */
    actionsAccessFilterHandler({ type }) {
      const permission = BUSINESS_USER_PERMISSIONS_ACTIONS_MAP.weather[type];

      if (!permission) {
        return true;
      }

      return this.checkAccess(permission);
    },

    prepareAckAction() {
      this.addAckActionToQueue({ entity: this.entity, output: WEATHER_ACK_EVENT_OUTPUT.ack });
    },

    prepareAssocTicketAction() {
      this.$modals.show({
        name: MODALS.textFieldEditor,
        config: {
          title: this.$t('modals.createDeclareTicket.title'),
          field: {
            name: 'ticket',
            label: this.$t('modals.createAssociateTicket.fields.ticket'),
            validationRules: 'required',
          },
          action: (ticket) => {
            this.addAckActionToQueue({ entity: this.entity, output: WEATHER_ACK_EVENT_OUTPUT.ack });
            this.addAssocTicketActionToQueue({ entity: this.entity, ticket });
          },
        },
      });
    },

    prepareCommentAction() {
      this.$modals.show({
        name: MODALS.createCommentEvent,
        config: {
          action: ({ output }) => {
            this.addCommentActionToQueue({ entity: this.entity, output });
          },
        },
      });
    },

    prepareValidateAction() {
      this.addAckActionToQueue({ entity: this.entity, output: WEATHER_ACK_EVENT_OUTPUT.validateOk });
      this.addValidateActionToQueue({ entity: this.entity });
    },

    prepareInvalidateAction() {
      this.addAckActionToQueue({ entity: this.entity, output: WEATHER_ACK_EVENT_OUTPUT.ack });
      this.addInvalidateActionToQueue({ entity: this.entity });
    },

    preparePauseAction() {
      this.$modals.show({
        name: MODALS.createServicePauseEvent,
        config: {
          action: async (pause) => {
            const defaultPbehaviorTypes = await this.fetchDefaultPbehaviorTypes();

            const pauseType = defaultPbehaviorTypes.find(({ type }) => type === PBEHAVIOR_TYPE_TYPES.pause);

            this.addPauseActionToQueue({
              entity: this.entity,
              comment: pause.comment,
              reason: pause.reason,
              type: pauseType,
            });
          },
        },
      });
    },

    async fetchDefaultPbehaviorTypes() {
      const { data } = await this.fetchPbehaviorTypesListWithoutStore({
        params: { default: true },
      });

      return data;
    },

    preparePlayAction() {
      this.addPlayActionToQueue({ entity: this.entity });
    },

    prepareCancelAction() {
      this.$modals.show({
        name: MODALS.textFieldEditor,
        config: {
          title: this.$t('common.output'),
          action: (output) => {
            this.addCancelActionToQueue({ entity: this.entity, output });
          },
        },
      });
    },
  },
};
