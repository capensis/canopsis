import {
  MODALS,
  EVENT_ENTITY_TYPES,
  WEATHER_ACK_EVENT_OUTPUT,
  BUSINESS_USER_RIGHTS_ACTIONS_MAP,
  WEATHER_AUTOREMOVE_BYPAUSE_OUTPUT, PBEHAVIOR_TYPE_TYPES,
} from '@/constants';

import authMixin from '@/mixins/auth';
import eventActionsWatcherEntityMixin from '@/mixins/event-actions/watcher-entity';
import entitiesPbehaviorMixin from '@/mixins/entities/pbehavior';
import entitiesPbehaviorTypesMixin from '@/mixins/entities/pbehavior/types';

export default {
  mixins: [
    authMixin,
    entitiesPbehaviorTypesMixin,
    eventActionsWatcherEntityMixin,
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
            this.actionsClicked.push(EVENT_ENTITY_TYPES.assocTicket);
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
            this.actionsClicked.push(EVENT_ENTITY_TYPES.comment);
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

    async preparePauseAction() {
      const defaultPbehaviorTypes = await this.fetchDefaultPbehaviorTypes();

      const pauseType = defaultPbehaviorTypes.find(({ type }) => type === PBEHAVIOR_TYPE_TYPES.pause);

      this.$modals.show({
        name: MODALS.createWatcherPauseEvent,
        config: {
          action: (pause) => {
            this.addPauseActionToQueue({
              entity: this.entity,
              comment: pause.comment,
              reason: pause.reason,
              type: pauseType,
            });
            this.addCancelActionToQueue({
              entity: this.entity,
              output: WEATHER_AUTOREMOVE_BYPAUSE_OUTPUT,
              fromSystem: true,
            });
            this.actionsClicked.push(EVENT_ENTITY_TYPES.pause);
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
      this.actionsClicked.push(EVENT_ENTITY_TYPES.play);
    },

    prepareCancelAction() {
      this.$modals.show({
        name: MODALS.textFieldEditor,
        config: {
          title: this.$t('modals.createCancelEvent.fields.output'),
          action: (output) => {
            this.addCancelActionToQueue({ entity: this.entity, output });
            this.actionsClicked.push(EVENT_ENTITY_TYPES.cancel);
          },
        },
      });
    },
  },
};
