<template>
  <actions-panel
    ref="sharedActionsPanel"
    v-bind="additionalProps"
    :actions="preparedActions"
    :small="small"
  />
</template>

<script>
import { find, keyBy } from 'lodash';
import { mapActions } from 'pinia';

import {
  ALARM_EXPORT_FILE_NAME_PREFIX,
  ALARM_LIST_ACTIONS_TYPES,
  ALARM_LIST_FAST_PBEHAVIOR_TIMEOUT,
  ALARM_LIST_TOGGLE_ACTIONS_TYPES_MAP,
  BUSINESS_USER_PERMISSIONS_ACTIONS_MAP,
  DEFAULT_ALARM_ACTIONS_INLINE_COUNT,
  LINK_RULE_ACTIONS,
  MODALS,
  PBEHAVIOR_ORIGINS,
  PBEHAVIOR_TYPE_TYPES,
  REMEDIATION_INSTRUCTION_EXECUTION_STATUSES,
} from '@/constants';

import { featuresService } from '@/services/features';

import { useAlarmStore } from '@/stores/modules/alarm';
import { useDeclareTicketRuleStore } from '@/stores/modules/declare-ticket-rule';
import { useMetaAlarmStore } from '@/stores/modules/meta-alarm';
import { usePbehaviorStore } from '@/stores/modules/pbehavior';
import { usePbehaviorReasonStore } from '@/stores/modules/pbehavior-reason';
import { usePbehaviorTypeStore } from '@/stores/modules/pbehavior-type';

import { sortActionsByQuickActions, getActionsInlineCount } from '@/helpers/actions-panel';
import { mapIds } from '@/helpers/array';
import { promisedWait } from '@/helpers/async';
import {
  formToChangeStateEvent,
  isAlarmStateOk,
  isAlarmStatusCancelled,
  isAlarmStatusClosed,
  isAlarmStatusFlapping,
  isAlarmStatusNoEvents,
  isAlarmStatusOngoing,
  isAlarmStatusUnknown,
  isResolvedAlarm,
  mapAlarmsEntities,
} from '@/helpers/entities/alarm/form';
import { getAlarmActionIcon } from '@/helpers/entities/alarm/icons';
import { harmonizeLinks, getLinkRuleLinkActionType } from '@/helpers/entities/link/list';
import {
  isAutoMetaAlarmRuleType,
  isManualGroupMetaAlarmRuleType,
} from '@/helpers/entities/meta-alarm/rule/form';
import { createEntityIdPatternByValue } from '@/helpers/entities/pattern/form';
import { createDowntimePbehavior } from '@/helpers/entities/pbehavior/form';
import { isInstructionExecutionIconInProgress } from '@/helpers/entities/remediation/instruction-execution/form';
import { isInstructionTypeManual } from '@/helpers/entities/remediation/instruction/form';
import { generatePreparedDefaultAlarmListWidget } from '@/helpers/entities/widget/form';
import { exportAlarmToPdf as exportAlarmToPdfFile } from '@/helpers/file/pdf';
import { convertObjectToTreeview } from '@/helpers/treeview';

import { authMixin } from '@/mixins/auth';

export default {
  inject: ['$system'],

  mixins: [authMixin],

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
    inlineCount: {
      type: Number,
      default: DEFAULT_ALARM_ACTIONS_INLINE_COUNT,
    },
    ignoreMediaQuery: {
      type: Boolean,
      default: false,
    },
    small: {
      type: Boolean,
      default: false,
    },
  },

  data() {
    return {
      pendingByActionsTypes: {},
    };
  },

  computed: {
    quickActions() {
      return this.widget.parameters.quickActions ?? [];
    },

    isCancelledAlarm() {
      return isAlarmStatusCancelled(this.item);
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

    isAlarmStatusNoEvents() {
      return isAlarmStatusNoEvents(this.item);
    },

    isAlarmStatusFlapping() {
      return isAlarmStatusFlapping(this.item);
    },

    isAlarmStatusUnknown() {
      return isAlarmStatusUnknown(this.item);
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

    isAckAndChangeStateAvailable() {
      return (
        (this.isAlarmStatusClosed && this.isActionsAllowWithOkState) ||
        this.isAlarmStatusOngoing ||
        this.isAlarmStatusNoEvents ||
        this.isAlarmStatusFlapping
      );
    },

    hasBookmark() {
      return !!this.item.bookmark;
    },

    visibleLinks() {
      return harmonizeLinks(this.item.links).filter(link => !link.hide_in_menu);
    },

    linksActions() {
      if (!this.checkAccess(BUSINESS_USER_PERMISSIONS_ACTIONS_MAP.alarmsList[ALARM_LIST_ACTIONS_TYPES.links])) {
        return [];
      }

      return this.visibleLinks.map((link) => {
        const type = getLinkRuleLinkActionType(link);

        return {
          type,
          link: true,
          icon: link.icon_name,
          title: this.$t('alarm.followLink', { title: link.label }),
          method: () => this.handleLinkClickAction(link, type),
        };
      });
    },

    instructionsActions() {
      const { assigned_instructions: assignedInstructions = [] } = this.item;

      if (!assignedInstructions.length) {
        return [];
      }

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
          disabled:
            hasRunningInstruction ||
            (Boolean(pausedInstructions.length) && !find(pausedInstructions, { _id: instruction._id })),
          title: this.$t(`remediation.instruction.${titlePrefix}Instruction`, {
            instructionName: instruction.name,
          }),
          method: () => this.showExecuteInstructionModal(instruction),
        };
      });
    },

    ticketsActions() {
      const actions = [];
      const { assigned_declare_ticket_rules: assignedDeclareTicketRules = [] } = this.item;
      const hasAssociatedTickets = this.item.v?.tickets?.length > 0;

      if (hasAssociatedTickets) {
        actions.unshift({
          type: ALARM_LIST_ACTIONS_TYPES.removeAssociatedTicket,
          iconSize: 20,
          title: this.$t('alarm.actions.titles.removeAssociatedTicket'),
          method: this.showRemoveAssociatedTicketModal,
        });
      }

      if (this.isResolvedAlarm || !this.isAckAndChangeStateAvailable || !this.item.v.ack) {
        return actions;
      }

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

    fastPbehaviorAction() {
      if (this.item.pbh_origin_icon) {
        return {
          type: ALARM_LIST_ACTIONS_TYPES.fastPbehaviorRemove,
          title: this.$t('alarm.actions.titles.fastPbehaviorRemove'),
          method: this.fastRemovePbehavior,
        };
      }

      const fastPbehaviorsParameters = this.widget.parameters.fast_pbehaviors ?? [];
      const fastPbehaviorAction = {
        type: ALARM_LIST_ACTIONS_TYPES.fastPbehaviorAdd,
        title: this.$t('alarm.actions.titles.fastPbehaviorAdd'),
      };

      if (fastPbehaviorsParameters.length > 1) {
        fastPbehaviorAction.items = fastPbehaviorsParameters.map(pbehaviorParameters => ({
          title: pbehaviorParameters.name_prefix,
          method: () => this.fastAddPbehavior(pbehaviorParameters),
        }));
      } else {
        fastPbehaviorAction.method = () => this.fastAddPbehavior(fastPbehaviorsParameters[0]);
      }

      return fastPbehaviorAction;
    },

    actions() {
      const actions = [];
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

      if (!this.isResolvedAlarm && this.isAckAndChangeStateAvailable) {
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
            actions.push(this.ackAction, this.fastAckAction);
          }
        } else {
          actions.push(this.ackAction, this.fastAckAction);
        }
      }

      if (
        !this.isAlarmStatusUnknown &&
        !this.isResolvedAlarm &&
        (this.isAckAndChangeStateAvailable ||
          (!this.isAlarmStateOk && !this.isAlarmStatusClosed && !this.isAlarmStatusCancelled))
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
        actions.push({
          type: ALARM_LIST_ACTIONS_TYPES.comment,
          title: this.$t('alarm.actions.titles.comment'),
          method: this.showCreateCommentEventModal,
        });
      }

      actions.push(...this.ticketsActions);

      if (isNotResolvedOpenedAlarm) {
        actions.push(
          {
            type: ALARM_LIST_ACTIONS_TYPES.snooze,
            title: this.$t('alarm.actions.titles.snooze'),
            method: this.showSnoozeModal,
          },
          this.fastPbehaviorAction,
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

      if (this.isOpenedAlarm && featuresService.has('components.alarmListActionPanel.computed.actions')) {
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

    ackAction() {
      return {
        type: ALARM_LIST_ACTIONS_TYPES.ack,
        title: this.$t('alarm.actions.titles.ack'),
        method: this.showAckModal,
      };
    },

    fastAckAction() {
      return {
        type: ALARM_LIST_ACTIONS_TYPES.fastAck,
        title: this.$t('alarm.actions.titles.fastAck'),
        method: this.createFastAckEvent,
      };
    },

    filteredActions() {
      return this.actions.filter(this.actionsAccessFilterHandler);
    },

    preparedActions() {
      return sortActionsByQuickActions(this.filteredActions, this.quickActions).map(action => ({
        ...action,
        icon: action.icon ?? getAlarmActionIcon(action.type),
        loading: this.isActionTypeInPending(action.type),
      }));
    },

    additionalProps() {
      if (this.ignoreMediaQuery) {
        return {
          inlineCount: this.inlineCount,
          ignoreMediaQuery: this.ignoreMediaQuery,
        };
      }

      return {
        inlineCount: getActionsInlineCount(this.preparedActions, this.quickActions) || this.inlineCount,
        ignoreMediaQuery: false,
      };
    },
  },

  methods: {
    ...mapActions(useAlarmStore, [
      'bulkCreateAlarmAckEvent',
      'bulkCreateAlarmAckremoveEvent',
      'bulkCreateAlarmSnoozeEvent',
      'bulkCreateAlarmAssocticketEvent',
      'bulkCreateAlarmTicketremoveEvent',
      'bulkCreateAlarmCommentEvent',
      'bulkCreateAlarmCancelEvent',
      'bulkCreateAlarmUnCancelEvent',
      'bulkCreateAlarmChangestateEvent',
      'addBookmarkToAlarm',
      'removeBookmarkFromAlarm',
      'fetchResolvedAlarmsListWithoutStore',
      'fetchAlarmLinksWithoutStore',
      'fetchAlarmDetailsListWithoutStore',
    ]),
    ...mapActions(useDeclareTicketRuleStore, {
      fetchAssignedTicketsWithoutStore: 'fetchAssignedTicketsWithoutStore',
    }),
    ...mapActions(useMetaAlarmStore, {
      createMetaAlarm: 'create',
      addAlarmsIntoMetaAlarm: 'addAlarms',
      removeAlarmsFromMetaAlarm: 'removeAlarms',
    }),
    ...mapActions(usePbehaviorStore, [
      'bulkCreateEntityPbehaviors',
      'bulkRemoveEntityPbehaviors',
    ]),
    ...mapActions(usePbehaviorTypeStore, {
      fetchDefaultPbehaviorTypes: 'fetchListWithoutStore',
    }),
    ...mapActions(usePbehaviorReasonStore, {
      fetchPbehaviorReasonsListWithoutStore: 'fetchListWithoutStore',
    }),

    afterSubmit() {
      this.refreshAlarmsList();
      this.$refs.sharedActionsPanel?.closeMenu?.();
    },

    isActionTypeInPending(type) {
      return (
        !!this.pendingByActionsTypes[type] ||
        !!this.pendingByActionsTypes[ALARM_LIST_TOGGLE_ACTIONS_TYPES_MAP[type]]
      );
    },

    setActionPending(type, value) {
      this.pendingByActionsTypes = {
        ...this.pendingByActionsTypes,
        [type]: value,
      };
    },

    async exportAlarmToPdf(alarm, template) {
      try {
        const response = await this.fetchAlarmDetailsListWithoutStore({
          params: [{
            _id: alarm._id,
            steps: { type: 'comment' },
          }],
        });

        const comments = response?.[0]?.data?.steps?.data ?? [];

        await exportAlarmToPdfFile(template, { ...alarm, comments }, this.$system.timezone);
      } catch (err) {
        console.error(err);
        this.$popups.error({ text: this.$t('errors.default') });
      }
    },

    async exportPdf() {
      this.setActionPending(ALARM_LIST_ACTIONS_TYPES.exportPdf, true);
      await this.exportAlarmToPdf(this.item, this.widget.parameters.exportPdfTemplate);
      this.setActionPending(ALARM_LIST_ACTIONS_TYPES.exportPdf, false);
    },

    showCreateChangeStateEventModal() {
      this.showEventModal(
        MODALS.createChangeStateEvent,
        [this.item],
        this.bulkCreateAlarmChangestateEvent,
        { prepareEvent: formToChangeStateEvent },
      );
    },

    showAckModal() {
      this.$modals.show({
        name: MODALS.createAckEvent,
        config: {
          items: [this.item],
          isNoteRequired: this.widget.parameters.isAckNoteRequired,
          action: async (event, { needDeclareTicket, needAssociateTicket } = {}) => {
            await this.bulkCreateAlarmAckEvent({ data: [{ ...event, _id: this.item._id }] });
            await this.afterSubmit();

            if (needAssociateTicket) {
              this.showAssociateTicketModal(true);
            } else if (needDeclareTicket && this.item.assigned_declare_ticket_rules?.length) {
              await this.showDeclareTicketModal();
            }
          },
        },
      });
    },

    async createFastAckEvent() {
      this.setActionPending(ALARM_LIST_ACTIONS_TYPES.fastAck, true);

      try {
        const { fastAckOutput } = this.widget.parameters;
        const event = { comment: fastAckOutput?.enabled ? fastAckOutput.value : '' };

        await this.bulkCreateAlarmAckEvent({ data: [{ ...event, _id: this.item._id }] });
        await this.afterSubmit();
      } catch (err) {
        console.error(err);
      } finally {
        this.setActionPending(ALARM_LIST_ACTIONS_TYPES.fastAck, false);
      }
    },

    showSnoozeModal() {
      this.showEventModal(MODALS.createSnoozeEvent, [this.item], this.bulkCreateAlarmSnoozeEvent, {
        isNoteRequired: this.widget.parameters.isSnoozeNoteRequired,
      });
    },

    showAssociateTicketModal(ignoreAck = false) {
      this.$modals.show({
        name: MODALS.createAssociateTicketEvent,
        config: {
          items: [this.item],
          ignoreAck,
          action: async (associateEvent) => {
            if (!ignoreAck && !this.item.v.ack) {
              const { fastAckOutput } = this.widget.parameters;

              await this.bulkCreateAlarmAckEvent({
                data: [{
                  comment: fastAckOutput?.enabled ? fastAckOutput.value : '',
                  _id: this.item._id,
                }],
              });
            }

            await this.bulkCreateAlarmAssocticketEvent({
              data: [{ ...associateEvent, _id: this.item._id }],
            });

            await this.afterSubmit();
          },
        },
      });
    },

    showRemoveAssociatedTicketModal() {
      this.showEventModal(
        MODALS.removeAssociatedTicketEvent,
        [this.item],
        this.bulkCreateAlarmTicketremoveEvent,
      );
    },

    async showDeclareTicketModal() {
      this.setActionPending(ALARM_LIST_ACTIONS_TYPES.declareTicket, true);

      try {
        const response = await this.fetchAssignedTicketsWithoutStore({
          params: { alarms: [this.item._id] },
        });

        this.$modals.show({
          name: MODALS.createDeclareTicketEvent,
          config: {
            items: [this.item],
            alarmsByTickets: response.by_rules,
            ticketsByAlarms: response.by_alarms,
            action: (events, singleMode) => {
              this.$modals.show({
                name: MODALS.executeDeclareTickets,
                config: {
                  singleMode,
                  executions: events,
                  tickets: events.map(({ _id: id }) => ({
                    _id: id,
                    name: response.by_rules[id].name,
                  })),
                  alarms: [this.item],
                  onExecute: this.afterSubmit,
                },
              });
            },
          },
        });
      } catch (err) {
        console.error(err);
      } finally {
        this.setActionPending(ALARM_LIST_ACTIONS_TYPES.declareTicket, false);
      }
    },

    showCreateCommentEventModal() {
      const commentTemplatesById = keyBy(this.widget.comment_templates, '_id');

      this.showEventModal(MODALS.createCommentEvent, [this.item], this.bulkCreateAlarmCommentEvent, {
        templates: (this.widget.parameters.comment_templates ?? [])
          .map(id => commentTemplatesById[id])
          .filter(Boolean),
      });
    },

    showAckRemoveModal() {
      this.showEventModal(MODALS.createEvent, [this.item], this.bulkCreateAlarmAckremoveEvent, {
        title: this.$t('modals.createAckRemove.title'),
      });
    },

    showCancelModal() {
      this.showEventModal(MODALS.createEvent, [this.item], this.bulkCreateAlarmCancelEvent, {
        title: this.$t('modals.createCancelEvent.title'),
      });
    },

    showUnCancelModal() {
      this.showEventModal(MODALS.createEvent, [this.item], this.bulkCreateAlarmUnCancelEvent, {
        title: this.$t('modals.createUnCancel.title'),
        isCommentRequired: this.widget.parameters.isUncancelAlarmsCommentRequired,
      });
    },

    async createFastCancel() {
      this.setActionPending(ALARM_LIST_ACTIONS_TYPES.fastCancel, true);

      try {
        const { fastCancelOutput } = this.widget.parameters;
        const event = { comment: fastCancelOutput?.enabled ? fastCancelOutput.value : '' };

        await this.bulkCreateAlarmCancelEvent({ data: [{ ...event, _id: this.item._id }] });
        await this.afterSubmit();
      } catch (err) {
        console.error(err);
      } finally {
        this.setActionPending(ALARM_LIST_ACTIONS_TYPES.fastCancel, false);
      }
    },

    showEventModal(name, alarms, action, config = {}) {
      const { prepareEvent = event => event, ...modalConfig } = config;

      this.$modals.show({
        name,
        config: {
          items: alarms,
          ...modalConfig,
          action: async (event) => {
            const preparedEvent = prepareEvent(event);

            await action({ data: alarms.map(alarm => ({ ...preparedEvent, _id: alarm._id })) });
            await this.afterSubmit();
          },
        },
      });
    },

    showVariablesHelperModal() {
      const { entity, pbehavior, infos, ...alarmOtherFields } = this.item;
      const variables = [{
        ...convertObjectToTreeview(alarmOtherFields, 'alarm'),
        original: this.item,
      }];

      if (entity) {
        variables.push(convertObjectToTreeview(entity, 'entity'));
      }

      if (pbehavior) {
        variables.push(convertObjectToTreeview(pbehavior, 'pbehavior'));
      }

      this.$modals.show({
        name: MODALS.variablesHelp,
        config: {
          exportEntity: this.item,
          exportEntityName: `${ALARM_EXPORT_FILE_NAME_PREFIX}-${this.item._id}`,
          exportPdfTemplate: this.widget.parameters.exportPdfTemplate,
          variables,
        },
      });
    },

    showAddPbehaviorModal() {
      this.$modals.show({
        name: MODALS.pbehaviorPlanning,
        config: {
          entityPattern: createEntityIdPatternByValue(this.item.entity?._id),
          entities: mapAlarmsEntities([this.item]),
          afterSubmit: async () => {
            this.$popups.success({ text: this.$t('modals.createPbehavior.success.create') });

            await promisedWait(ALARM_LIST_FAST_PBEHAVIOR_TIMEOUT);

            return this.afterSubmit();
          },
        },
      });
    },

    async fastAddPbehavior(pbehaviorParameters = {}) {
      this.setActionPending(ALARM_LIST_ACTIONS_TYPES.fastPbehaviorAdd, true);

      try {
        const { name_prefix: namePrefix = '' } = pbehaviorParameters ?? {};
        let { type, reason } = pbehaviorParameters ?? {};

        if (!type) {
          const { data: defaultPbehaviorTypes = [] } = await this.fetchDefaultPbehaviorTypes();
          const pauseType = defaultPbehaviorTypes.find(
            ({ type: pbehaviorType }) => pbehaviorType === PBEHAVIOR_TYPE_TYPES.pause,
          );

          type = pauseType?._id;
        }

        if (!reason) {
          const { data: reasons = [] } = await this.fetchPbehaviorReasonsListWithoutStore();

          reason = reasons?.[0]?._id;
        }

        await this.bulkCreateEntityPbehaviors({
          data: mapAlarmsEntities([this.item]).map(entity => createDowntimePbehavior({
              entity,
              prefix: namePrefix,
              type,
              reason,
              origin: PBEHAVIOR_ORIGINS.alarmList,
            }),
          ),
        });

        await promisedWait(ALARM_LIST_FAST_PBEHAVIOR_TIMEOUT);
        await this.afterSubmit();
      } catch (err) {
        console.error(err);
      } finally {
        this.setActionPending(ALARM_LIST_ACTIONS_TYPES.fastPbehaviorAdd, false);
      }
    },

    async fastRemovePbehavior() {
      this.setActionPending(ALARM_LIST_ACTIONS_TYPES.fastPbehaviorAdd, true);

      try {
        await this.bulkRemoveEntityPbehaviors({
          data: mapAlarmsEntities([this.item]).map(({ _id: id }) => ({
            origin: PBEHAVIOR_ORIGINS.alarmList,
            entity: id,
          })),
        });

        await promisedWait(ALARM_LIST_FAST_PBEHAVIOR_TIMEOUT);
        await this.afterSubmit();
      } catch (err) {
        console.error(err);
      } finally {
        this.setActionPending(ALARM_LIST_ACTIONS_TYPES.fastPbehaviorAdd, false);
      }
    },

    showHistoryModal() {
      const widget = generatePreparedDefaultAlarmListWidget();

      widget.parameters.widgetColumns = this.widget.parameters.widgetColumns;
      widget.parameters.widgetGroupColumns = this.widget.parameters.widgetGroupColumns;
      widget.parameters.serviceDependenciesColumns = this.widget.parameters.serviceDependenciesColumns;

      this.$modals.show({
        name: MODALS.alarmsList,
        config: {
          widget,
          title: this.$t('modals.alarmsList.prefixTitle', { prefix: this.item.entity._id }),
          fetchList: params => this.fetchResolvedAlarmsListWithoutStore({
              params: { ...params, _id: this.item.entity._id },
            }),
        },
      });
    },

    showRemoveAlarmsFromManualMetaAlarmModal() {
      this.showRemoveAlarmsFromMetaAlarmModal(
        this.$t('alarm.actions.titles.removeAlarmsFromManualMetaAlarm'),
      );
    },

    showRemoveAlarmsFromAutoMetaAlarmModal() {
      this.showRemoveAlarmsFromMetaAlarmModal(
        this.$t('alarm.actions.titles.removeAlarmsFromAutoMetaAlarm'),
      );
    },

    showRemoveAlarmsFromMetaAlarmModal(title) {
      this.$modals.show({
        name: MODALS.removeAlarmsFromMetaAlarm,
        config: {
          items: [this.item],
          title,
          isCommentRequired: this.widget.parameters.isRemoveAlarmsFromMetaAlarmCommentRequired,
          action: async (data) => {
            await this.removeAlarmsFromMetaAlarm({
              id: this.parentAlarm?._id,
              data,
            });

            await this.afterSubmit();
          },
        },
      });
    },

    async handleLinkClickAction(link, type) {
      try {
        this.setActionPending(type, true);

        const links = await this.fetchAlarmLinksWithoutStore({
          id: link.rule_id,
          params: { ids: mapIds([this.item]) },
        });

        const summaryLink = links.find(({ icon_name: iconName, label = '' } = {}) => iconName === link.icon_name && label.includes(link.label),
        );

        if (!summaryLink) {
          return;
        }

        if (link.action === LINK_RULE_ACTIONS.copy) {
          await navigator.clipboard?.writeText(summaryLink.url);

          return;
        }

        window.open(summaryLink.url, '_blank');
      } catch (err) {
        console.error(err);
      } finally {
        this.setActionPending(type, false);
      }
    },

    async addBookmark() {
      await this.addBookmarkByAlarm(this.item);
    },

    async removeBookmark() {
      await this.removeBookmarkByAlarm(this.item);
    },

    async addBookmarkByAlarm(alarm) {
      try {
        this.setActionPending(ALARM_LIST_ACTIONS_TYPES.addBookmark, true);
        await this.addBookmarkToAlarm({ id: alarm._id });
        this.$popups.success({ text: this.$t('alarm.popups.addBookmarkSuccess') });
        await this.afterSubmit();
      } catch (err) {
        console.error(err);
        this.$popups.error({ text: this.$t('alarm.popups.addBookmarkFailed') });
      } finally {
        this.setActionPending(ALARM_LIST_ACTIONS_TYPES.addBookmark, false);
      }
    },

    async removeBookmarkByAlarm(alarm) {
      try {
        this.setActionPending(ALARM_LIST_ACTIONS_TYPES.addBookmark, true);
        await this.removeBookmarkFromAlarm({ id: alarm._id });
        this.$popups.success({ text: this.$t('alarm.popups.removeBookmarkSuccess') });
        await this.afterSubmit();
      } catch (err) {
        console.error(err);
        this.$popups.error({ text: this.$t('alarm.popups.removeBookmarkFailed') });
      } finally {
        this.setActionPending(ALARM_LIST_ACTIONS_TYPES.addBookmark, false);
      }
    },

    actionsAccessFilterHandler({ type }) {
      const permission = BUSINESS_USER_PERMISSIONS_ACTIONS_MAP.alarmsList[type];

      if (!permission) {
        return true;
      }

      return this.checkAccess(permission);
    },

    showExecuteInstructionModal(assignedInstruction) {
      const refreshAlarm = () => this.refreshAlarmsList();

      this.$modals.show({
        id: `${this.item._id}${assignedInstruction._id}`,
        name: isInstructionTypeManual(assignedInstruction?.type)
          ? MODALS.executeRemediationInstruction
          : MODALS.executeRemediationSimpleInstruction,
        config: {
          assignedInstruction,
          alarmId: this.item._id,
          onClose: refreshAlarm,
          onComplete: refreshAlarm,
        },
      });
    },
  },
};
</script>
