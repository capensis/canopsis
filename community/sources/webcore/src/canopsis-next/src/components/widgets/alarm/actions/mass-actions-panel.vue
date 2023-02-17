<template lang="pug">
  shared-mass-actions-panel(:actions="filteredActions")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import {
  MODALS,
  EVENT_ENTITY_TYPES,
  ALARM_LIST_ACTIONS_TYPES,
} from '@/constants';

import featuresService from '@/services/features';

import { getEntityEventIcon } from '@/helpers/icon';
import { createEntityIdPatternByValue } from '@/helpers/pattern';

import { widgetActionsPanelAlarmMixin } from '@/mixins/widget/actions-panel/alarm';

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
  mixins: [widgetActionsPanelAlarmMixin],
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
          type: ALARM_LIST_ACTIONS_TYPES.associateTicket,
          icon: getEntityEventIcon(EVENT_ENTITY_TYPES.assocTicket),
          title: this.$t('alarm.actions.titles.associateTicket'),
          method: this.showCreateAssociateTicketEventModal,
        },
        {
          type: ALARM_LIST_ACTIONS_TYPES.snooze,
          icon: getEntityEventIcon(EVENT_ENTITY_TYPES.snooze),
          title: this.$t('alarm.actions.titles.snooze'),
          method: this.showSnoozeModal,
        },
      ];

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

    showCreateAssociateTicketEventModal() {
      this.$modals.show({
        name: MODALS.createAssociateTicketEvent,
        config: {
          ...this.modalConfig,

          fastAckOutput: this.widget.parameters.fastAckOutput,
        },
      });
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
