import { createNamespacedHelpers } from 'vuex';

import {
  MODALS,
  EVENT_ENTITY_TYPES,
  BUSINESS_USER_PERMISSIONS_ACTIONS_MAP,
} from '@/constants';

import { convertObjectToTreeview } from '@/helpers/treeview';

import { generatePreparedDefaultAlarmListWidget, mapIds } from '@/helpers/entities';
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

    showActionModal(name) {
      return () => this.$modals.show({
        name,
        config: this.modalConfig,
      });
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
                  onExecute: this.afterSubmit,
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

    showAssociateTicketModalByAlarms(alarms, ignoreAck = false) {
      this.$modals.show({
        name: MODALS.createAssociateTicketEvent,
        config: {
          items: alarms,
          ignoreAck,
          action: async (event) => {
            const events = [];

            if (!ignoreAck) {
              const itemsWithoutAck = alarms.filter(alarm => !alarm.v.ack);

              const { fastAckOutput } = this.widget.parameters;

              events.push(...prepareEventsByAlarms(
                EVENT_ENTITY_TYPES.ack,
                itemsWithoutAck,
                { output: fastAckOutput?.enabled ? fastAckOutput.value : '' },
              ));
            }

            events.push(...prepareEventsByAlarms(EVENT_ENTITY_TYPES.assocTicket, alarms, event));

            await this.createEventAction({ data: events });

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
      this.showAckModalByAlarms([this.item]);
    },

    showAckModalByAlarms(alarms) {
      this.$modals.show({
        name: MODALS.createAckEvent,
        config: {
          items: alarms,
          isNoteRequired: this.widget.parameters.isAckNoteRequired,
          action: async (event, { needDeclareTicket, needAssociateTicket }) => {
            const ackEvents = prepareEventsByAlarms(
              EVENT_ENTITY_TYPES.ack,
              alarms,
              event,
            );

            await this.createEventAction({ data: ackEvents });

            await this.$emit('clear:items');
            await this.refreshAlarmsList();

            if (needAssociateTicket) {
              this.showAssociateTicketModalByAlarms(alarms, true);
            } else if (needDeclareTicket) {
              const alarmsWithRules = alarms.filter(
                ({ assigned_declare_ticket_rules: assignedDeclareTicketRules }) => assignedDeclareTicketRules?.length,
              );
              await this.showDeclareTicketModalByAlarms(alarmsWithRules);
            }
          },
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
      const widget = generatePreparedDefaultAlarmListWidget();

      widget.parameters.widgetColumns = this.widget.parameters.widgetColumns;
      widget.parameters.widgetGroupColumns = this.widget.parameters.widgetGroupColumns;
      widget.parameters.serviceDependenciesColumns = this.widget.parameters.serviceDependenciesColumns;

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

    actionsAccessFilterHandler({ type }) {
      const permission = BUSINESS_USER_PERMISSIONS_ACTIONS_MAP.alarmsList[type];

      if (!permission) {
        return true;
      }

      return this.checkAccess(permission);
    },
  },
};
