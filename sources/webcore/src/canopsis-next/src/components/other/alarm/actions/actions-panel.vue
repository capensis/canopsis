<template lang="pug">
  shared-actions-panel(:actions="actions.inline", :dropDownActions="actions.dropDown")
</template>

<script>
import { pickBy } from 'lodash';

import {
  MODALS,
  ENTITIES_TYPES,
  ENTITIES_STATUSES,
  EVENT_ENTITY_TYPES,
  EVENT_ENTITY_STYLE,
  WIDGETS_ACTIONS_TYPES,
} from '@/constants';

import authMixin from '@/mixins/auth';
import entitiesAlarmMixin from '@/mixins/entities/alarm';
import widgetActionsPanelAlarmMixin from '@/mixins/widget/actions-panel/alarm';

import SharedActionsPanel from '@/components/other/shared/actions-panel/actions-panel.vue';

/**
 * Component to regroup actions (actions-panel-item) for each alarm on the alarms list
 *
 * @module alarm
 *
 * @prop {Object} item - Object representing an alarm
 * @prop {Object} widget - Full widget object
 * @prop {boolean} [isEditingMode=false] - Is editing mode enable on a view
 */
export default {
  components: { SharedActionsPanel },
  mixins: [authMixin, entitiesAlarmMixin, widgetActionsPanelAlarmMixin],
  props: {
    item: {
      type: Object,
      required: true,
    },
    widget: {
      type: Object,
      required: true,
    },
    isEditingMode: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    const { alarmsList: alarmsListActionsTypes } = WIDGETS_ACTIONS_TYPES;

    return {
      actionsMap: {
        ack: {
          type: alarmsListActionsTypes.ack,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.ack].icon,
          title: this.$t('alarmList.actions.titles.ack'),
          method: this.showActionModal(MODALS.createAckEvent),
        },
        fastAck: {
          type: alarmsListActionsTypes.fastAck,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.fastAck].icon,
          title: this.$t('alarmList.actions.titles.fastAck'),
          method: this.createAckEvent,
        },
        ackRemove: {
          type: alarmsListActionsTypes.ackRemove,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.ackRemove].icon,
          title: this.$t('alarmList.actions.titles.ackRemove'),
          method: this.showAckRemoveModal,
        },
        pbehaviorAdd: {
          type: alarmsListActionsTypes.pbehaviorAdd,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.pbehaviorAdd].icon,
          title: this.$t('alarmList.actions.titles.pbehavior'),
          method: this.showActionModal(MODALS.createPbehavior),
        },
        snooze: {
          type: alarmsListActionsTypes.snooze,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.snooze].icon,
          title: this.$t('alarmList.actions.titles.snooze'),
          method: this.showActionModal(MODALS.createSnoozeEvent),
        },
        pbehaviorList: {
          type: alarmsListActionsTypes.pbehaviorList,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.pbehaviorList].icon,
          title: this.$t('alarmList.actions.titles.pbehaviorList'),
          method: this.showActionModal(MODALS.pbehaviorList),
        },
        declareTicket: {
          type: alarmsListActionsTypes.declareTicket,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.declareTicket].icon,
          title: this.$t('alarmList.actions.titles.declareTicket'),
          method: this.showActionModal(MODALS.createDeclareTicketEvent),
        },
        associateTicket: {
          type: alarmsListActionsTypes.associateTicket,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.assocTicket].icon,
          title: this.$t('alarmList.actions.titles.associateTicket'),
          method: this.showActionModal(MODALS.createAssociateTicketEvent),
        },
        cancel: {
          type: alarmsListActionsTypes.cancel,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.delete].icon,
          title: this.$t('alarmList.actions.titles.cancel'),
          method: this.showActionModal(MODALS.createCancelEvent),
        },
        changeState: {
          type: alarmsListActionsTypes.changeState,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.changeState].icon,
          title: this.$t('alarmList.actions.titles.changeState'),
          method: this.showActionModal(MODALS.createChangeStateEvent),
        },
        moreInfos: {
          type: alarmsListActionsTypes.moreInfos,
          icon: 'more_horiz',
          title: this.$t('alarmList.actions.titles.moreInfos'),
          method: this.showMoreInfosModal(),
        },
        variablesHelp: {
          type: alarmsListActionsTypes.variablesHelp,
          icon: 'help',
          title: this.$t('alarmList.actions.titles.variablesHelp'),
          method: this.showVariablesHelperModal(),
        },
      },
    };
  },
  computed: {
    filteredActionsMap() {
      return pickBy(this.actionsMap, this.actionsAccessFilterHandler);
    },
    modalConfig() {
      return {
        itemsType: ENTITIES_TYPES.alarm,
        itemsIds: [this.item._id],
        afterSubmit: () => this.fetchAlarmsListWithPreviousParams({ widgetId: this.widget._id }),
      };
    },
    actions() {
      const { filteredActionsMap } = this;

      let inlineActions = [filteredActionsMap.pbehaviorList];
      let dropDownActions = [];

      if ([ENTITIES_STATUSES.ongoing, ENTITIES_STATUSES.flapping].includes(this.item.v.status.val)) {
        if (this.item.v.ack) {
          inlineActions = [
            filteredActionsMap.declareTicket,
            filteredActionsMap.associateTicket,
            filteredActionsMap.cancel,
          ];

          dropDownActions = [
            filteredActionsMap.ackRemove,
            filteredActionsMap.snooze,
            filteredActionsMap.changeState,
            filteredActionsMap.pbehaviorAdd,
            filteredActionsMap.pbehaviorList,
            filteredActionsMap.moreInfos,
          ];
        } else {
          inlineActions = [
            filteredActionsMap.ack,
            filteredActionsMap.fastAck,
          ];

          dropDownActions = [
            filteredActionsMap.moreInfos,
          ];
        }

        if (this.isEditingMode) {
          inlineActions.push(filteredActionsMap.variablesHelp);
        }
      }

      return {
        inline: inlineActions.filter(action => !!action),
        dropDown: dropDownActions.filter(action => !!action),
      };
    },
  },
  methods: {
    createAckEvent() {
      return this.createEvent(EVENT_ENTITY_TYPES.ack, this.item);
    },
  },
};
</script>
