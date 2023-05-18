import { createNamespacedHelpers } from 'vuex';
import { isEmpty, pick } from 'lodash';

import {
  MODALS,
  WEATHER_ACK_EVENT_OUTPUT,
  BUSINESS_USER_PERMISSIONS_ACTIONS_MAP,
  WEATHER_ACTIONS_TYPES,
  PBEHAVIOR_TYPE_TYPES, PBEHAVIOR_ORIGINS,
} from '@/constants';

import { mapIds } from '@/helpers/array';
import { isActionTypeAvailableForEntity } from '@/helpers/entities/entity/actions';
import { createDowntimePbehavior } from '@/helpers/entities/pbehavior/form';

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
    ...mapAlarmActions({
      fetchAlarmItemWithoutStore: 'fetchItemWithoutStore',
      bulkCreateAlarmAckEvent: 'bulkCreateAlarmAckEvent',
      bulkCreateAlarmAckremoveEvent: 'bulkCreateAlarmAckremoveEvent',
      bulkCreateAlarmSnoozeEvent: 'bulkCreateAlarmSnoozeEvent',
      bulkCreateAlarmAssocticketEvent: 'bulkCreateAlarmAssocticketEvent',
      bulkCreateAlarmCommentEvent: 'bulkCreateAlarmCommentEvent',
      bulkCreateAlarmCancelEvent: 'bulkCreateAlarmCancelEvent',
    }),

    refreshEntities() {
      this.$emit('refresh');
    },

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

    getAvailableActions(type, entities) {
      const {
        availableEntities,
        unavailableEntities,
      } = entities.reduce((acc, entity) => {
        if (isActionTypeAvailableForEntity(type, entity)) {
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

      return availableEntities;
    },

    applyEntityAction(actionType, entities) {
      const handler = this.actionsMethodsMap[actionType];

      if (handler) {
        handler(entities);
      }
    },

    showPbehaviorResponseErrorPopups(response) {
      if (response?.length) {
        response.forEach(({ error, errors }) => {
          if (error || !isEmpty(errors)) {
            this.$popups.error({ text: error || Object.values(errors).join('\n') });
          }
        });
      }
    },

    async createPauseEvent(entities, payload) {
      const response = await this.createEntityPbehaviors({
        data: entities.map(entity => createDowntimePbehavior({
          entity,
          ...pick(payload, ['comment', 'reason', 'type']),
        }), []),
      });

      this.showPbehaviorResponseErrorPopups(response);
    },

    async createPlayEvent(entities) {
      const response = await this.removeEntityPbehaviors({
        data: entities.map(({ _id: id }) => ({
          entity: id,
          origin: PBEHAVIOR_ORIGINS.serviceWeather,
        })),
      });

      this.showPbehaviorResponseErrorPopups(response);
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

    async applyAckAction(entities) {
      this.setActionPendingByType(WEATHER_ACTIONS_TYPES.entityAck, true);

      const availableEntities = this.getAvailableActions(WEATHER_ACTIONS_TYPES.entityAck, entities);

      const preparedRequestData = availableEntities.map(
        ({ alarm_id: alarmId }) => ({ _id: alarmId, comment: WEATHER_ACK_EVENT_OUTPUT.ack }),
      );

      await this.bulkCreateAlarmAckEvent({ data: preparedRequestData });

      this.refreshEntities();

      this.setActionPendingByType(WEATHER_ACTIONS_TYPES.entityAck, false);
    },

    showAckRemoveModal(entities) {
      this.$modals.show({
        name: MODALS.textFieldEditor,
        config: {
          title: this.$t('modals.createAckRemove.title'),
          field: {
            name: 'comment',
            label: this.$t('common.note'),
            validationRules: 'required',
          },
          action: async (comment) => {
            const availableEntities = this.getAvailableActions(WEATHER_ACTIONS_TYPES.entityAckRemove, entities);

            const preparedRequestData = availableEntities.map(
              ({ alarm_id: alarmId }) => ({ _id: alarmId, comment }),
            );

            await this.bulkCreateAlarmAckremoveEvent({ data: preparedRequestData });

            this.refreshEntities();
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
            action: async (event) => {
              const availableEntities = this.getAvailableActions(WEATHER_ACTIONS_TYPES.entityAssocTicket, entities);

              const preparedRequestData = availableEntities.map(
                ({ alarm_id: alarmId }) => ({ _id: alarmId, ...event, comment: WEATHER_ACK_EVENT_OUTPUT.ack }),
              );

              await this.bulkCreateAlarmAssocticketEvent({ data: preparedRequestData });

              this.refreshEntities();
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
                  onExecute: () => this.refreshEntities(),
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
          action: async ({ comment }) => {
            const availableEntities = this.getAvailableActions(WEATHER_ACTIONS_TYPES.entityComment, entities);

            const preparedRequestData = availableEntities.map(({ alarm_id: alarmId }) => ({ _id: alarmId, comment }));

            await this.bulkCreateAlarmCommentEvent({ data: preparedRequestData });

            this.refreshEntities();
          },
        },
      });
    },

    addValidateActionToQueue(/*  entities */) {
      /**
       * TODO: Event should be changed
       */
      /*
      this.applyAction({
        actionType: WEATHER_ACTIONS_TYPES.entityValidate,
        entities,
        payload: { comment: WEATHER_ACK_EVENT_OUTPUT.ack },
      });
       */
    },

    addInvalidateActionToQueue(/*  entities */) {
      /**
       * TODO: Event should be changed
       */
      /*
      this.applyAction({
        entities,
        actionType: WEATHER_ACTIONS_TYPES.entityInvalidate,
        payload: {
          comment: WEATHER_ACK_EVENT_OUTPUT.ack,
        },
      });
       */
    },

    showCreateServicePauseEventModal(entities) {
      this.$modals.show({
        name: MODALS.createServicePauseEvent,
        config: {
          action: async (pause) => {
            const defaultPbehaviorTypes = await this.fetchDefaultPbehaviorTypes();

            const pauseType = defaultPbehaviorTypes.find(({ type }) => type === PBEHAVIOR_TYPE_TYPES.pause);

            await this.createPauseEvent(entities, {
              comment: pause.comment,
              reason: pause.reason,
              type: pauseType,
            });

            this.refreshEntities();
          },
        },
      });
    },

    async addPlayActionToQueue(entities) {
      this.setActionPendingByType(WEATHER_ACTIONS_TYPES.entityPlay, true);

      await this.createPlayEvent(entities);

      this.refreshEntities();

      this.setActionPendingByType(WEATHER_ACTIONS_TYPES.entityPlay, false);
    },

    showCancelModal(entities) {
      this.$modals.show({
        name: MODALS.textFieldEditor,
        config: {
          title: this.$t('common.note'),
          action: async (comment) => {
            const availableEntities = this.getAvailableActions(WEATHER_ACTIONS_TYPES.entityCancel, entities);

            const preparedRequestData = availableEntities.map(
              ({ alarm_id: alarmId }) => ({ _id: alarmId, comment }),
            );

            await this.bulkCreateAlarmCancelEvent({ data: preparedRequestData });

            this.refreshEntities();
          },
        },
      });
    },
  },
};
