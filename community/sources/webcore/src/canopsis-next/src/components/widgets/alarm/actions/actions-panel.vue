<template>
  <shared-actions-panel v-bind="$attrs" :actions="preparedActions" />
</template>

<script>
import { find } from 'lodash';

import {
  MODALS,
  ALARM_LIST_ACTIONS_TYPES,
  LINK_RULE_ACTIONS,
  REMEDIATION_INSTRUCTION_EXECUTION_STATUSES,
} from '@/constants';

import featuresService from '@/services/features';

import { getAlarmActionIcon } from '@/helpers/entities/alarm/icons';
import { isManualGroupMetaAlarmRuleType, isAutoMetaAlarmRuleType } from '@/helpers/entities/meta-alarm/rule/form';
import { isInstructionExecutionIconInProgress } from '@/helpers/entities/remediation/instruction-execution/form';
import { isInstructionManual } from '@/helpers/entities/remediation/instruction/form';
import { harmonizeLinks, getLinkRuleLinkActionType } from '@/helpers/entities/link/list';
import {
  isCancelledAlarmStatus,
  isResolvedAlarm,
  isAlarmStateOk,
  isAlarmStatusCancelled,
  isAlarmStatusClosed,
  isAlarmStatusFlapping,
  isAlarmStatusOngoing,
} from '@/helpers/entities/alarm/form';

import { entitiesAlarmMixin } from '@/mixins/entities/alarm';
import { entitiesMetaAlarmMixin } from '@/mixins/entities/meta-alarm';
import { widgetActionsPanelAlarmMixin } from '@/mixins/widget/actions-panel/alarm';
import { clipboardMixin } from '@/mixins/clipboard';
import { widgetActionsPanelAlarmExportPdfMixin } from '@/mixins/widget/actions-panel/alarm-export-pdf';

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
    entitiesMetaAlarmMixin,
    widgetActionsPanelAlarmMixin,
    clipboardMixin,
    widgetActionsPanelAlarmExportPdfMixin,
  ],
  inheritAttrs: false,
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
    refreshAlarmsList: {
      type: Function,
      default: () => {},
    },
  },
  computed: {
    isCancelledAlarm() {
      return isCancelledAlarmStatus(this.item);
    },

    isResolvedAlarm() {
      return isResolvedAlarm(this.item);
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
      if (this.isAlarmStatusFlapping) {
        return this.isAlarmStateOk;
      }

      return this.widget.parameters.isActionsAllowWithOkState && this.isAlarmStateOk;
    },

    isAlarmOpenedOrActionAllowedWithStateOk() {
      return this.isOpenedAlarm || this.isActionsAllowWithOkState;
    },

    isParentAlarmManualMetaAlarm() {
      return isManualGroupMetaAlarmRuleType(this.parentAlarm?.meta_alarm_rule?.type);
    },

    isParentAlarmAutoMetaAlarm() {
      return isAutoMetaAlarmRuleType(this.parentAlarm?.meta_alarm_rule?.type);
    },

    hasBookmark() {
      return !!this.item.bookmark;
    },

    visibleLinks() {
      return harmonizeLinks(this.item.links).filter(link => !link.hide_in_menu);
    },

    linksActions() {
      return this.visibleLinks.map((link) => {
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

    instructionsActions() {
      const {
        assigned_instructions: assignedInstructions = [],
      } = this.item;

      if (assignedInstructions.length) {
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
            cssClass,
            type: ALARM_LIST_ACTIONS_TYPES.executeInstruction,
            disabled: hasRunningInstruction
              || (Boolean(pausedInstructions.length) && !find(pausedInstructions, { _id: instruction._id })),
            title: this.$t(`remediation.instruction.${titlePrefix}Instruction`, {
              instructionName: instruction.name,
            }),
            method: () => this.showExecuteInstructionModal(instruction),
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

      if (!this.item.v?.ticket || this.widget.parameters.isMultiDeclareTicketEnabled) {
        actions.unshift({
          type: ALARM_LIST_ACTIONS_TYPES.associateTicket,
          title: this.$t('alarm.actions.titles.associateTicket'),
          method: this.showAssociateTicketModal,
        });

        if (assignedDeclareTicketRules.length) {
          actions.unshift({
            type: ALARM_LIST_ACTIONS_TYPES.declareTicket,
            title: this.$t('alarm.actions.titles.declareTicket'),
            method: this.showDeclareTicketModal,
          });
        }
      }

      return actions;
    },

    actions() {
      const actions = [];

      const isAckAndChangeStateAvailable = (this.isAlarmStatusClosed && this.isActionsAllowWithOkState)
        || this.isAlarmStatusOngoing
        || this.isAlarmStatusFlapping;
      const isNotResolvedOpenedAlarm = !this.isResolvedAlarm && this.isOpenedAlarm;

      const variablesHelpAction = {
        type: ALARM_LIST_ACTIONS_TYPES.variablesHelp,
        title: this.$t('alarm.actions.titles.variablesHelp'),
        method: this.showVariablesHelperModal,
      };
      const exportPdfAction = {
        type: ALARM_LIST_ACTIONS_TYPES.exportPdf,
        title: this.$t('alarm.actions.titles.exportPdf'),
        method: this.exportPdf,
      };
      const ackAction = {
        type: ALARM_LIST_ACTIONS_TYPES.ack,
        title: this.$t('alarm.actions.titles.ack'),
        method: this.showAckModal,
      };

      if (!this.isResolvedAlarm && isAckAndChangeStateAvailable) {
        if (this.item.v.ack) {
          actions.push(
            {
              type: ALARM_LIST_ACTIONS_TYPES.ackRemove,
              title: this.$t('alarm.actions.titles.ackRemove'),
              method: this.showAckRemoveModal,
            },
            {
              type: ALARM_LIST_ACTIONS_TYPES.changeState,
              title: this.$t('alarm.actions.titles.changeState'),
              method: this.showCreateChangeStateEventModal,
            },
          );

          if (this.widget.parameters.isMultiAckEnabled) {
            actions.push(ackAction);
          }
        } else {
          actions.push(
            ackAction,
            {
              type: ALARM_LIST_ACTIONS_TYPES.fastAck,
              title: this.$t('alarm.actions.titles.fastAck'),
              method: this.createFastAckEvent,
            },
          );
        }
      }

      if (
        !this.isResolvedAlarm && (
          /**
           * Save previous behavior
           */
          isAckAndChangeStateAvailable
          /**
           * Add behavior like in mass actions
           */
          || (
            !this.isAlarmStateOk
            && !this.isAlarmStatusClosed
            && !this.isAlarmStatusCancelled
          )
        )
      ) {
        actions.push(
          {
            type: ALARM_LIST_ACTIONS_TYPES.cancel,
            icon: '$vuetify.icons.list_delete',
            title: this.$t('alarm.actions.titles.cancel'),
            method: this.showCancelModal,
          },
          {
            type: ALARM_LIST_ACTIONS_TYPES.fastCancel,
            icon: 'delete',
            title: this.$t('alarm.actions.titles.fastCancel'),
            method: this.createFastCancel,
          },
        );
      }

      if (this.isCancelledAlarm && !this.isResolvedAlarm) {
        actions.push({
          type: ALARM_LIST_ACTIONS_TYPES.unCancel,
          title: this.$t('alarm.actions.titles.unCancel'),
          method: this.showUnCancelModal,
        });
      }

      if (!this.isResolvedAlarm && this.isAlarmOpenedOrActionAllowedWithStateOk) {
        actions.push(
          {
            type: ALARM_LIST_ACTIONS_TYPES.comment,
            title: this.$t('alarm.actions.titles.comment'),
            method: this.showCreateCommentEventModal,
          },
        );
      }

      if (!this.isResolvedAlarm && isAckAndChangeStateAvailable && this.item.v.ack) {
        actions.push(...this.ticketsActions);
      }

      if (isNotResolvedOpenedAlarm) {
        actions.push(
          {
            type: ALARM_LIST_ACTIONS_TYPES.snooze,
            title: this.$t('alarm.actions.titles.snooze'),
            method: this.showSnoozeModal,
          },
          {
            type: ALARM_LIST_ACTIONS_TYPES.fastPbehaviorAdd,
            title: this.$t('alarm.actions.titles.fastPbehaviorAdd'),
            method: this.fastAddPbehavior,
          },
          {
            type: ALARM_LIST_ACTIONS_TYPES.pbehaviorAdd,
            title: this.$t('alarm.actions.titles.pbehavior'),
            method: this.showAddPbehaviorModal,
          },
        );
      }

      actions.push(
        this.hasBookmark
          ? {
            type: ALARM_LIST_ACTIONS_TYPES.removeBookmark,
            title: this.$t('alarm.actions.titles.removeBookmark'),
            method: this.removeBookmark,
          }
          : {
            type: ALARM_LIST_ACTIONS_TYPES.addBookmark,
            title: this.$t('alarm.actions.titles.addBookmark'),
            method: this.addBookmark,
          },
      );

      if (isNotResolvedOpenedAlarm && this.item.entity) {
        actions.push({
          type: ALARM_LIST_ACTIONS_TYPES.history,
          title: this.$t('alarm.actions.titles.history'),
          method: this.showHistoryModal,
        });
      }

      if (this.isOpenedAlarm) {
        actions.push(variablesHelpAction, exportPdfAction);
      }

      if (this.isParentAlarmManualMetaAlarm) {
        actions.push({
          type: ALARM_LIST_ACTIONS_TYPES.removeAlarmsFromManualMetaAlarm,
          title: this.$t('alarm.actions.titles.removeAlarmsFromManualMetaAlarm'),
          method: this.showRemoveAlarmsFromManualMetaAlarmModal,
        });
      }

      if (!this.isResolvedAlarm && this.isAlarmOpenedOrActionAllowedWithStateOk && this.isParentAlarmAutoMetaAlarm) {
        actions.push({
          type: ALARM_LIST_ACTIONS_TYPES.removeAlarmsFromAutoMetaAlarm,
          title: this.$t('alarm.actions.titles.removeAlarmsFromAutoMetaAlarm'),
          method: this.showRemoveAlarmsFromAutoMetaAlarmModal,
        });
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

      actions.push(...this.linksActions);

      if (this.isOpenedAlarm) {
        if (!this.isResolvedAlarm) {
          actions.push(...this.instructionsActions);
        }
      } else {
        actions.push(variablesHelpAction, exportPdfAction);
      }

      return actions;
    },

    filteredActions() {
      return this.actions.filter(this.actionsAccessFilterHandler);
    },

    preparedActions() {
      return this.filteredActions.map(action => ({
        ...action,

        icon: action.icon ?? getAlarmActionIcon(action.type),
        loading: this.isActionTypeInPending(action.type),
      }));
    },
  },
  methods: {
    afterSubmit() {
      this.refreshAlarmsList();
    },

    showCreateChangeStateEventModal() {
      this.showCreateChangeStateEventModalByAlarms([this.item]);
    },

    showSnoozeModal() {
      this.showSnoozeModalByAlarms([this.item]);
    },

    showAckModal() {
      this.showAckModalByAlarms([this.item]);
    },

    createFastAckEvent() {
      this.createFastAckActionByAlarms([this.item]);
    },

    async exportPdf() {
      this.setActionPending(ALARM_LIST_ACTIONS_TYPES.exportPdf, true);
      await this.exportAlarmToPdf(this.item, this.widget.parameters.exportPdfTemplate);
      this.setActionPending(ALARM_LIST_ACTIONS_TYPES.exportPdf, false);
    },

    showAssociateTicketModal() {
      this.showAssociateTicketModalByAlarms([this.item]);
    },

    showDeclareTicketModal() {
      this.showDeclareTicketModalByAlarms([this.item]);
    },

    showCreateCommentEventModal() {
      this.showCreateCommentModalByAlarms([this.item]);
    },

    showAckRemoveModal() {
      this.showAckRemoveModalByAlarms([this.item]);
    },

    showCancelModal() {
      this.showCancelModalByAlarms([this.item]);
    },

    showUnCancelModal() {
      this.showUnCancelModalByAlarms([this.item]);
    },

    createFastCancel() {
      this.createFastCancelActionByAlarms([this.item]);
    },

    showRemoveAlarmsFromManualMetaAlarmModal() {
      this.showRemoveAlarmsFromManualMetaAlarmModalByAlarms([this.item]);
    },

    showRemoveAlarmsFromAutoMetaAlarmModal() {
      this.showRemoveAlarmsFromAutoMetaAlarmModalByAlarms([this.item]);
    },

    showVariablesHelperModal() {
      this.showVariablesHelperModalByAlarm(this.item);
    },

    showAddPbehaviorModal() {
      this.showAddPbehaviorModalByAlarms([this.item]);
    },

    fastAddPbehavior() {
      this.addFastPbehaviorByAlarms([this.item]);
    },

    showHistoryModal() {
      this.showHistoryModalByAlarm(this.item);
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

    async addBookmark() {
      return this.addBookmarkByAlarm(this.item);
    },

    async removeBookmark() {
      return this.removeBookmarkByAlarm(this.item);
    },
  },
};
</script>
