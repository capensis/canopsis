<template lang="pug">
  shared-actions-panel(:actions="actions.inline", :dropDownActions="actions.dropDown")
</template>

<script>
import { get, pickBy, compact } from 'lodash';

import {
  MODALS,
  ENTITIES_TYPES,
  ENTITIES_STATUSES,
  EVENT_ENTITY_TYPES,
  EVENT_ENTITY_STYLE,
  WIDGETS_ACTIONS_TYPES,
  META_ALARMS_RULE_TYPES,
  REMEDIATION_INSTRUCTION_EXECUTION_STATUSES,
} from '@/constants';

import authMixin from '@/mixins/auth';
import entitiesAlarmMixin from '@/mixins/entities/alarm';
import widgetActionsPanelAlarmMixin from '@/mixins/widget/actions-panel/alarm';

import SharedActionsPanel from '@/components/common/actions-panel/actions-panel.vue';

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
    parentAlarm: {
      type: Object,
      default: null,
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
          method: this.showSnoozeModal,
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
          method: this.showCancelEventModal,
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
        comment: {
          type: alarmsListActionsTypes.comment,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.comment].icon,
          title: this.$t('alarmList.actions.titles.comment'),
          method: this.showCreateCommentModal,
        },
        manualMetaAlarmUngroup: {
          type: alarmsListActionsTypes.manualMetaAlarmUngroup,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.manualMetaAlarmUngroup].icon,
          title: this.$t('alarmList.actions.titles.manualMetaAlarmUngroup'),
          method: this.showManualMetaAlarmUngroupModal,
        },
        executeInstruction: {
          type: alarmsListActionsTypes.executeInstruction,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.executeInstruction].icon,
          method: this.showExecuteInstructionModal,
        },
      },
    };
  },
  computed: {
    isParentAlarmManualMetaAlarm() {
      return get(this.parentAlarm, 'rule.type') === META_ALARMS_RULE_TYPES.manualgroup;
    },
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
      const { assigned_instructions: assignedInstructions = [] } = this.item;

      const actions = [
        filteredActionsMap.snooze,
        filteredActionsMap.pbehaviorAdd,
        filteredActionsMap.pbehaviorList,
        filteredActionsMap.comment,
      ];

      if (this.item.entity) {
        actions.push(filteredActionsMap.history);
      }

      actions.push(filteredActionsMap.variablesHelp);

      if (this.isParentAlarmManualMetaAlarm) {
        actions.push(filteredActionsMap.manualMetaAlarmUngroup);
      }

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

      /**
       * Add actions for available instructions
       */
      if (assignedInstructions.length && filteredActionsMap.executeInstruction) {
        assignedInstructions.forEach((instruction) => {
          const { execution } = instruction;
          const titlePrefix = execution ? 'resume' : 'execute';

          const action = {
            ...filteredActionsMap.executeInstruction,

            disabled: get(execution, 'status') === REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.running,
            title: this.$t(`alarmList.actions.titles.${titlePrefix}Instruction`, {
              instructionName: instruction.name,
            }),
            method: () => filteredActionsMap.executeInstruction.method(instruction),
          };

          actions.push(action);
        });
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
  methods: {
    async showExecuteInstructionModal(assignedInstruction) {
      const refreshAlarm = () => this.refreshAlarmById(this.item._id);

      this.$modals.show({
        id: `${this.item._id}${assignedInstruction._id}`,
        name: MODALS.executeRemediationInstruction,
        config: {
          assignedInstruction,
          alarm: this.item,
          onOpen: refreshAlarm,
          onClose: refreshAlarm,
          onComplete: async (instructionExecute) => {
            await refreshAlarm();
            this.showRateInstructionModal(instructionExecute._id);
          },
        },
      });
    },
  },
};
</script>
