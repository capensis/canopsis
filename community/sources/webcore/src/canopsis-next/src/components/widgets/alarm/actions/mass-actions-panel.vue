<template>
  <shared-mass-actions-panel :actions="preparedActions" />
</template>

<script>
import { difference } from 'lodash';

import { VUETIFY_ANIMATION_DELAY } from '@/config';
import { EVENT_ENTITY_TYPES, ALARM_LIST_ACTIONS_TYPES, BUSINESS_USER_PERMISSIONS_ACTIONS_MAP } from '@/constants';

import featuresService from '@/services/features';

import { getEntityEventIcon } from '@/helpers/entities/entity/icons';
import { harmonizeAlarmsLinks, getLinkRuleLinkActionType } from '@/helpers/entities/link/list';
import {
  isAlarmStateOk,
  isCancelledAlarmStatus,
  isClosedAlarmStatus,
  isResolvedAlarm,
} from '@/helpers/entities/alarm/form';

import { widgetActionsPanelAlarmMixin } from '@/mixins/widget/actions-panel/alarm';
import { entitiesDeclareTicketRuleMixin } from '@/mixins/entities/declare-ticket-rule';

import SharedMassActionsPanel from '@/components/common/actions-panel/mass-actions-panel.vue';

/**
 * Panel regrouping mass actions icons
 *
 * @module alarm
 *
 * @prop {Array} [itemIds] - Items selected for the mass action
 */
export default {
  components: { SharedMassActionsPanel },
  mixins: [
    widgetActionsPanelAlarmMixin,
    entitiesDeclareTicketRuleMixin,
  ],
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
  },
  data() {
    return {
      localItems: this.items,
    };
  },
  computed: {
    openedAlarms() {
      return this.localItems.filter(item => !isCancelledAlarmStatus(item) && !isClosedAlarmStatus(item));
    },

    alarmsForActions() {
      return this.localItems.filter((item) => {
        if (this.widget.parameters.isActionsAllowWithOkState && isAlarmStateOk(item)) {
          return true;
        }

        return !isCancelledAlarmStatus(item) && !isClosedAlarmStatus(item);
      });
    },

    openedAndUnResolvedAlarms() {
      return this.openedAlarms.filter(alarm => !isResolvedAlarm(alarm));
    },

    cancelledAndUnResolvedAlarms() {
      return this.localItems.filter(alarm => isCancelledAlarmStatus(alarm) && !isResolvedAlarm(alarm));
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

    alarmsWithAck() {
      return this.alarmsForActions.filter(item => item.v?.ack);
    },

    alarmsWithoutAck() {
      return difference(this.alarmsForActions, this.alarmsWithAck);
    },

    hasOpenedAlarms() {
      return !!this.alarmsForActions.length;
    },

    hasCancelledAndUnResolvedAlarms() {
      return !!this.cancelledAndUnResolvedAlarms.length;
    },

    hasAlarmsWithoutTickets() {
      return !!this.alarmsWithoutTickets.length;
    },

    hasAlarmsWithAck() {
      return !!this.alarmsWithAck.length;
    },

    hasAlarmsWithoutAck() {
      return !!this.alarmsWithoutAck.length;
    },

    hasMetaAlarm() {
      return this.alarmsForActions.some(item => item.is_meta_alarm);
    },

    isActionsAllowWithOkState() {
      return this.widget.parameters.isActionsAllowWithOkState;
    },

    actions() {
      const unCancelAction = {
        type: ALARM_LIST_ACTIONS_TYPES.unCancel,
        icon: getEntityEventIcon(EVENT_ENTITY_TYPES.uncancel),
        title: this.$t('alarm.actions.titles.unCancel'),
        method: this.showUnCancelEventModal,
      };

      const actions = [];

      if (this.hasOpenedAlarms) {
        actions.push(
          {
            type: ALARM_LIST_ACTIONS_TYPES.pbehaviorAdd,
            icon: getEntityEventIcon(EVENT_ENTITY_TYPES.pbehaviorAdd),
            title: this.$t('alarm.actions.titles.pbehavior'),
            method: this.showAddPbehaviorModal,
          },
        );
      }

      if (this.hasAlarmsWithoutAck || this.widget.parameters.isMultiAckEnabled) {
        actions.push(
          {
            type: ALARM_LIST_ACTIONS_TYPES.ack,
            icon: getEntityEventIcon(EVENT_ENTITY_TYPES.ack),
            title: this.$t('alarm.actions.titles.ack'),
            method: this.showAckModal,
          },
          {
            type: ALARM_LIST_ACTIONS_TYPES.fastAck,
            icon: getEntityEventIcon(EVENT_ENTITY_TYPES.fastAck),
            title: this.$t('alarm.actions.titles.fastAck'),
            method: this.createMassFastAckEvent,
          },
        );
      }

      if (this.hasAlarmsWithAck) {
        actions.push(
          {
            type: ALARM_LIST_ACTIONS_TYPES.ackRemove,
            icon: getEntityEventIcon(EVENT_ENTITY_TYPES.ackRemove),
            title: this.$t('alarm.actions.titles.ackRemove'),
            method: this.showAckRemoveModal,
          },
        );
      }

      if (this.hasOpenedAlarms) {
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
        actions.push(unCancelAction);
      }

      if (this.hasOpenedAlarms) {
        actions.push({
          type: ALARM_LIST_ACTIONS_TYPES.comment,
          icon: getEntityEventIcon(EVENT_ENTITY_TYPES.comment),
          title: this.$t('alarm.actions.titles.comment'),
          method: this.showCreateCommentEventModal,
        });
      }

      if (this.hasAlarmsWithoutTickets || this.widget.parameters.isMultiDeclareTicketEnabled) {
        if (this.alarmsWithAssignedDeclareTicketRules.length) {
          actions.push({
            type: ALARM_LIST_ACTIONS_TYPES.declareTicket,
            icon: getEntityEventIcon(EVENT_ENTITY_TYPES.declareTicket),
            title: this.$t('alarm.actions.titles.declareTicket'),
            method: this.showCreateDeclareTicketModal,
          });
        }

        actions.push({
          type: ALARM_LIST_ACTIONS_TYPES.associateTicket,
          icon: getEntityEventIcon(EVENT_ENTITY_TYPES.assocTicket),
          title: this.$t('alarm.actions.titles.associateTicket'),
          method: this.showCreateAssociateTicketModal,
        });
      }

      if (this.openedAndUnResolvedAlarms.length) {
        actions.push(
          {
            type: ALARM_LIST_ACTIONS_TYPES.snooze,
            icon: getEntityEventIcon(EVENT_ENTITY_TYPES.snooze),
            title: this.$t('alarm.actions.titles.snooze'),
            method: this.showSnoozeModal,
          },
        );
      }

      if (!this.hasMetaAlarm) {
        actions.push(
          {
            type: ALARM_LIST_ACTIONS_TYPES.createManualMetaAlarm,
            icon: getEntityEventIcon(EVENT_ENTITY_TYPES.createManualMetaAlarm),
            title: this.$t('alarm.actions.titles.createManualMetaAlarm'),
            method: this.showCreateManualMetaAlarmModal,
          },
        );
      }

      /**
       * If we have actions for resolved alarms in the features we should move this condition to
       * the every feature's repositories
       */
      if (featuresService.has('components.alarmListMassActionsPanel.computed.actions')) {
        const featuresActions = featuresService.call('components.alarmListMassActionsPanel.computed.actions', this, []);

        if (featuresActions?.length) {
          actions.push(...featuresActions);
        }
      }

      return actions;
    },

    filteredActions() {
      return this.actions.filter(this.actionsAccessFilterHandler);
    },

    linksActions() {
      if (!this.checkAccess(BUSINESS_USER_PERMISSIONS_ACTIONS_MAP.alarmsList[ALARM_LIST_ACTIONS_TYPES.links])) {
        return [];
      }

      return harmonizeAlarmsLinks(this.alarmsForActions).map((link) => {
        const type = getLinkRuleLinkActionType(link);

        return {
          type,
          icon: link.icon_name,
          title: this.$t('alarm.followLink', { title: link.label }),
          method: () => this.linkAction(link, type),
        };
      });
    },

    preparedActions() {
      return [
        ...this.filteredActions,
        ...this.linksActions,
      ].map(action => ({
        ...action,
        loading: this.isActionTypeInPending(action.type),
      }));
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
  methods: {
    ...featuresService.get('components.alarmListMassActionsPanel.methods', {}),

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

    showSnoozeModal() {
      this.showSnoozeModalByAlarms(this.openedAndUnResolvedAlarms);
    },

    showAddPbehaviorModal() {
      this.showAddPbehaviorModalByAlarms(this.alarmsForActions);
    },

    showCreateAssociateTicketModal() {
      this.showAssociateTicketModalByAlarms(
        this.widget.parameters.isMultiDeclareTicketEnabled
          ? this.alarmsForActions
          : this.alarmsWithoutTickets,
      );
    },

    showCreateDeclareTicketModal() {
      this.showDeclareTicketModalByAlarms(
        this.widget.parameters.isMultiDeclareTicketEnabled
          ? this.alarmsForActions
          : this.alarmsWithoutTickets,
      );
    },

    showAckModal() {
      this.showAckModalByAlarms(
        this.widget.parameters.isMultiAckEnabled
          ? this.alarmsForActions
          : this.alarmsWithoutAck,
      );
    },

    showAckRemoveModal() {
      this.showAckRemoveModalByAlarms(this.alarmsForActions);
    },

    showCancelEventModal() {
      this.showCancelModalByAlarms(this.alarmsForActions);
    },

    showUnCancelEventModal() {
      this.showUnCancelModalByAlarms(this.cancelledAndUnResolvedAlarms);
    },

    showCreateManualMetaAlarmModal() {
      this.showCreateManualMetaAlarmModalByAlarms(this.alarmsForActions);
    },

    createMassFastAckEvent() {
      this.createFastAckActionByAlarms(
        this.widget.parameters.isMultiAckEnabled
          ? this.alarmsForActions
          : this.alarmsWithoutAck,
      );
    },

    createFastCancelEvent() {
      this.createFastCancelActionByAlarms(this.alarmsForActions);
    },

    showCreateCommentEventModal() {
      this.showCreateCommentModalByAlarms(this.alarmsForActions);
    },

    linkAction(link, type) {
      this.handleLinkClickActionByAlarms(this.alarmsForActions, link, type);
    },
  },
};
</script>
