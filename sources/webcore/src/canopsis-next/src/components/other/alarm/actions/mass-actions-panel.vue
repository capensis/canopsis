<template lang="pug">
  div
    shared-mass-actions-panel(:actions="filteredActions")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MODALS, ENTITIES_TYPES, EVENT_ENTITY_TYPES, EVENT_ENTITY_STYLE, WIDGETS_ACTIONS_TYPES } from '@/constants';

import authMixin from '@/mixins/auth';
import widgetActionsPanelAlarmMixin from '@/mixins/widget/actions-panel/alarm';

import SharedMassActionsPanel from '@/components/other/shared/actions-panel/mass-actions-panel.vue';

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
  mixins: [authMixin, widgetActionsPanelAlarmMixin],
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
  data() {
    const { alarmsList: alarmsListActionsTypes } = WIDGETS_ACTIONS_TYPES;

    return {
      actions: [
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
          method: this.showActionModal(MODALS.createCancelEvent),
        },
        {
          type: alarmsListActionsTypes.snooze,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.snooze].icon,
          title: this.$t('alarmList.actions.titles.snooze'),
          method: this.showActionModal(MODALS.createSnoozeEvent),
        },
      ],
    };
  },
  computed: {
    ...entitiesMapGetters({
      getEntitiesList: 'getList',
    }),

    filteredActions() {
      return this.actions.filter(this.actionsAccessFilterHandler);
    },

    items() {
      return this.getEntitiesList(ENTITIES_TYPES.alarm, this.itemsIds);
    },

    modalConfig() {
      return {
        itemsType: ENTITIES_TYPES.alarm,
        itemsIds: this.itemsIds,
      };
    },
  },
};
</script>
