import { createNamespacedHelpers } from 'vuex';

import {
  MODALS,
  WEATHER_ACK_EVENT_OUTPUT,
  BUSINESS_USER_PERMISSIONS_ACTIONS_MAP,
  WEATHER_ACTIONS_TYPES,
  PBEHAVIOR_TYPE_TYPES,
} from '@/constants';

import { mapIds } from '@/helpers/entities';
import { isActionTypeAvailableForEntity } from '@/helpers/entities/entity';

import { authMixin } from '@/mixins/auth';
import { entitiesPbehaviorMixin } from '@/mixins/entities/pbehavior';
import { entitiesPbehaviorTypeMixin } from '@/mixins/entities/pbehavior/types';
import { entitiesDeclareTicketRuleMixin } from '@/mixins/entities/declare-ticket-rule';

const { mapActions: mapAlarmActions } = createNamespacedHelpers('alarm');

export const widgetActionPanelServiceEntityMixin = {
  mixins: [
    authMixin,
    entitiesPbehaviorTypeMixin,
    entitiesPbehaviorMixin,
    entitiesPbehaviorTypeMixin,
    entitiesDeclareTicketRuleMixin,
    entitiesDeclareTicketRuleMixin,
  ],
  data() {
    return {
      unavailableEntitiesAction: {},
      pendingByActionType: Object.values(WEATHER_ACTIONS_TYPES).reduce((acc, type) => {
        acc[type] = false;

        return acc;
      }, {}),
    };
  },
  computed: {
    /**
     * @return {Object.<string, Function>}
     */
    actionsMethodsMap() {
      return {
        [WEATHER_ACTIONS_TYPES.entityAck]: this.applyAckAction,
        [WEATHER_ACTIONS_TYPES.entityAckRemove]: this.showAckRemoveModal,
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
    ...mapAlarmActions({ fetchAlarmItemWithoutStore: 'fetchItemWithoutStore' }),

    fetchAlarmsByEntities(entities) {
      return Promise.all(
        entities.map(({ alarm_id: alarmId }) => this.fetchAlarmItemWithoutStore({ id: alarmId })),
      );
    },

    setActionPendingByType(type, value) {
      this.pendingByActionType[type] = value;
    },

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

    applyAckAction(entities) {
      this.applyAction({
        entities,
        actionType: WEATHER_ACTIONS_TYPES.entityAck,
      });
    },

    showAckRemoveModal(entities) {
      this.$modals.show({
        name: MODALS.textFieldEditor,
        config: {
          title: this.$t('modals.createAckRemove.title'),
          field: {
            name: 'output',
            label: this.$t('common.note'),
            validationRules: 'required',
          },
          action: ({ output }) => {
            this.applyAction({
              entities,
              payload: { output },
              actionType: WEATHER_ACTIONS_TYPES.entityAckRemove,
            });
          },
        },
      });
    },

    async showCreateAssociateTicketModal(entities) {
      this.setActionPendingByType(WEATHER_ACTIONS_TYPES.entityAssocTicket, true);

      try {
        const alarms = await this.fetchAlarmsByEntities(entities);

        this.$modals.show({
          name: MODALS.createAssociateTicketEvent,
          config: {
            items: alarms,
            action: (event) => {
              this.applyAction({
                entities,
                actionType: WEATHER_ACTIONS_TYPES.entityAssocTicket,
                payload: event,
              });
            },
          },
        });
      } catch (err) {
        console.error(err);
      } finally {
        this.setActionPendingByType(WEATHER_ACTIONS_TYPES.entityAssocTicket, false);
      }
    },

    async showCreateDeclareTicketModal(entities) {
      this.setActionPendingByType(WEATHER_ACTIONS_TYPES.declareTicket, true);

      try {
        const alarms = await this.fetchAlarmsByEntities(entities);

        const {
          by_rules: alarmsByTickets,
          by_alarms: ticketsByAlarms,
        } = await this.fetchAssignedDeclareTicketsWithoutStore({
          params: {
            alarms: mapIds(alarms),
          },
        });

        this.$modals.show({
          name: MODALS.createDeclareTicketEvent,
          config: {
            items: alarms,
            alarmsByTickets,
            ticketsByAlarms,
            action: (events) => {
              this.$modals.show({
                name: MODALS.executeDeclareTickets,
                config: {
                  executions: events,
                  tickets: events.map(({ _id: id }) => ({
                    _id: id,
                    name: alarmsByTickets[id].name,
                  })),
                  alarms,
                  onExecute: () => this.$emit('refresh'),
                },
              });
            },
          },
        });
      } catch (err) {
        console.error(err);
      } finally {
        this.setActionPendingByType(WEATHER_ACTIONS_TYPES.declareTicket, false);
      }
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
