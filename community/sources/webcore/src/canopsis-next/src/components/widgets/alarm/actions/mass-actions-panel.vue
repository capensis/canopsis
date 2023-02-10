<template lang="pug">
  shared-mass-actions-panel(:actions="filteredActions")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import {
  MODALS,
  EVENT_ENTITY_TYPES,
  EVENT_ENTITY_STYLE,
  ALARM_LIST_ACTIONS_TYPES,
} from '@/constants';

import featuresService from '@/services/features';

import { createEntityIdPatternByValue } from '@/helpers/pattern';

import { widgetActionsPanelAlarmMixin } from '@/mixins/widget/actions-panel/alarm';
import { entitiesDeclareTicketRuleMixin } from '@/mixins/entities/declare-ticket-rule';

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
  mixins: [widgetActionsPanelAlarmMixin, entitiesDeclareTicketRuleMixin],
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
  computed: {
    ...entitiesMapGetters({
      getEntitiesList: 'getList',
    }),

    actions() {
      const actions = [
        {
          type: ALARM_LIST_ACTIONS_TYPES.pbehaviorAdd,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.pbehaviorAdd].icon,
          title: this.$t('alarm.actions.titles.pbehavior'),
          method: this.showAddPbehaviorModal,
        },
        {
          type: ALARM_LIST_ACTIONS_TYPES.ack,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.ack].icon,
          title: this.$t('alarm.actions.titles.ack'),
          method: this.showAckModal,
        },
        {
          type: ALARM_LIST_ACTIONS_TYPES.fastAck,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.fastAck].icon,
          title: this.$t('alarm.actions.titles.fastAck'),
          method: this.createMassFastAckEvent,
        },
        {
          type: ALARM_LIST_ACTIONS_TYPES.ackRemove,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.ackRemove].icon,
          title: this.$t('alarm.actions.titles.ackRemove'),
          method: this.showAckRemoveModal,
        },
        {
          type: ALARM_LIST_ACTIONS_TYPES.cancel,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.delete].icon,
          title: this.$t('alarm.actions.titles.cancel'),
          method: this.showCancelEventModal,
        },
        {
          type: ALARM_LIST_ACTIONS_TYPES.declareTicket,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.declareTicket].icon,
          title: this.$t('alarm.actions.titles.declareTicket'),
          loading: this.ticketsForAlarmsPending,
          method: this.showCreateDeclareTicketModal,
        },
        {
          type: ALARM_LIST_ACTIONS_TYPES.associateTicket,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.assocTicket].icon,
          title: this.$t('alarm.actions.titles.associateTicket'),
          method: this.showCreateAssociateTicketModal,
        },
        {
          type: ALARM_LIST_ACTIONS_TYPES.snooze,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.snooze].icon,
          title: this.$t('alarm.actions.titles.snooze'),
          method: this.showSnoozeModal,
        },
      ];

      if (!this.hasMetaAlarm) {
        actions.push(
          {
            type: ALARM_LIST_ACTIONS_TYPES.groupRequest,
            icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.groupRequest].icon,
            title: this.$t('alarm.actions.titles.groupRequest'),
            method: this.showCreateGroupRequestEventModal,
          },
          {
            type: ALARM_LIST_ACTIONS_TYPES.manualMetaAlarmGroup,
            icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.manualMetaAlarmGroup].icon,
            title: this.$t('alarm.actions.titles.manualMetaAlarmGroup'),
            method: this.showCreateManualMetaAlarmModal,
          },
        );
      }

      if (featuresService.has('components.alarmListMassActionsPanel.computed.actions')) {
        actions.push(...featuresService.call('components.alarmListMassActionsPanel.computed.actions', this, []));
      }

      return actions;
    },

    filteredActions() {
      return this.actions.filter(this.actionsAccessFilterHandler);
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
      this.showAssociateTicketModalByAlarms(this.items);
    },

    showCreateDeclareTicketModal() {
      this.showDeclareTicketModalByAlarms(this.items);
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
  },
};
</script>
