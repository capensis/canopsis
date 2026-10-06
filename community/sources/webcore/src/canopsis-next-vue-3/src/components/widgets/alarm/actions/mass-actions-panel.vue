<template>
  <actions-panel
    :actions="preparedActions"
    :inline-count="inlineCount"
    :small="small"
    ignore-media-query
  />
</template>

<script>
import { difference, keyBy } from 'lodash';
import { mapActions } from 'pinia';

import { VUETIFY_ANIMATION_DELAY } from '@/config';
import {
  ALARM_LIST_ACTIONS_TYPES,
  ALARM_LIST_FAST_PBEHAVIOR_TIMEOUT,
  ALARM_LIST_TOGGLE_ACTIONS_TYPES_MAP,
  BUSINESS_USER_PERMISSIONS_ACTIONS_MAP,
  FEATURES_SERVICE_KEYS,
  LINK_RULE_ACTIONS,
  MODALS,
  PBEHAVIOR_ORIGINS,
  PBEHAVIOR_TYPE_TYPES,
} from '@/constants';

import { featuresService } from '@/services/features';

import { useAlarmStore } from '@/stores/modules/alarm';
import { useDeclareTicketRuleStore } from '@/stores/modules/declare-ticket-rule';
import { useMetaAlarmStore } from '@/stores/modules/meta-alarm';
import { usePbehaviorStore } from '@/stores/modules/pbehavior';
import { usePbehaviorReasonStore } from '@/stores/modules/pbehavior-reason';
import { usePbehaviorTypeStore } from '@/stores/modules/pbehavior-type';

import { sortActionsByQuickActions, getActionsInlineCount } from '@/helpers/actions-panel';
import { mapIds } from '@/helpers/array';
import { promisedWait } from '@/helpers/async';
import {
  formToChangeStateEvent,
  isAlarmStateOk,
  isAlarmStatusCancelled,
  isAlarmStatusClosed,
  isAlarmStatusUnknown,
  isResolvedAlarm,
  mapAlarmsEntities,
} from '@/helpers/entities/alarm/form';
import { getAlarmActionIcon } from '@/helpers/entities/alarm/icons';
import {
  isResumeAvailableForEntity,
  isEntityPaused,
} from '@/helpers/entities/entity/pause';
import {
  getLinkRuleLinkActionType,
  harmonizeAlarmsLinks,
} from '@/helpers/entities/link/list';
import { createEntityIdPatternByValue } from '@/helpers/entities/pattern/form';
import { createDowntimePbehavior } from '@/helpers/entities/pbehavior/form';

import { authMixin } from '@/mixins/auth';
import { widgetActionsPanelAlarmPauseMixin } from '@/mixins/widget/actions-panel/alarm-pause';

import ActionsPanel from '@/components/common/actions-panel/actions-panel.vue';

export default {
  components: { ActionsPanel },

  mixins: [authMixin, widgetActionsPanelAlarmPauseMixin],

  props: {
    items: {
      type: Array,
      default: () => [],
    },
    widget: {
      type: Object,
      required: true,
    },
    refreshAlarmsList: {
      type: Function,
      default: () => {},
    },
    small: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['clear:items'],

  data() {
    return {
      itemsTimer: null,
      localItems: this.items,
      pendingByActionsTypes: {},
    };
  },

  computed: {
    quickMassActions() {
      return this.widget.parameters.quickMassActions ?? [];
    },

    openedAlarms() {
      return this.localItems.filter(
        item => !isAlarmStatusCancelled(item) && !isAlarmStatusClosed(item),
      );
    },

    alarmsForActions() {
      return this.localItems.filter(
        item => (this.widget.parameters.isActionsAllowWithOkState && isAlarmStateOk(item)) ||
          (!isAlarmStatusCancelled(item) && !isAlarmStatusClosed(item)),
      );
    },

    alarmsWithFastPbehavior() {
      return this.alarmsForActions.filter(item => item.pbh_origin_icon);
    },

    alarmsWithoutFastPbehavior() {
      return this.alarmsForActions.filter(item => !item.pbh_origin_icon);
    },

    alarmsWithOwnPause() {
      return this.alarmsForActions.filter(item => isResumeAvailableForEntity(item.entity));
    },

    alarmsWithoutPause() {
      return this.alarmsForActions.filter(item => !isEntityPaused(item.entity));
    },

    alarmsForActionsWithoutUnknown() {
      return this.alarmsForActions.filter(item => !isAlarmStatusUnknown(item));
    },

    openedAndUnResolvedAlarms() {
      return this.openedAlarms.filter(alarm => !isResolvedAlarm(alarm));
    },

    cancelledAndUnResolvedAlarms() {
      return this.localItems.filter(
        alarm => isAlarmStatusCancelled(alarm) && !isResolvedAlarm(alarm),
      );
    },

    alarmsWithAssignedDeclareTicketRules() {
      return this.alarmsForActions.filter(item => item.assigned_declare_ticket_rules?.length);
    },

    alarmsWithTickets() {
      return this.alarmsForActions.filter(item => item.v?.ticket);
    },

    alarmsWithoutTickets() {
      return difference(this.alarmsForActions, this.alarmsWithTickets);
    },

    alarmsWithAssociatedTickets() {
      return this.alarmsForActions.filter(item => item.v?.tickets?.length > 0);
    },

    alarmsWithAck() {
      return this.alarmsForActions.filter(item => item.v?.ack);
    },

    alarmsWithoutAck() {
      return difference(this.alarmsForActions, this.alarmsWithAck);
    },

    hasOpenedAlarms() {
      return !!this.alarmsForActions.length;
    },

    hasAlarmsWithoutAck() {
      return !!this.alarmsWithoutAck.length;
    },

    hasAlarmsWithAck() {
      return !!this.alarmsWithAck.length;
    },

    hasAlarmsWithoutUnknownStatus() {
      return !!this.alarmsForActionsWithoutUnknown.length;
    },

    hasCancelledAndUnResolvedAlarms() {
      return !!this.cancelledAndUnResolvedAlarms.length;
    },

    hasAlarmsWithoutTickets() {
      return !!this.alarmsWithoutTickets.length;
    },

    hasAlarmsWithAssociatedTickets() {
      return !!this.alarmsWithAssociatedTickets.length;
    },

    hasMetaAlarm() {
      return this.alarmsForActions.some(item => item.is_meta_alarm);
    },

    pauseActions() {
      const actions = [];

      if (this.isPauseActionMode) {
        if (this.alarmsWithOwnPause.length) {
          actions.push(this.getEntityResumeAction(() =>
            this.resumeAlarmEntities(mapAlarmsEntities(this.alarmsWithOwnPause), { bulk: true })));
        }

        if (this.alarmsWithoutPause.length) {
          const pauseTypes = this.widget.parameters.pause?.types ?? [];
          actions.push(this.getEntityPauseAction(
            pauseTypes,
            pauseType => this.pauseAlarmEntities(
              mapAlarmsEntities(this.alarmsWithoutPause),
              pauseType,
              { bulk: true },
            ),
          ));
        }

        return actions;
      }

      if (this.alarmsWithFastPbehavior.length) {
        actions.push({
          type: ALARM_LIST_ACTIONS_TYPES.fastPbehaviorRemove,
          title: this.$t('alarm.actions.titles.fastPbehaviorRemove'),
          method: this.fastRemovePbehavior,
        });
      }

      if (this.alarmsWithoutFastPbehavior.length) {
        const fastPbehaviorsParameters = this.widget.parameters.fast_pbehaviors ?? [];
        actions.push(this.getFastPbehaviorAction(
          fastPbehaviorsParameters,
          pbehaviorParameters => this.fastAddPbehavior(pbehaviorParameters),
        ));
      }

      return actions;
    },

    actions() {
      const actions = [];

      if (this.hasOpenedAlarms) {
        actions.push(
          ...this.pauseActions,
          {
            type: ALARM_LIST_ACTIONS_TYPES.pbehaviorAdd,
            title: this.$t('alarm.actions.titles.pbehavior'),
            method: this.showAddPbehaviorModal,
          },
        );
      }

      if (this.hasAlarmsWithoutAck || this.widget.parameters.isMultiAckEnabled) {
        actions.push(
          {
            type: ALARM_LIST_ACTIONS_TYPES.ack,
            title: this.$t('alarm.actions.titles.ack'),
            method: this.showAckModal,
          },
          {
            type: ALARM_LIST_ACTIONS_TYPES.fastAck,
            title: this.$t('alarm.actions.titles.fastAck'),
            method: this.createMassFastAckEvent,
          },
        );
      }

      if (this.hasAlarmsWithAck) {
        actions.push(
          {
            type: ALARM_LIST_ACTIONS_TYPES.ackRemove,
            title: this.$t('alarm.actions.titles.ackRemove'),
            method: this.showAckRemoveModal,
          },
          {
            type: ALARM_LIST_ACTIONS_TYPES.changeState,
            title: this.$t('alarm.actions.titles.changeState'),
            method: this.showCreateChangeStateEventModal,
          },
        );
      }

      if (this.hasAlarmsWithoutUnknownStatus) {
        actions.push(
          {
            type: ALARM_LIST_ACTIONS_TYPES.cancel,
            icon: '$vuetify.icons.list_delete',
            title: this.$t('alarm.actions.titles.cancel'),
            method: this.showCancelEventModal,
          },
          {
            type: ALARM_LIST_ACTIONS_TYPES.fastCancel,
            icon: 'delete',
            title: this.$t('alarm.actions.titles.fastCancel'),
            method: this.createFastCancelEvent,
          },
        );
      }

      if (this.hasCancelledAndUnResolvedAlarms) {
        actions.push({
          type: ALARM_LIST_ACTIONS_TYPES.unCancel,
          title: this.$t('alarm.actions.titles.unCancel'),
          method: this.showUnCancelEventModal,
        });
      }

      if (this.hasOpenedAlarms) {
        actions.push({
          type: ALARM_LIST_ACTIONS_TYPES.comment,
          title: this.$t('alarm.actions.titles.comment'),
          method: this.showCreateCommentEventModal,
        });
      }

      if (this.hasAlarmsWithAssociatedTickets) {
        actions.push({
          type: ALARM_LIST_ACTIONS_TYPES.removeAssociatedTicket,
          iconSize: 20,
          title: this.$t('alarm.actions.titles.removeAssociatedTicket'),
          method: this.showRemoveAssociatedTicketModal,
        });
      }

      if (this.hasAlarmsWithoutTickets || this.widget.parameters.isMultiDeclareTicketEnabled) {
        if (this.alarmsWithAssignedDeclareTicketRules.length) {
          actions.push({
            type: ALARM_LIST_ACTIONS_TYPES.declareTicket,
            title: this.$t('alarm.actions.titles.declareTicket'),
            method: this.showCreateDeclareTicketModal,
          });
        }

        actions.push({
          type: ALARM_LIST_ACTIONS_TYPES.associateTicket,
          title: this.$t('alarm.actions.titles.associateTicket'),
          method: this.showCreateAssociateTicketModal,
        });
      }

      if (this.openedAndUnResolvedAlarms.length) {
        actions.push({
          type: ALARM_LIST_ACTIONS_TYPES.snooze,
          title: this.$t('alarm.actions.titles.snooze'),
          method: this.showSnoozeModal,
        });
      }

      if (!this.hasMetaAlarm) {
        actions.push({
          type: ALARM_LIST_ACTIONS_TYPES.linkToMetaAlarm,
          title: this.$t('alarm.actions.titles.linkToMetaAlarm'),
          method: this.showLinkToMetaAlarmModal,
        });
      }

      const featuresActionsKey =
        FEATURES_SERVICE_KEYS.components.alarmListMassActionsPanel.computed.actions;

      if (featuresService.has(featuresActionsKey)) {
        const featuresActions = featuresService.call(
          featuresActionsKey,
          this,
          [],
        );

        if (featuresActions?.length) {
          actions.push(...featuresActions);
        }
      }

      return actions;
    },

    linksActions() {
      const permission = BUSINESS_USER_PERMISSIONS_ACTIONS_MAP.alarmsList[ALARM_LIST_ACTIONS_TYPES.links];

      if (!this.checkAccess(permission)) {
        return [];
      }

      return harmonizeAlarmsLinks(this.alarmsForActions).map((link) => {
        const type = getLinkRuleLinkActionType(link);

        return {
          type,
          link: true,
          icon: link.icon_name,
          title: this.$t('alarm.followLink', { title: link.label }),
          method: () => this.linkAction(link, type),
        };
      });
    },

    filteredActions() {
      return this.actions.filter(this.actionsAccessFilterHandler);
    },

    preparedActions() {
      return sortActionsByQuickActions(
        [...this.filteredActions, ...this.linksActions],
        this.quickMassActions,
      ).map(action => ({
        ...action,
        icon: action.icon ?? getAlarmActionIcon(action.type),
        loading: this.isActionTypeInPending(action.type),
      }));
    },

    inlineCount() {
      if (!this.widget.parameters.hideMassActions) {
        return this.preparedActions.length;
      }

      if (this.preparedActions.length && !this.quickMassActions.length) {
        return 1;
      }

      return getActionsInlineCount(this.preparedActions, this.quickMassActions);
    },
  },

  watch: {
    items(items) {
      if (!items?.length) {
        this.itemsTimer = setTimeout(this.setLocalItems, VUETIFY_ANIMATION_DELAY);
      } else {
        clearTimeout(this.itemsTimer);
        this.setLocalItems();
      }
    },
  },

  beforeUnmount() {
    clearTimeout(this.itemsTimer);
  },

  methods: {
    ...featuresService.get(FEATURES_SERVICE_KEYS.components.alarmListMassActionsPanel.methods, {}),
    ...mapActions(useAlarmStore, [
      'bulkCreateAlarmAckEvent',
      'bulkCreateAlarmAckremoveEvent',
      'bulkCreateAlarmSnoozeEvent',
      'bulkCreateAlarmAssocticketEvent',
      'bulkCreateAlarmTicketremoveEvent',
      'bulkCreateAlarmCommentEvent',
      'bulkCreateAlarmCancelEvent',
      'bulkCreateAlarmUnCancelEvent',
      'bulkCreateAlarmChangestateEvent',
      'fetchAlarmLinksWithoutStore',
    ]),
    ...mapActions(useDeclareTicketRuleStore, {
      fetchAssignedTicketsWithoutStore: 'fetchAssignedTicketsWithoutStore',
    }),
    ...mapActions(useMetaAlarmStore, {
      createMetaAlarm: 'create',
      addAlarmsIntoMetaAlarm: 'addAlarms',
    }),
    ...mapActions(usePbehaviorStore, [
      'bulkCreateEntityPbehaviors',
      'bulkRemoveEntityPbehaviors',
    ]),
    ...mapActions(usePbehaviorTypeStore, {
      fetchDefaultPbehaviorTypes: 'fetchListWithoutStore',
    }),
    ...mapActions(usePbehaviorReasonStore, {
      fetchPbehaviorReasonsListWithoutStore: 'fetchListWithoutStore',
    }),

    setLocalItems() {
      this.localItems = this.items;
    },

    clearItems() {
      this.$emit('clear:items');
    },

    afterSubmit() {
      this.clearItems();
      this.refreshAlarmsList();
    },

    isActionTypeInPending(type) {
      return (
        !!this.pendingByActionsTypes[type] ||
        !!this.pendingByActionsTypes[ALARM_LIST_TOGGLE_ACTIONS_TYPES_MAP[type]]
      );
    },

    setActionPending(type, value) {
      this.pendingByActionsTypes = {
        ...this.pendingByActionsTypes,
        [type]: value,
      };
    },

    showEventModal(name, alarms, action, config = {}) {
      const { prepareEvent = event => event, ...modalConfig } = config;

      this.$modals.show({
        name,
        config: {
          items: alarms,
          ...modalConfig,
          action: async (event) => {
            const preparedEvent = prepareEvent(event);

            await action({ data: alarms.map(alarm => ({ ...preparedEvent, _id: alarm._id })) });
            await this.afterSubmit();
          },
        },
      });
    },

    showSnoozeModal() {
      this.showEventModal(
        MODALS.createSnoozeEvent,
        this.openedAndUnResolvedAlarms,
        this.bulkCreateAlarmSnoozeEvent,
        {
          isNoteRequired: this.widget.parameters.isSnoozeNoteRequired,
        },
      );
    },

    showAckModal() {
      const alarms = this.widget.parameters.isMultiAckEnabled
        ? this.alarmsForActions
        : this.alarmsWithoutAck;

      this.$modals.show({
        name: MODALS.createAckEvent,
        config: {
          items: alarms,
          isNoteRequired: this.widget.parameters.isAckNoteRequired,
          action: async (event, { needDeclareTicket, needAssociateTicket } = {}) => {
            await this.bulkCreateAlarmAckEvent({
              data: alarms.map(alarm => ({ ...event, _id: alarm._id })),
            });
            await this.afterSubmit();

            if (needAssociateTicket) {
              this.showAssociateTicketModal(alarms, true);
            } else if (needDeclareTicket) {
              await this.showDeclareTicketModal(this.alarmsWithAssignedDeclareTicketRules);
            }
          },
        },
      });
    },

    async createMassFastAckEvent() {
      this.setActionPending(ALARM_LIST_ACTIONS_TYPES.fastAck, true);

      try {
        const alarms = this.widget.parameters.isMultiAckEnabled
          ? this.alarmsForActions
          : this.alarmsWithoutAck;
        const { fastAckOutput } = this.widget.parameters;
        const event = { comment: fastAckOutput?.enabled ? fastAckOutput.value : '' };

        await this.bulkCreateAlarmAckEvent({
          data: alarms.map(alarm => ({ ...event, _id: alarm._id })),
        });
        await this.afterSubmit();
      } catch (err) {
        console.error(err);
      } finally {
        this.setActionPending(ALARM_LIST_ACTIONS_TYPES.fastAck, false);
      }
    },

    showCreateAssociateTicketModal() {
      this.showAssociateTicketModal(
        this.widget.parameters.isMultiDeclareTicketEnabled
          ? this.alarmsForActions
          : this.alarmsWithoutTickets,
      );
    },

    showAssociateTicketModal(alarms, ignoreAck = false) {
      this.$modals.show({
        name: MODALS.createAssociateTicketEvent,
        config: {
          items: alarms,
          ignoreAck,
          action: async (associateEvent) => {
            if (!ignoreAck) {
              const { fastAckOutput } = this.widget.parameters;
              const alarmsWithoutAck = alarms.filter(alarm => !alarm.v?.ack);

              await this.bulkCreateAlarmAckEvent({
                data: alarmsWithoutAck.map(alarm => ({
                  comment: fastAckOutput?.enabled ? fastAckOutput.value : '',
                  _id: alarm._id,
                })),
              });
            }

            await this.bulkCreateAlarmAssocticketEvent({
              data: alarms.map(alarm => ({ ...associateEvent, _id: alarm._id })),
            });

            await this.afterSubmit();
          },
        },
      });
    },

    showRemoveAssociatedTicketModal() {
      this.showEventModal(
        MODALS.removeAssociatedTicketEvent,
        this.alarmsWithAssociatedTickets,
        this.bulkCreateAlarmTicketremoveEvent,
      );
    },

    showCreateDeclareTicketModal() {
      return this.showDeclareTicketModal(
        this.widget.parameters.isMultiDeclareTicketEnabled
          ? this.alarmsForActions
          : this.alarmsWithoutTickets,
      );
    },

    async showDeclareTicketModal(alarms) {
      this.setActionPending(ALARM_LIST_ACTIONS_TYPES.declareTicket, true);

      try {
        const response = await this.fetchAssignedTicketsWithoutStore({
          params: { alarms: mapIds(alarms) },
        });

        this.$modals.show({
          name: MODALS.createDeclareTicketEvent,
          config: {
            items: alarms,
            alarmsByTickets: response.by_rules,
            ticketsByAlarms: response.by_alarms,
            action: (events, singleMode) => {
              this.$modals.show({
                name: MODALS.executeDeclareTickets,
                config: {
                  singleMode,
                  executions: events,
                  tickets: events.map(({ _id: id }) => ({
                    _id: id,
                    name: response.by_rules[id].name,
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

    showCreateCommentEventModal() {
      const commentTemplatesById = keyBy(this.widget.comment_templates, '_id');

      this.showEventModal(MODALS.createCommentEvent, this.alarmsForActions, this.bulkCreateAlarmCommentEvent, {
        templates: (this.widget.parameters.comment_templates ?? [])
          .map(id => commentTemplatesById[id])
          .filter(Boolean),
      });
    },

    showAckRemoveModal() {
      this.showEventModal(MODALS.createEvent, this.alarmsForActions, this.bulkCreateAlarmAckremoveEvent, {
        title: this.$t('modals.createAckRemove.title'),
      });
    },

    showCancelEventModal() {
      this.showEventModal(
        MODALS.createEvent,
        this.alarmsForActionsWithoutUnknown,
        this.bulkCreateAlarmCancelEvent,
        {
          title: this.$t('modals.createCancelEvent.title'),
        },
      );
    },

    showUnCancelEventModal() {
      this.showEventModal(
        MODALS.createEvent,
        this.cancelledAndUnResolvedAlarms,
        this.bulkCreateAlarmUnCancelEvent,
        {
          title: this.$t('modals.createUnCancel.title'),
          isCommentRequired: this.widget.parameters.isUncancelAlarmsCommentRequired,
        },
      );
    },

    showCreateChangeStateEventModal() {
      this.showEventModal(
        MODALS.createChangeStateEvent,
        this.alarmsWithAck,
        this.bulkCreateAlarmChangestateEvent,
        { prepareEvent: formToChangeStateEvent },
      );
    },

    async createFastCancelEvent() {
      this.setActionPending(ALARM_LIST_ACTIONS_TYPES.fastCancel, true);

      try {
        const { fastCancelOutput } = this.widget.parameters;
        const event = { comment: fastCancelOutput?.enabled ? fastCancelOutput.value : '' };

        await this.bulkCreateAlarmCancelEvent({
          data: this.alarmsForActionsWithoutUnknown.map(alarm => ({ ...event, _id: alarm._id })),
        });
        await this.afterSubmit();
      } catch (err) {
        console.error(err);
      } finally {
        this.setActionPending(ALARM_LIST_ACTIONS_TYPES.fastCancel, false);
      }
    },

    showAddPbehaviorModal() {
      const entities = mapAlarmsEntities(this.alarmsForActions);
      const entityIds = entities.map(({ _id }) => _id);

      this.$modals.show({
        name: MODALS.pbehaviorPlanning,
        config: {
          entityPattern: createEntityIdPatternByValue(
            entityIds.length === 1 ? entityIds[0] : entityIds,
          ),
          entities,
          afterSubmit: async () => {
            this.$popups.success({ text: this.$t('modals.createPbehavior.success.create') });

            await promisedWait(ALARM_LIST_FAST_PBEHAVIOR_TIMEOUT);

            return this.afterSubmit();
          },
        },
      });
    },

    async fastAddPbehavior(pbehaviorParameters = {}) {
      this.setActionPending(ALARM_LIST_ACTIONS_TYPES.fastPbehaviorAdd, true);

      try {
        const { name_prefix: namePrefix = '' } = pbehaviorParameters ?? {};
        let { type, reason } = pbehaviorParameters ?? {};

        if (!type) {
          const { data: defaultPbehaviorTypes = [] } = await this.fetchDefaultPbehaviorTypes();
          const pauseType = defaultPbehaviorTypes.find(
            ({ type: pbehaviorType }) => pbehaviorType === PBEHAVIOR_TYPE_TYPES.pause,
          );

          type = pauseType?._id;
        }

        if (!reason) {
          const { data: reasons = [] } = await this.fetchPbehaviorReasonsListWithoutStore();

          reason = reasons?.[0]?._id;
        }

        await this.bulkCreateEntityPbehaviors({
          data: mapAlarmsEntities(this.alarmsWithoutFastPbehavior).map(entity => createDowntimePbehavior({
              entity,
              prefix: namePrefix,
              type,
              reason,
              origin: PBEHAVIOR_ORIGINS.alarmList,
            }),
          ),
        });

        await promisedWait(ALARM_LIST_FAST_PBEHAVIOR_TIMEOUT);
        await this.afterSubmit();
      } catch (err) {
        console.error(err);
      } finally {
        this.setActionPending(ALARM_LIST_ACTIONS_TYPES.fastPbehaviorAdd, false);
      }
    },

    async fastRemovePbehavior() {
      this.setActionPending(ALARM_LIST_ACTIONS_TYPES.fastPbehaviorRemove, true);

      try {
        await this.bulkRemoveEntityPbehaviors({
          data: mapAlarmsEntities(this.alarmsWithFastPbehavior).map(({ _id: id }) => ({
            origin: PBEHAVIOR_ORIGINS.alarmList,
            entity: id,
          })),
        });

        await promisedWait(ALARM_LIST_FAST_PBEHAVIOR_TIMEOUT);
        await this.afterSubmit();
      } catch (err) {
        console.error(err);
      } finally {
        this.setActionPending(ALARM_LIST_ACTIONS_TYPES.fastPbehaviorRemove, false);
      }
    },

    showLinkToMetaAlarmModal() {
      this.$modals.show({
        name: MODALS.linkToMetaAlarm,
        config: {
          items: this.alarmsForActions,
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

    async linkAction(link, type) {
      try {
        this.setActionPending(type, true);

        const links = await this.fetchAlarmLinksWithoutStore({
          id: link.rule_id,
          params: { ids: mapIds(this.alarmsForActions) },
        });

        const summaryLink = links.find(
          ({ icon_name: iconName, label = '' } = {}) => iconName === link.icon_name && label.includes(link.label),
        );

        if (!summaryLink) {
          return;
        }

        if (link.action === LINK_RULE_ACTIONS.copy) {
          await navigator.clipboard?.writeText(summaryLink.url);

          return;
        }

        window.open(summaryLink.url, '_blank');
      } catch (err) {
        console.error(err);
      } finally {
        this.setActionPending(type, false);
      }
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
</script>
