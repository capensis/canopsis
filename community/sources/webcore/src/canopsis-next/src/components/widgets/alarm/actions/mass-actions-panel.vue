<template lang="pug">
  shared-mass-actions-panel(:actions="preparedActions")
</template>

<script>
import { difference, find, pick } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import {
  MODALS,
  EVENT_ENTITY_TYPES,
  ALARM_LIST_ACTIONS_TYPES,
  BUSINESS_USER_PERMISSIONS_ACTIONS_MAP,
} from '@/constants';

import featuresService from '@/services/features';

import { mapIds } from '@/helpers/entities';
import { getEntityEventIcon } from '@/helpers/icon';
import { createEntityIdPatternByValue } from '@/helpers/pattern';
import { harmonizeAlarmsLinks, getLinkRuleLinkActionType } from '@/helpers/links';

import { widgetActionsPanelAlarmMixin } from '@/mixins/widget/actions-panel/alarm';
import { entitiesDeclareTicketRuleMixin } from '@/mixins/entities/declare-ticket-rule';
import { entitiesAlarmLinksMixin } from '@/mixins/entities/alarm/links';

import SharedMassActionsPanel from '@/components/common/actions-panel/mass-actions-panel.vue';

const { mapGetters: entitiesMapGetters } = createNamespacedHelpers('entities');

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
    entitiesAlarmLinksMixin,
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
      pendingByActionsTypes: {},
    };
  },
  computed: {
    ...entitiesMapGetters({
      getEntitiesList: 'getList',
    }),

    actions() {
      const actions = [
        {
          type: ALARM_LIST_ACTIONS_TYPES.pbehaviorAdd,
          icon: getEntityEventIcon(EVENT_ENTITY_TYPES.pbehaviorAdd),
          title: this.$t('alarm.actions.titles.pbehavior'),
          method: this.showAddPbehaviorModal,
        },
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
        {
          type: ALARM_LIST_ACTIONS_TYPES.ackRemove,
          icon: getEntityEventIcon(EVENT_ENTITY_TYPES.ackRemove),
          title: this.$t('alarm.actions.titles.ackRemove'),
          method: this.showAckRemoveModal,
        },
        {
          type: ALARM_LIST_ACTIONS_TYPES.cancel,
          icon: getEntityEventIcon(EVENT_ENTITY_TYPES.delete),
          title: this.$t('alarm.actions.titles.cancel'),
          method: this.showCancelEventModal,
        },
        {
          type: ALARM_LIST_ACTIONS_TYPES.comment,
          icon: getEntityEventIcon(EVENT_ENTITY_TYPES.comment),
          title: this.$t('alarm.actions.titles.comment'),
          method: this.showCreateCommentEventModal,
        },
      ];

      if (this.hasAlarmsWithoutTickets || this.widget.parameters.isMultiDeclareTicketEnabled) {
        if (this.alarmsWithAssignedDeclareTicketRules.length) {
          actions.push({
            type: ALARM_LIST_ACTIONS_TYPES.declareTicket,
            icon: getEntityEventIcon(EVENT_ENTITY_TYPES.declareTicket),
            title: this.$t('alarm.actions.titles.declareTicket'),
            loading: this.ticketsForAlarmsPending,
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

      actions.push(
        {
          type: ALARM_LIST_ACTIONS_TYPES.snooze,
          icon: getEntityEventIcon(EVENT_ENTITY_TYPES.snooze),
          title: this.$t('alarm.actions.titles.snooze'),
          method: this.showSnoozeModal,
        },
      );

      if (!this.hasMetaAlarm) {
        actions.push(
          {
            type: ALARM_LIST_ACTIONS_TYPES.groupRequest,
            icon: getEntityEventIcon(EVENT_ENTITY_TYPES.groupRequest),
            title: this.$t('alarm.actions.titles.groupRequest'),
            method: this.showCreateGroupRequestEventModal,
          },
          {
            type: ALARM_LIST_ACTIONS_TYPES.createManualMetaAlarm,
            icon: getEntityEventIcon(EVENT_ENTITY_TYPES.createManualMetaAlarm),
            title: this.$t('alarm.actions.titles.createManualMetaAlarm'),
            method: this.showCreateManualMetaAlarmModal,
          },
        );
      }

      /**
       * If we will have actions for resolved alarms in the features we should move this condition to
       * the every features repositories
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

      return harmonizeAlarmsLinks(this.items).map((link) => {
        const type = getLinkRuleLinkActionType(link);

        return {
          type,
          icon: link.icon_name,
          title: this.$t('alarm.followLink', { title: link.label }),
          loading: this.pendingByActionsTypes[type],
          method: () => this.openLink(link, type),
        };
      });
    },

    preparedActions() {
      return [
        ...this.filteredActions,
        ...this.linksActions,
      ];
    },

    alarmsWithAssignedDeclareTicketRules() {
      return this.items.filter(item => item.assigned_declare_ticket_rules?.length);
    },

    alarmsWithTickets() {
      return this.items.filter(item => item.v?.tickets?.length);
    },

    alarmsWithoutTickets() {
      return difference(this.items, this.alarmsWithTickets);
    },

    hasAlarmsWithoutTickets() {
      return !!this.alarmsWithoutTickets.length;
    },

    hasMetaAlarm() {
      return this.items.some(item => item.is_meta_alarm);
    },

    modalConfig() {
      return {
        items: this.items,
        afterSubmit: this.afterSubmit,
      };
    },
  },

  methods: {
    ...featuresService.get('components.alarmListMassActionsPanel.methods', {}),

    clearItems() {
      this.$emit('clear:items');
    },

    afterSubmit() {
      this.clearItems();

      return this.refreshAlarmsList();
    },

    showAddPbehaviorModal() {
      this.$modals.show({
        name: MODALS.pbehaviorPlanning,
        config: {
          entityPattern: createEntityIdPatternByValue(this.items.map(item => item.entity._id)),
          afterSubmit: this.clearItems,
        },
      });
    },

    showCreateAssociateTicketModal() {
      this.showAssociateTicketModalByAlarms(
        this.widget.parameters.isMultiDeclareTicketEnabled
          ? this.items
          : this.alarmsWithoutTickets,
      );
    },

    showCreateDeclareTicketModal() {
      this.showDeclareTicketModalByAlarms(
        this.widget.parameters.isMultiDeclareTicketEnabled
          ? this.items
          : this.alarmsWithoutTickets,
      );
    },

    showAckModal() {
      this.showAckModalByAlarms(this.items);
    },

    showCreateGroupRequestEventModal() {
      this.$modals.show({
        name: MODALS.createEvent,
        config: {
          ...this.modalConfig,

          title: this.$t('modals.createGroupRequestEvent.title'),
          eventType: EVENT_ENTITY_TYPES.groupRequest,
        },
      });
    },

    showCreateManualMetaAlarmModal() {
      this.$modals.show({
        name: MODALS.createManualMetaAlarm,
        config: {
          ...this.modalConfig,

          title: this.$t('modals.createManualMetaAlarm.title'),
        },
      });
    },

    async createMassFastAckEvent() {
      let eventData = {};

      if (this.widget.parameters.fastAckOutput && this.widget.parameters.fastAckOutput.enabled) {
        eventData = { output: this.widget.parameters.fastAckOutput.value };
      }

      await this.createEvent(EVENT_ENTITY_TYPES.ack, this.items, eventData);

      return this.afterSubmit();
    },

    showCreateCommentEventModal() {
      this.$modals.show({
        name: MODALS.createCommentEvent,
        config: {
          ...this.modalConfig,
          action: data => this.createEvent(EVENT_ENTITY_TYPES.comment, this.items, data),
        },
      });
    },

    async openLink(link, type) {
      try {
        this.$set(this.pendingByActionsTypes, type, true);

        const links = await this.fetchAlarmLinkWithoutStore({
          id: link.rule_id,
          params: { ids: mapIds(this.items) },
        });

        const summaryLink = find(links, pick(link, ['icon_name, label']));

        if (!summaryLink) {
          return;
        }

        window.open(summaryLink.url, '_blank');
      } catch (err) {
        console.error(err);
      } finally {
        this.$set(this.pendingByActionsTypes, type, false);
      }
    },
  },
};
</script>
