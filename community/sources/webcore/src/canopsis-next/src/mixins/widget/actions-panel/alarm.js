import { createNamespacedHelpers } from 'vuex';

import {
  MODALS,
  EVENT_ENTITY_TYPES,
  BUSINESS_USER_PERMISSIONS_ACTIONS_MAP,
} from '@/constants';

import { convertObjectToTreeview } from '@/helpers/treeview';
import { generateDefaultAlarmListWidget, mapIds } from '@/helpers/entities';
import { createEntityIdPatternByValue } from '@/helpers/pattern';
import { prepareEventsByAlarms } from '@/helpers/forms/event';

import { authMixin } from '@/mixins/auth';
import { queryMixin } from '@/mixins/query';
import { eventActionsAlarmMixin } from '@/mixins/event-actions/alarm';
import { entitiesPbehaviorMixin } from '@/mixins/entities/pbehavior';
import { entitiesDeclareTicketRuleMixin } from '@/mixins/entities/declare-ticket-rule';

const { mapActions } = createNamespacedHelpers('alarm');

export const widgetActionsPanelAlarmMixin = {
  mixins: [
    authMixin,
    queryMixin,
    eventActionsAlarmMixin,
    entitiesPbehaviorMixin,
    entitiesDeclareTicketRuleMixin,
  ],
  data() {
    return {
      ticketsForAlarmsPending: false,
    };
  },
  methods: {
    ...mapActions({
      fetchResolvedAlarmsListWithoutStore: 'fetchResolvedAlarmsListWithoutStore',
    }),

    async createFastAckEvent() {
      let eventData = {};

      if (this.widget.parameters.fastAckOutput && this.widget.parameters.fastAckOutput.enabled) {
        eventData = { output: this.widget.parameters.fastAckOutput.value };
      }

      await this.createEvent(EVENT_ENTITY_TYPES.ack, this.item, eventData);

      return this.refreshAlarmsList();
    },

    showCreateCommentModal() {
      this.$modals.show({
        name: MODALS.createCommentEvent,
        config: {
          ...this.modalConfig,
          action: data => this.createEvent(EVENT_ENTITY_TYPES.comment, this.item, data),
        },
      });
    },

    showActionModal(name) {
      return () => this.$modals.show({
        name,
        config: this.modalConfig,
      });
    },

    showDeclareTicketModal() {
      this.showDeclareTicketModalByAlarms([this.item]);
    },

    async showDeclareTicketModalByAlarms(alarms) {
      this.ticketsForAlarmsPending = true;

      try {
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
            ...this.modalConfig,
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
                },
              });
            },
          },
        });
      } catch (err) {
        console.error(err);
      } finally {
        this.ticketsForAlarmsPending = false;
      }
    },

    showAssociateTicketModal() {
      this.showAssociateTicketModalByAlarms([this.item]);
    },

    showAssociateTicketModalByAlarms(alarms) {
      this.$modals.show({
        name: MODALS.createAssociateTicketEvent,
        config: {
          items: alarms,
          action: async (event) => {
            const itemsWithoutAck = alarms.filter(alarm => !alarm.v.ack);

            const { fastAckOutput } = this.widget.parameters;

            const ackEvents = prepareEventsByAlarms(
              EVENT_ENTITY_TYPES.ack,
              itemsWithoutAck,
              { output: fastAckOutput?.enabled ? fastAckOutput.value : '' },
            );

            const assocTicketEvents = prepareEventsByAlarms(EVENT_ENTITY_TYPES.assocTicket, alarms, event);

            await this.createEventAction({ data: [...ackEvents, ...assocTicketEvents] });

            this.afterSubmit();
          },
        },
      });
    },

    showSnoozeModal() {
      this.$modals.show({
        name: MODALS.createSnoozeEvent,
        config: {
          ...this.modalConfig,
          isNoteRequired: this.widget.parameters.isSnoozeNoteRequired,
        },
      });
    },

    showAckModal() {
      this.$modals.show({
        name: MODALS.createAckEvent,
        config: {
          ...this.modalConfig,

          isNoteRequired: this.widget.parameters.isAckNoteRequired,
        },
      });
    },

    showCancelEventModal() {
      this.$modals.show({
        name: MODALS.createEvent,
        config: {
          ...this.modalConfig,

          title: this.$t('modals.createCancelEvent.title'),
          eventType: EVENT_ENTITY_TYPES.cancel,
        },
      });
    },

    showAckRemoveModal() {
      this.$modals.show({
        name: MODALS.createEvent,
        config: {
          ...this.modalConfig,

          title: this.$t('modals.createAckRemove.title'),
          eventType: EVENT_ENTITY_TYPES.ackRemove,
        },
      });
    },

    showVariablesHelperModal() {
      const {
        entity,
        pbehavior,
        infos,
        ...alarm
      } = this.item;
      const variables = [];

      variables.push(convertObjectToTreeview(alarm, 'alarm'));

      if (entity) {
        variables.push(convertObjectToTreeview(entity, 'entity'));
      }

      if (pbehavior) {
        variables.push(convertObjectToTreeview(pbehavior, 'pbehavior'));
      }

      this.$modals.show({
        name: MODALS.variablesHelp,
        config: {
          ...this.modalConfig,

          variables,
        },
      });
    },

    showAddPbehaviorModal() {
      this.$modals.show({
        name: MODALS.pbehaviorPlanning,
        config: {
          entityPattern: createEntityIdPatternByValue(this.item.entity._id),
        },
      });
    },

    showHistoryModal() {
      const widget = generateDefaultAlarmListWidget();

      widget.parameters.widgetColumns = this.widget.parameters.widgetColumns;

      this.$modals.show({
        name: MODALS.alarmsList,
        config: {
          widget,
          title: this.$t('modals.alarmsList.prefixTitle', { prefix: this.item.entity._id }),
          fetchList: params => this.fetchResolvedAlarmsListWithoutStore({
            params: { ...params, _id: this.item.entity._id },
          }),
        },
      });
    },

    showManualMetaAlarmUngroupModal() {
      this.$modals.show({
        name: MODALS.createEvent,
        config: {
          ...this.modalConfig,

          title: this.$t('alarm.actions.titles.manualMetaAlarmUngroup'),
          eventType: EVENT_ENTITY_TYPES.manualMetaAlarmUngroup,
          parentsIds: [this.parentAlarm.d],
        },
      });
    },

    actionsAccessFilterHandler({ type }) {
      const permission = BUSINESS_USER_PERMISSIONS_ACTIONS_MAP.alarmsList[type];

      if (!permission) {
        return true;
      }

      return this.checkAccess(permission);
    },
  },
};
