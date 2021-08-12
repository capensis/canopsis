<template lang="pug">
  div
    shared-mass-actions-panel(:actions="filteredActions")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MODALS, ENTITIES_TYPES, EVENT_ENTITY_TYPES, EVENT_ENTITY_STYLE, WIDGETS_ACTIONS_TYPES } from '@/constants';

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
    itemsIds: {
      type: Array,
      default: () => [],
    },
    widget: {
      type: Object,
      required: true,
    },
  },
  computed: {
    ...entitiesMapGetters({
      getEntitiesList: 'getList',
    }),

    actions() {
      const { alarmsList: alarmsListActionsTypes } = WIDGETS_ACTIONS_TYPES;

      const actions = [
        {
          type: alarmsListActionsTypes.pbehaviorAdd,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.pbehaviorAdd].icon,
          title: this.$t('alarmList.actions.titles.pbehavior'),
          method: this.showAddPbehaviorModal,
        },
        {
          type: alarmsListActionsTypes.ack,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.ack].icon,
          title: this.$t('alarmList.actions.titles.ack'),
          method: this.showAckModal,
        },
        {
          type: alarmsListActionsTypes.fastAck,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.fastAck].icon,
          title: this.$t('alarmList.actions.titles.fastAck'),
          method: this.createMassFastAckEvent,
        },
        {
          type: alarmsListActionsTypes.ackRemove,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.ackRemove].icon,
          title: this.$t('alarmList.actions.titles.ackRemove'),
          method: this.showAckRemoveModal,
        },
        {
          type: alarmsListActionsTypes.cancel,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.delete].icon,
          title: this.$t('alarmList.actions.titles.cancel'),
          method: this.showCancelEventModal,
        },
        {
          type: alarmsListActionsTypes.associateTicket,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.assocTicket].icon,
          title: this.$t('alarmList.actions.titles.associateTicket'),
          method: this.showCreateAssociateTicketEventModal,
        },
        {
          type: alarmsListActionsTypes.snooze,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.snooze].icon,
          title: this.$t('alarmList.actions.titles.snooze'),
          method: this.showSnoozeModal,
        },
      ];

      if (!this.hasMetaAlarm) {
        actions.push(
          {
            type: alarmsListActionsTypes.groupRequest,
            icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.groupRequest].icon,
            title: this.$t('alarmList.actions.titles.groupRequest'),
            method: this.showCreateGroupRequestEventModal,
          },
          {
            type: alarmsListActionsTypes.manualMetaAlarmGroup,
            icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.manualMetaAlarmGroup].icon,
            title: this.$t('alarmList.actions.titles.manualMetaAlarmGroup'),
            method: this.showCreateManualMetaAlarmModal,
          },
        );
      }

      return actions;
    },

    filteredActions() {
      return this.actions.filter(this.actionsAccessFilterHandler);
    },

    items() {
      return this.getEntitiesList(ENTITIES_TYPES.alarm, this.itemsIds);
    },

    hasMetaAlarm() {
      return this.items.some(item => item.metaalarm);
    },

    modalConfig() {
      return {
        itemsType: ENTITIES_TYPES.alarm,
        itemsIds: this.itemsIds,
      };
    },
  },

  methods: {
    showAddPbehaviorModal() {
      this.$modals.show({
        name: MODALS.pbehaviorPlanning,
        config: {
          filter: {
            _id: { $in: this.items.map(item => item.entity._id) },
          },
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
  },
};
</script>
