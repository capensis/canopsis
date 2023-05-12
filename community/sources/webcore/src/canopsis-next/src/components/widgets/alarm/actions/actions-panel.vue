<template lang="pug">
  shared-actions-panel(:actions="actions", :small="small")
</template>

<script>
import { pickBy, compact, find } from 'lodash';

import {
  MODALS,
  EVENT_ENTITY_TYPES,
  ALARM_LIST_ACTIONS_TYPES,
  LINK_RULE_ACTIONS,
  REMEDIATION_INSTRUCTION_EXECUTION_STATUSES,
} from '@/constants';

import { getEntityEventIcon } from '@/helpers/icon';

import featuresService from '@/services/features';

import { isManualGroupMetaAlarmRuleType } from '@/helpers/forms/meta-alarm-rule';
import { isInstructionExecutionIconInProgress } from '@/helpers/forms/remediation-instruction-execution';
import { isInstructionManual } from '@/helpers/forms/remediation-instruction';
import { harmonizeLinks, getLinkRuleLinkActionType } from '@/helpers/links';
import {
  isAlarmStateOk,
  isAlarmStatusCancelled,
  isAlarmStatusClosed,
  isAlarmStatusFlapping,
  isAlarmStatusOngoing,
} from '@/helpers/entities/alarm';

import { entitiesAlarmMixin } from '@/mixins/entities/alarm';
import { widgetActionsPanelAlarmMixin } from '@/mixins/widget/actions-panel/alarm';
import { clipboardMixin } from '@/mixins/clipboard';

import SharedActionsPanel from '@/components/common/actions-panel/actions-panel.vue';

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
    entitiesAlarmMixin,
    widgetActionsPanelAlarmMixin,
    clipboardMixin,
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
    small: {
      type: Boolean,
      default: false,
    },
    refreshAlarmsList: {
      type: Function,
      default: () => {},
    },
  },
  computed: {
    actionsMap() {
      /**
       * !!!IMPORTANT!!! TODO: We need check all features
       */
      const featuresActionsMap = featuresService.has('components.alarmListActionPanel.computed.actionsMap')
        ? featuresService.call('components.alarmListActionPanel.computed.actionsMap', this, [])
        : {};

      return {
        ack: {
          type: ALARM_LIST_ACTIONS_TYPES.ack,
          icon: getEntityEventIcon(EVENT_ENTITY_TYPES.ack),
          title: this.$t('alarm.actions.titles.ack'),
          method: this.showAckModal,
        },
        fastAck: {
          type: ALARM_LIST_ACTIONS_TYPES.fastAck,
          icon: getEntityEventIcon(EVENT_ENTITY_TYPES.fastAck),
          title: this.$t('alarm.actions.titles.fastAck'),
          method: this.createFastAckEvent,
        },
        ackRemove: {
          type: ALARM_LIST_ACTIONS_TYPES.ackRemove,
          icon: getEntityEventIcon(EVENT_ENTITY_TYPES.ackRemove),
          title: this.$t('alarm.actions.titles.ackRemove'),
          method: this.showAckRemoveModal,
        },
        pbehaviorAdd: {
          type: ALARM_LIST_ACTIONS_TYPES.pbehaviorAdd,
          icon: getEntityEventIcon(EVENT_ENTITY_TYPES.pbehaviorAdd),
          title: this.$t('alarm.actions.titles.pbehavior'),
          method: this.showAddPbehaviorModal,
        },
        snooze: {
          type: ALARM_LIST_ACTIONS_TYPES.snooze,
          icon: getEntityEventIcon(EVENT_ENTITY_TYPES.snooze),
          title: this.$t('alarm.actions.titles.snooze'),
          method: this.showSnoozeModal,
        },
        declareTicket: {
          type: ALARM_LIST_ACTIONS_TYPES.declareTicket,
          icon: getEntityEventIcon(EVENT_ENTITY_TYPES.declareTicket),
          title: this.$t('alarm.actions.titles.declareTicket'),
          loading: this.ticketsForAlarmsPending,
          method: this.showDeclareTicketModal,
        },
        associateTicket: {
          type: ALARM_LIST_ACTIONS_TYPES.associateTicket,
          icon: getEntityEventIcon(EVENT_ENTITY_TYPES.assocTicket),
          title: this.$t('alarm.actions.titles.associateTicket'),
          method: this.showAssociateTicketModal,
        },
        cancel: {
          type: ALARM_LIST_ACTIONS_TYPES.cancel,
          icon: '$vuetify.icons.list_delete',
          title: this.$t('alarm.actions.titles.cancel'),
          method: this.showCancelEventModal,
        },
        fastCancel: {
          type: ALARM_LIST_ACTIONS_TYPES.fastCancel,
          icon: 'delete',
          title: this.$t('alarm.actions.titles.fastCancel'),
          method: this.createFastCancelEvent,
        },
        changeState: {
          type: ALARM_LIST_ACTIONS_TYPES.changeState,
          icon: getEntityEventIcon(EVENT_ENTITY_TYPES.changeState),
          title: this.$t('alarm.actions.titles.changeState'),
          method: this.showActionModal(MODALS.createChangeStateEvent),
        },
        variablesHelp: {
          type: ALARM_LIST_ACTIONS_TYPES.variablesHelp,
          icon: 'help',
          title: this.$t('alarm.actions.titles.variablesHelp'),
          method: this.showVariablesHelperModal,
        },
        history: {
          type: ALARM_LIST_ACTIONS_TYPES.history,
          icon: 'history',
          title: this.$t('alarm.actions.titles.history'),
          method: this.showHistoryModal,
        },
        comment: {
          type: ALARM_LIST_ACTIONS_TYPES.comment,
          icon: getEntityEventIcon(EVENT_ENTITY_TYPES.comment),
          title: this.$t('alarm.actions.titles.comment'),
          method: this.showCreateCommentEventModal,
        },
        removeAlarmsFromManualMetaAlarm: {
          type: ALARM_LIST_ACTIONS_TYPES.removeAlarmsFromManualMetaAlarm,
          icon: getEntityEventIcon(EVENT_ENTITY_TYPES.removeAlarmsFromManualMetaAlarm),
          title: this.$t('alarm.actions.titles.removeAlarmsFromManualMetaAlarm'),
          method: this.showRemoveAlarmsFromManualMetaAlarmModal,
        },
        executeInstruction: {
          type: ALARM_LIST_ACTIONS_TYPES.executeInstruction,
          icon: getEntityEventIcon(EVENT_ENTITY_TYPES.executeInstruction),
          method: this.showExecuteInstructionModal,
        },
        ...featuresActionsMap,
      };
    },

    isAlarmStatusClosed() {
      return isAlarmStatusClosed(this.item);
    },

    isAlarmStatusCancelled() {
      return isAlarmStatusCancelled(this.item);
    },

    isAlarmStatusOngoing() {
      return isAlarmStatusOngoing(this.item);
    },

    isAlarmStatusFlapping() {
      return isAlarmStatusFlapping(this.item);
    },

    isOpenedAlarm() {
      return !this.isAlarmStatusClosed && !this.isAlarmStatusCancelled;
    },

    isAlarmStateOk() {
      return isAlarmStateOk(this.item);
    },

    isActionsAllowWithOkState() {
      return this.widget.parameters.isActionsAllowWithOkState && this.isAlarmStateOk;
    },

    isAlarmOpenedOrActionAllowedWithStateOk() {
      return this.isOpenedAlarm || this.isActionsAllowWithOkState;
    },

    isParentAlarmManualMetaAlarm() {
      return isManualGroupMetaAlarmRuleType(this.parentAlarm?.meta_alarm_rule?.type);
    },

    filteredActionsMap() {
      return pickBy(this.actionsMap, this.actionsAccessFilterHandler);
    },

    linksActions() {
      return harmonizeLinks(this.item.links).map((link) => {
        const type = getLinkRuleLinkActionType(link);

        return {
          type,
          icon: link.icon_name,
          title: this.$t('alarm.followLink', { title: link.label }),
          method: link.action === LINK_RULE_ACTIONS.copy
            ? () => this.writeTextToClipboard(link.url)
            : () => window.open(link.url, '_blank'),
        };
      });
    },

    modalConfig() {
      return {
        items: [this.item],
        afterSubmit: this.afterSubmit,
      };
    },

    instructionsActions() {
      const {
        assigned_instructions: assignedInstructions = [],
      } = this.item;

      if (assignedInstructions.length && this.filteredActionsMap.executeInstruction) {
        const pausedInstructions = assignedInstructions.filter(instruction => instruction.execution);
        const hasRunningInstruction = isInstructionExecutionIconInProgress(this.item.instruction_execution_icon);

        return assignedInstructions.map((instruction) => {
          const { execution } = instruction;
          let titlePrefix = 'execute';
          let cssClass = '';

          if (execution) {
            if (execution.status === REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.running) {
              titlePrefix = 'inProgress';
              cssClass = 'font-italic';
            } else {
              titlePrefix = 'resume';
            }
          }

          return {
            ...this.filteredActionsMap.executeInstruction,

            cssClass,
            disabled: hasRunningInstruction
              || (Boolean(pausedInstructions.length) && !find(pausedInstructions, { _id: instruction._id })),
            title: this.$t(`remediation.instruction.${titlePrefix}Instruction`, {
              instructionName: instruction.name,
            }),
            method: () => this.filteredActionsMap.executeInstruction.method(instruction),
          };
        });
      }

      return [];
    },

    ticketsActions() {
      const actions = [];
      const {
        assigned_declare_ticket_rules: assignedDeclareTicketRules = [],
      } = this.item;

      if (!this.item.v?.tickets?.length || this.widget.parameters.isMultiDeclareTicketEnabled) {
        actions.unshift(this.filteredActionsMap.associateTicket);

        if (assignedDeclareTicketRules.length) {
          actions.unshift(this.filteredActionsMap.declareTicket);
        }
      }

      return actions;
    },

    actions() {
      const { filteredActionsMap } = this;
      const actions = [];

      if (this.isOpenedAlarm) {
        actions.push(
          filteredActionsMap.snooze,
          filteredActionsMap.pbehaviorAdd,
        );
      }

      if (this.isAlarmOpenedOrActionAllowedWithStateOk) {
        actions.push(
          filteredActionsMap.comment,
        );
      }

      if (this.isOpenedAlarm && this.item.entity) {
        actions.push(filteredActionsMap.history);
      }

      if (this.isOpenedAlarm) {
        actions.push(filteredActionsMap.variablesHelp);
      }

      if (this.isAlarmOpenedOrActionAllowedWithStateOk && this.isParentAlarmManualMetaAlarm) {
        actions.push(filteredActionsMap.removeAlarmsFromManualMetaAlarm);
      }

      /**
         * If we will have actions for resolved alarms in the features we should move this condition to
         * the every features repositories
         */
      if (
        this.isOpenedAlarm
        && featuresService.has('components.alarmListActionPanel.computed.actions')
      ) {
        const featuresActions = featuresService.call('components.alarmListActionPanel.computed.actions', this, []);

        if (featuresActions?.length) {
          actions.unshift(...featuresActions);
        }
      }

      if (
        (this.isAlarmStatusClosed && this.isActionsAllowWithOkState)
        || this.isAlarmStatusOngoing
        || this.isAlarmStatusFlapping
      ) {
        if (this.item.v.ack) {
          if (this.widget.parameters.isMultiAckEnabled) {
            actions.unshift(filteredActionsMap.ack);
          }

          actions.unshift(
            filteredActionsMap.ackRemove,
            filteredActionsMap.changeState,
          );

          if (!this.isAlarmStateOk) {
            actions.unshift(
              filteredActionsMap.cancel,
              filteredActionsMap.fastCancel,
            );
          }

          actions.unshift(...this.ticketsActions);
        } else {
          actions.unshift(
            filteredActionsMap.ack,
            filteredActionsMap.fastAck,
          );
        }
      }

      actions.push(...this.linksActions);

      /**
         * Add actions for available instructions
         */
      if (this.isOpenedAlarm && filteredActionsMap.executeInstruction) {
        actions.push(...this.instructionsActions);
      }

      if (!this.isOpenedAlarm) {
        actions.push(filteredActionsMap.variablesHelp);
      }

      return compact(actions);
    },
  },
  methods: {
    afterSubmit() {
      this.refreshAlarmsList();
    },

    showExecuteInstructionModal(assignedInstruction) {
      const refreshAlarm = () => this.refreshAlarmsList();

      this.$modals.show({
        id: `${this.item._id}${assignedInstruction._id}`,
        name: isInstructionManual(assignedInstruction)
          ? MODALS.executeRemediationInstruction
          : MODALS.executeRemediationSimpleInstruction,
        config: {
          assignedInstruction,
          alarmId: this.item._id,
          onClose: refreshAlarm,
          onComplete: refreshAlarm,
          onExecute: refreshAlarm,
        },
      });
    },

    showAssociateTicketModal() {
      this.showAssociateTicketModalByAlarms([this.item]);
    },

    showDeclareTicketModal() {
      this.showDeclareTicketModalByAlarms([this.item]);
    },

    showCreateCommentEventModal() {
      this.$modals.show({
        name: MODALS.createCommentEvent,
        config: {
          ...this.modalConfig,
          action: data => this.createEvent(EVENT_ENTITY_TYPES.comment, this.item, data),
        },
      });
    },

    showRemoveAlarmsFromManualMetaAlarmModal() {
      this.$modals.show({
        name: MODALS.removeAlarmsFromManualMetaAlarm,
        config: {
          ...this.modalConfig,

          title: this.$t('alarm.actions.titles.removeAlarmsFromManualMetaAlarm'),
          parentAlarm: this.parentAlarm,
        },
      });
    },

    async createFastAckEvent() {
      await this.createFastAckActionByAlarms([this.item]);

      return this.refreshAlarmsList();
    },

    async createFastCancelEvent() {
      await this.createFastCancelActionByAlarms([this.item]);

      return this.refreshAlarmsList();
    },
  },
};
</script>
