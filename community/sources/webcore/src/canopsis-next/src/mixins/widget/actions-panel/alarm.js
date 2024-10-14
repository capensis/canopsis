import { createNamespacedHelpers } from 'vuex';
import { find, pick } from 'lodash';

import {
  MODALS,
  BUSINESS_USER_PERMISSIONS_ACTIONS_MAP,
  LINK_RULE_ACTIONS,
  ALARM_LIST_ACTIONS_TYPES,
  ALARM_EXPORT_FILE_NAME_PREFIX,
} from '@/constants';

import { convertObjectToTreeview } from '@/helpers/treeview';
import { mapIds } from '@/helpers/array';
import { generatePreparedDefaultAlarmListWidget } from '@/helpers/entities/widget/form';
import { createEntityIdPatternByValue } from '@/helpers/entities/pattern/form';
import { mapAlarmsEntities } from '@/helpers/entities/alarm/form';

import { authMixin } from '@/mixins/auth';
import { queryMixin } from '@/mixins/query';
import { entitiesPbehaviorMixin } from '@/mixins/entities/pbehavior';
import { entitiesDeclareTicketRuleMixin } from '@/mixins/entities/declare-ticket-rule';
import { entitiesMetaAlarmMixin } from '@/mixins/entities/meta-alarm';
import { entitiesAlarmLinksMixin } from '@/mixins/entities/alarm/links';
import { clipboardMixin } from '@/mixins/clipboard';

const { mapActions: mapAlarmActions } = createNamespacedHelpers('alarm');

export const widgetActionsPanelAlarmMixin = {
  mixins: [
    authMixin,
    queryMixin,
    clipboardMixin,
    entitiesMetaAlarmMixin,
    entitiesPbehaviorMixin,
    entitiesAlarmLinksMixin,
    entitiesDeclareTicketRuleMixin,
  ],
  data() {
    return {
      ticketsForAlarmsPending: false,
      pendingByActionsTypes: {},
    };
  },
  methods: {
    ...mapAlarmActions({
      fetchResolvedAlarmsListWithoutStore: 'fetchResolvedAlarmsListWithoutStore',
      bulkCreateAlarmAckEvent: 'bulkCreateAlarmAckEvent',
      bulkCreateAlarmAckremoveEvent: 'bulkCreateAlarmAckremoveEvent',
      bulkCreateAlarmSnoozeEvent: 'bulkCreateAlarmSnoozeEvent',
      bulkCreateAlarmAssocticketEvent: 'bulkCreateAlarmAssocticketEvent',
      bulkCreateAlarmCommentEvent: 'bulkCreateAlarmCommentEvent',
      bulkCreateAlarmCancelEvent: 'bulkCreateAlarmCancelEvent',
      bulkCreateAlarmUnCancelEvent: 'bulkCreateAlarmUnCancelEvent',
      bulkCreateAlarmChangestateEvent: 'bulkCreateAlarmChangestateEvent',
      addBookmarkToAlarm: 'addBookmarkToAlarm',
      removeBookmarkFromAlarm: 'removeBookmarkFromAlarm',
    }),

    isActionTypeInPending(type) {
      return !!this.pendingByActionsTypes[type];
    },

    setActionPending(type, value) {
      this.$set(this.pendingByActionsTypes, type, value);
    },

    showCreateChangeStateEventModalByAlarms(alarms) {
      this.$modals.show({
        name: MODALS.createChangeStateEvent,
        config: {
          items: alarms,
          action: async (changeStateEvent) => {
            await this.bulkCreateAlarmChangestateEvent({
              data: alarms.map(alarm => ({ ...changeStateEvent, _id: alarm._id })),
            });

            await this.afterSubmit();
          },
        },
      });
    },

    showAckModalByAlarms(alarms) {
      this.$modals.show({
        name: MODALS.createAckEvent,
        config: {
          items: alarms,
          isNoteRequired: this.widget.parameters.isAckNoteRequired,
          action: async (ackEvent, { needDeclareTicket, needAssociateTicket }) => {
            await this.bulkCreateAlarmAckEvent({
              data: alarms.map(alarm => ({ ...ackEvent, _id: alarm._id })),
            });

            await this.afterSubmit();

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

    async createFastAckActionByAlarms(alarms) {
      this.setActionPending(ALARM_LIST_ACTIONS_TYPES.fastAck, true);

      try {
        const { fastAckOutput } = this.widget.parameters;

        const ackEvent = { comment: fastAckOutput?.enabled ? fastAckOutput.value : '' };

        await this.bulkCreateAlarmAckEvent({
          data: alarms.map(alarm => ({ ...ackEvent, _id: alarm._id })),
        });

        await this.afterSubmit();
      } catch (err) {
        console.error(err);
      } finally {
        this.setActionPending(ALARM_LIST_ACTIONS_TYPES.fastAck, false);
      }
    },

    async createFastCancelActionByAlarms(alarms) {
      this.setActionPending(ALARM_LIST_ACTIONS_TYPES.fastCancel, true);

      try {
        const { fastCancelOutput } = this.widget.parameters;

        const cancelEvent = { comment: fastCancelOutput?.enabled ? fastCancelOutput.value : '' };

        await this.bulkCreateAlarmCancelEvent({
          data: alarms.map(alarm => ({ ...cancelEvent, _id: alarm._id })),
        });

        this.afterSubmit();
      } catch (err) {
        console.error(err);
      } finally {
        this.setActionPending(ALARM_LIST_ACTIONS_TYPES.fastCancel, false);
      }
    },

    showCancelModalByAlarms(alarms) {
      this.$modals.show({
        name: MODALS.createEvent,
        config: {
          items: alarms,

          title: this.$t('modals.createCancelEvent.title'),
          action: async (cancelEvent) => {
            await this.bulkCreateAlarmCancelEvent({
              data: alarms.map(alarm => ({ ...cancelEvent, _id: alarm._id })),
            });

            await this.afterSubmit();
          },
        },
      });
    },

    showSnoozeModalByAlarms(alarms) {
      this.$modals.show({
        name: MODALS.createSnoozeEvent,
        config: {
          items: alarms,
          isNoteRequired: this.widget.parameters.isSnoozeNoteRequired,
          action: async (snoozeEvent) => {
            await this.bulkCreateAlarmSnoozeEvent({
              data: alarms.map(alarm => ({ ...snoozeEvent, _id: alarm._id })),
            });

            await this.afterSubmit();
          },
        },
      });
    },

    showCreateCommentModalByAlarms(alarms) {
      this.$modals.show({
        name: MODALS.createCommentEvent,
        config: {
          items: alarms,
          action: async (commentEvent) => {
            await this.bulkCreateAlarmCommentEvent({
              data: alarms.map(alarm => ({ ...commentEvent, _id: alarm._id })),
            });

            await this.afterSubmit();
          },
        },
      });
    },

    async showDeclareTicketModalByAlarms(alarms) {
      this.setActionPending(ALARM_LIST_ACTIONS_TYPES.declareTicket, true);

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
            action: (events, singleMode) => {
              this.$modals.show({
                name: MODALS.executeDeclareTickets,
                config: {
                  singleMode,
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
        this.setActionPending(ALARM_LIST_ACTIONS_TYPES.declareTicket, false);
      }
    },

    showAssociateTicketModalByAlarms(alarms, ignoreAck = false) {
      this.$modals.show({
        name: MODALS.createAssociateTicketEvent,
        config: {
          items: alarms,
          ignoreAck,
          action: async (associateEvent) => {
            if (!ignoreAck) {
              const itemsWithoutAck = alarms.filter(alarm => !alarm.v.ack);

              const { fastAckOutput } = this.widget.parameters;

              await this.bulkCreateAlarmAckEvent({
                data: itemsWithoutAck.map(alarm => ({
                  comment: fastAckOutput?.enabled ? fastAckOutput.value : '',
                  _id: alarm._id,
                })),
              });
            }

            await this.bulkCreateAlarmAssocticketEvent({
              data: alarms.map(alarm => ({ ...associateEvent, _id: alarm._id })),
            });

            this.afterSubmit();
          },
        },
      });
    },

    showAckRemoveModalByAlarms(alarms) {
      this.$modals.show({
        name: MODALS.createEvent,
        config: {
          items: alarms,
          title: this.$t('modals.createAckRemove.title'),
          action: async (ackRemoveEvent) => {
            await this.bulkCreateAlarmAckremoveEvent({
              data: alarms.map(alarm => ({ ...ackRemoveEvent, _id: alarm._id })),
            });

            await this.afterSubmit();
          },
        },
      });
    },

    showUnCancelModalByAlarms(alarms) {
      this.$modals.show({
        name: MODALS.createEvent,
        config: {
          items: alarms,
          title: this.$t('modals.createUnCancel.title'),
          isCommentRequired: this.widget.parameters.isUncancelAlarmsCommentRequired,
          action: async (unCancelEvent) => {
            await this.bulkCreateAlarmUnCancelEvent({
              data: alarms.map(alarm => ({ ...unCancelEvent, _id: alarm._id })),
            });

            await this.afterSubmit();
          },
        },
      });
    },

    showLinkToMetaAlarmModalByAlarms(alarms) {
      this.$modals.show({
        name: MODALS.linkToMetaAlarm,
        config: {
          items: alarms,

          title: this.$t('modals.linkToMetaAlarm.title'),
          action: async (data) => {
            if (data.id) {
              await this.addAlarmsIntoMetaAlarm({ id: data.id, data });
            } else {
              await this.createMetaAlarm({ data });
            }

            await this.afterSubmit();
          },
        },
      });
    },

    showRemoveAlarmsFromManualMetaAlarmModalByAlarms(alarms) {
      this.$modals.show({
        name: MODALS.removeAlarmsFromMetaAlarm,
        config: {
          items: alarms,
          title: this.$t('alarm.actions.titles.removeAlarmsFromManualMetaAlarm'),
          isCommentRequired: this.widget.parameters.isRemoveAlarmsFromMetaAlarmCommentRequired,
          action: async (data) => {
            await this.removeAlarmsFromMetaAlarm({
              id: this.parentAlarm?._id,
              data,
            });

            await this.afterSubmit();
          },
        },
      });
    },

    showRemoveAlarmsFromAutoMetaAlarmModalByAlarms(alarms) {
      this.$modals.show({
        name: MODALS.removeAlarmsFromMetaAlarm,
        config: {
          items: alarms,
          title: this.$t('alarm.actions.titles.removeAlarmsFromAutoMetaAlarm'),
          isCommentRequired: this.widget.parameters.isRemoveAlarmsFromMetaAlarmCommentRequired,
          action: async (data) => {
            await this.removeAlarmsFromMetaAlarm({
              id: this.parentAlarm?._id,
              data,
            });

            await this.afterSubmit();
          },
        },
      });
    },

    async handleLinkClickActionByAlarms(alarms, link, type) {
      try {
        this.setActionPending(type, true);

        const links = await this.fetchAlarmLinkWithoutStore({
          id: link.rule_id,
          params: { ids: mapIds(alarms) },
        });

        const summaryLink = find(links, pick(link, ['icon_name', 'label']));

        if (!summaryLink) {
          return;
        }

        if (link.action === LINK_RULE_ACTIONS.copy) {
          this.writeTextToClipboard(summaryLink.url);

          return;
        }

        window.open(summaryLink.url, '_blank');
      } catch (err) {
        console.error(err);
      } finally {
        this.setActionPending(type, false);
      }
    },

    showVariablesHelperModalByAlarm(alarm) {
      const {
        entity,
        pbehavior,
        infos,
        ...alarmOtherFields
      } = alarm;
      const variables = [{
        ...convertObjectToTreeview(alarmOtherFields, 'alarm'),

        original: this.item,
      }];

      if (entity) {
        variables.push(convertObjectToTreeview(entity, 'entity'));
      }

      if (pbehavior) {
        variables.push(convertObjectToTreeview(pbehavior, 'pbehavior'));
      }

      this.$modals.show({
        name: MODALS.variablesHelp,
        config: {
          exportEntity: alarm,
          exportEntityName: `${ALARM_EXPORT_FILE_NAME_PREFIX}-${alarm._id}`,
          variables,
        },
      });
    },

    showAddPbehaviorModalByAlarms(alarms) {
      this.$modals.show({
        name: MODALS.pbehaviorPlanning,
        config: {
          entityPattern: createEntityIdPatternByValue(
            alarms.length === 1 ? alarms[0].entity?._id : alarms.map(item => item.entity._id),
          ),
          entities: mapAlarmsEntities(alarms),
          afterSubmit: this.afterSubmit,
        },
      });
    },

    showHistoryModalByAlarm(alarm) {
      const widget = generatePreparedDefaultAlarmListWidget();

      widget.parameters.widgetColumns = this.widget.parameters.widgetColumns;
      widget.parameters.widgetGroupColumns = this.widget.parameters.widgetGroupColumns;
      widget.parameters.serviceDependenciesColumns = this.widget.parameters.serviceDependenciesColumns;

      this.$modals.show({
        name: MODALS.alarmsList,
        config: {
          widget,
          title: this.$t('modals.alarmsList.prefixTitle', { prefix: alarm.entity._id }),
          fetchList: params => this.fetchResolvedAlarmsListWithoutStore({
            params: { ...params, _id: alarm.entity._id },
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

    async addBookmarkByAlarm(alarm) {
      try {
        this.setActionPending(ALARM_LIST_ACTIONS_TYPES.addBookmark, true);

        await this.addBookmarkToAlarm({ id: alarm._id });

        this.$popups.success({ text: this.$t('alarm.popups.addBookmarkSuccess') });

        await this.afterSubmit();
      } catch (err) {
        console.error(err);

        this.$popups.error({ text: this.$t('alarm.popups.addBookmarkFailed') });
      } finally {
        this.setActionPending(ALARM_LIST_ACTIONS_TYPES.addBookmark, false);
      }
    },

    async removeBookmarkByAlarm(alarm) {
      try {
        this.setActionPending(ALARM_LIST_ACTIONS_TYPES.removeBookmark, true);

        await this.removeBookmarkFromAlarm({ id: alarm._id });

        this.$popups.success({ text: this.$t('alarm.popups.removeBookmarkSuccess') });

        await this.afterSubmit();
      } catch (err) {
        console.error(err);

        this.$popups.error({ text: this.$t('alarm.popups.removeBookmarkFailed') });
      } finally {
        this.setActionPending(ALARM_LIST_ACTIONS_TYPES.removeBookmark, false);
      }
    },
  },
};
