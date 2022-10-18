import {
  MODALS,
  WEATHER_ACK_EVENT_OUTPUT,
  BUSINESS_USER_PERMISSIONS_ACTIONS_MAP,
  WEATHER_ACTIONS_TYPES,
  PBEHAVIOR_TYPE_TYPES,
} from '@/constants';

import { isActionTypeAvailableForEntity } from '@/helpers/entities/entity';

import { authMixin } from '@/mixins/auth';
import { entitiesPbehaviorMixin } from '@/mixins/entities/pbehavior';
import { entitiesPbehaviorTypeMixin } from '@/mixins/entities/pbehavior/types';

export const widgetActionPanelServiceEntityMixin = {
  mixins: [
    authMixin,
    entitiesPbehaviorTypeMixin,
    entitiesPbehaviorMixin,
    entitiesPbehaviorTypeMixin,
  ],
  data() {
    return {
      unavailableEntitiesAction: {},
    };
  },
  computed: {
    /**
     * @return {Object.<string, Function>}
     */
    actionsMethodsMap() {
      return {
        [WEATHER_ACTIONS_TYPES.entityAck]: this.addAckActionToQueue,
        [WEATHER_ACTIONS_TYPES.entityAssocTicket]: this.showCreateAssociateTicketModal,
        [WEATHER_ACTIONS_TYPES.entityValidate]: this.addValidateActionToQueue,
        [WEATHER_ACTIONS_TYPES.entityInvalidate]: this.addInvalidateActionToQueue,
        [WEATHER_ACTIONS_TYPES.entityPause]: this.showCreateServicePauseEventModal,
        [WEATHER_ACTIONS_TYPES.entityPlay]: this.addPlayActionToQueue,
        [WEATHER_ACTIONS_TYPES.entityCancel]: this.showCancelModal,
        [WEATHER_ACTIONS_TYPES.entityComment]: this.showCreateCommentEventModal,
        [WEATHER_ACTIONS_TYPES.declareTicket]: this.showCreateDeclareTicketModal,
      };
    },
  },
  methods: {
    removeEntityFromUnavailable(entity) {
      this.unavailableEntitiesAction[entity._id] = false;
    },

    applyAction(action) {
      const {
        availableEntities,
        unavailableEntities,
      } = action.entities.reduce((acc, entity) => {
        if (isActionTypeAvailableForEntity(action.actionType, entity)) {
          acc.availableEntities.push(entity);
        } else {
          acc.unavailableEntities.push(entity);
        }

        return acc;
      }, {
        availableEntities: [],
        unavailableEntities: [],
      });

      this.unavailableEntitiesAction = unavailableEntities.reduce((acc, { _id: id }) => {
        acc[id] = true;

        return acc;
      }, {});

      this.$emit('apply:action', {
        ...action,
        entities: availableEntities,
      });
    },

    applyEntityAction(actionType, entities) {
      const handler = this.actionsMethodsMap[actionType];

      if (handler) {
        handler(entities);
      }
    },

    /**
     * Filter for available entity actions
     *
     * @param {string} type
     * @return {boolean}
     */
    actionsAccessFilterHandler({ type }) {
      const permission = BUSINESS_USER_PERMISSIONS_ACTIONS_MAP.weather[type];

      return permission
        ? this.checkAccess(permission)
        : true;
    },

    addAckActionToQueue(entities) {
      this.applyAction({
        entities,
        actionType: WEATHER_ACTIONS_TYPES.entityAck,
      });
    },

    showCreateAssociateTicketModal(entities) {
      this.$modals.show({
        name: MODALS.textFieldEditor,
        config: {
          title: this.$t('modals.createAssociateTicket.title'),
          field: {
            name: 'ticket',
            label: this.$t('modals.createAssociateTicket.fields.ticket'),
            validationRules: 'required',
          },
          action: (ticket) => {
            this.applyAction({
              entities,
              actionType: WEATHER_ACTIONS_TYPES.entityAssocTicket,
              payload: { ticket },
            });
          },
        },
      });
    },

    showCreateDeclareTicketModal(entities) {
      this.applyAction({
        actionType: WEATHER_ACTIONS_TYPES.declareTicket,
        entities,
      });
    },

    showCreateCommentEventModal(entities) {
      this.$modals.show({
        name: MODALS.createCommentEvent,
        config: {
          action: ({ output }) => {
            this.applyAction({
              entities,
              actionType: WEATHER_ACTIONS_TYPES.entityComment,
              payload: { output },
            });
          },
        },
      });
    },

    addValidateActionToQueue(entities) {
      this.applyAction({
        actionType: WEATHER_ACTIONS_TYPES.entityValidate,
        entities,
      });
    },

    addInvalidateActionToQueue(entities) {
      this.applyAction({
        entities,
        actionType: WEATHER_ACTIONS_TYPES.entityInvalidate,
        payload: {
          output: WEATHER_ACK_EVENT_OUTPUT.ack,
        },
      });
    },

    showCreateServicePauseEventModal(entities) {
      this.$modals.show({
        name: MODALS.createServicePauseEvent,
        config: {
          action: async (pause) => {
            const defaultPbehaviorTypes = await this.fetchDefaultPbehaviorTypes();

            const pauseType = defaultPbehaviorTypes.find(({ type }) => type === PBEHAVIOR_TYPE_TYPES.pause);

            this.applyAction({
              entities,
              actionType: WEATHER_ACTIONS_TYPES.entityPause,
              payload: {
                comment: pause.comment,
                reason: pause.reason,
                type: pauseType,
              },
            });
          },
        },
      });
    },

    addPlayActionToQueue(entities) {
      this.applyAction({
        entities,
        actionType: WEATHER_ACTIONS_TYPES.entityPlay,
      });
    },

    showCancelModal(entities) {
      this.$modals.show({
        name: MODALS.textFieldEditor,
        config: {
          title: this.$t('common.note'),
          action: (output) => {
            this.applyAction({
              entities,
              actionType: WEATHER_ACTIONS_TYPES.entityCancel,
              payload: { output },
            });
          },
        },
      });
    },
  },
};
