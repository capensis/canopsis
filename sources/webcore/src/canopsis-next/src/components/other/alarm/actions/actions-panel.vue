<template lang="pug">
  shared-actions-panel(:actions="actions.inline", :dropDownActions="actions.dropDown")
</template>

<script>
import { pickBy, compact } from 'lodash';

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

import featuresService from '@/services/features';

/**
 * Component to regroup actions (actions-panel-item) for each alarm on the alarms list
 *
 * @module alarm
 *
 * @prop {Object} item - Object representing an alarm
 * @prop {Object} widget - Full widget object
 */
export default {
  components: { SharedActionsPanel },
  mixins: [
    authMixin,
    entitiesAlarmMixin,
    widgetActionsPanelAlarmMixin,

    ...featuresService.get('components.alarmListActionPanel.mixins', []),
  ],
  props: {
    item: {
      type: Object,
      required: true,
    },
    widget: {
      type: Object,
      required: true,
    },
    isResolvedAlarm: {
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
          method: this.showAckModal,
        },
        fastAck: {
          type: alarmsListActionsTypes.fastAck,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.fastAck].icon,
          title: this.$t('alarmList.actions.titles.fastAck'),
          method: this.createFastAckEvent,
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
          method: this.showAddPbehaviorModal,
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
          method: this.showPbehaviorsListModal,
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
        variablesHelp: {
          type: alarmsListActionsTypes.variablesHelp,
          icon: 'help',
          title: this.$t('alarmList.actions.titles.variablesHelp'),
          method: this.showVariablesHelperModal,
        },
        history: {
          type: alarmsListActionsTypes.history,
          icon: 'history',
          title: this.$t('alarmList.actions.titles.history'),
          method: this.showHistoryModal,
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
    resolvedActions() {
      const { pbehaviorList, variablesHelp } = this.filteredActionsMap;

      return [pbehaviorList, variablesHelp];
    },
    unresolvedActions() {
      const { filteredActionsMap } = this;

      const actions = [
        filteredActionsMap.snooze,
        filteredActionsMap.pbehaviorAdd,
        filteredActionsMap.pbehaviorList,
      ];

      if (this.item.entity) {
        actions.push(filteredActionsMap.history);
      }

      actions.push(filteredActionsMap.variablesHelp);

      if ([ENTITIES_STATUSES.ongoing, ENTITIES_STATUSES.flapping].includes(this.item.v.status.val)) {
        if (this.item.v.ack) {
          if (this.widget.parameters.isMultiAckEnabled) {
            actions.unshift(filteredActionsMap.ack);
          }

          actions.unshift(
            filteredActionsMap.declareTicket,
            filteredActionsMap.associateTicket,
            filteredActionsMap.cancel,
            filteredActionsMap.ackRemove,
            filteredActionsMap.changeState,
          );
        } else {
          actions.unshift(
            filteredActionsMap.ack,
            filteredActionsMap.fastAck,
          );
        }
      }

      return actions;
    },
    actions() {
      let actions = this.isResolvedAlarm ? this.resolvedActions : this.unresolvedActions;

      actions = compact(actions);

      const result = {
        inline: actions.slice(0, 3),
        dropDown: actions.slice(3),
      };

      /**
       * If we will have actions for resolved alarms in the features we should move this condition to
       * the every features repositories
       */
      if (!this.isResolvedAlarm && featuresService.has('components.alarmListActionPanel.computed.actions')) {
        return featuresService.call('components.alarmListActionPanel.computed.actions', this, result);
      }

      return result;
    },
  },
};
</script>
