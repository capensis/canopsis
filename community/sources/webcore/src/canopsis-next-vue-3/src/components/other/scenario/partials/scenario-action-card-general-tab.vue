<template>
  <c-layout class="scenario-action-card-general-tab" wrap>
    <div
      v-for="(item, index) in items"
      :key="index"
      :class="item.flexClass || 'w-100'"
      class="mt-1"
    >
      <scenario-info-item-multiline
        v-if="item.multiline"
        :icon="item.icon"
        :label="item.label"
        :value="item.value"
      />
      <scenario-info-item
        v-else
        :icon="item.icon"
        :label="item.label"
        :value="item.value"
      />
      <component v-if="item.subcomponent" :is="item.subcomponent" v-bind="item.subcomponentProps" />
    </div>
  </c-layout>
</template>

<script>
import { rrulestr } from 'rrule';

import { ACTION_TYPES } from '@/constants';

import { convertDateToString } from '@/helpers/date/date';
import { convertDurationToString } from '@/helpers/date/duration';

import ScenarioActionCardCompiledTemplate from './scenario-action-card-compiled-template.vue';
import ScenarioInfoItemMultiline from './scenario-info-item-multiline.vue';
import ScenarioInfoItem from './scenario-info-item.vue';

export default {
  components: {
    ScenarioInfoItem,
    ScenarioInfoItemMultiline,
    ScenarioActionCardCompiledTemplate,
  },

  props: {
    action: {
      type: Object,
      required: true,
    },
  },

  computed: {
    outputItem() {
      return {
        icon: 'message',
        label: this.$t('scenario.output'),
        subcomponent: 'scenario-action-card-compiled-template',
        subcomponentProps: {
          template: this.action.parameters.output,
        },
      };
    },

    items() {
      const getItems =
        {
          [ACTION_TYPES.ack]: this.getItemsDefault,
          [ACTION_TYPES.ackremove]: this.getItemsDefault,
          [ACTION_TYPES.cancel]: this.getItemsDefault,
          [ACTION_TYPES.assocticket]: this.getItemsForAssocTicketType,
          [ACTION_TYPES.changeState]: this.getItemsForChangeStateType,
          [ACTION_TYPES.snooze]: this.getItemsForSnoozeType,
          [ACTION_TYPES.pbehavior]: this.getItemsForPbehaviorType,
          [ACTION_TYPES.webhook]: this.getItemsForWebhookType,
          [ACTION_TYPES.pause]: this.getItemsForPauseType,
        }[this.action.type] ?? this.getItemsDefault;

      const items = getItems.call(this);

      const dropScenarioItemField = this.action.drop_scenario_if_not_matched
        ? {
            icon: 'stop',
            label: this.$t('scenario.actionDidNotMatch'),
            value: this.$t('common.stop'),
            multiline: true,
          }
        : {
            icon: 'trending_flat',
            label: this.$t('scenario.actionDidNotMatch'),
            value: this.$t('scenario.remainingAction'),
            multiline: true,
          };

      const workflowFields = [{ ...dropScenarioItemField, flexClass: 'w-50' }];

      if (this.action.type === ACTION_TYPES.webhook) {
        const { parameters = {} } = this.action;
        const { stop_on_fail: stopOnFail, stop_on_success: stopOnSuccess } = parameters;

        workflowFields.push(
          {
            icon: stopOnFail ? 'stop' : 'trending_flat',
            value: stopOnFail ? this.$t('common.stop') : this.$t('scenario.remainingStep'),
            label: this.$t('scenario.caseFailure'),
            flexClass: 'w-50 ml-auto',
            multiline: true,
          },
          {
            icon: stopOnSuccess ? 'stop' : 'trending_flat',
            value: stopOnSuccess ? this.$t('common.stop') : this.$t('scenario.remainingStep'),
            label: this.$t('scenario.caseSuccess'),
            flexClass: 'w-50 ml-auto',
            multiline: true,
          },
        );
      }

      const triggerFields = [
        {
          icon: 'bolt',
          label: this.$t('common.emitTriggerSuccess'),
          flexClass: 'w-50',
        },
      ];

      if (this.action.type === ACTION_TYPES.webhook && this.action.parameters?.declare_ticket) {
        triggerFields.push({
          icon: 'bolt',
          label: this.$t('common.emitTriggerFail'),
          flexClass: 'w-50',
        });
      }

      items.push(
        {
          icon: 'message',
          label: this.$t('common.comment'),
          subcomponent: 'scenario-action-card-compiled-template',
          subcomponentProps: {
            template: this.action.comment,
          },
        },
        ...triggerFields,
        ...workflowFields,
      );

      return items;
    },
  },

  methods: {
    /**
     * Get items for rendering for action types which have only `output` parameter
     *
     * @returns {Object[]}
     */
    getItemsDefault() {
      return [this.outputItem];
    },

    /**
     * Get items for an entity pause action.
     *
     * @returns {Object[]}
     */
    getItemsForPauseType() {
      return [
        this.outputItem,
        {
          icon: 'label',
          label: this.$t('common.type'),
          value: this.action.parameters.type,
        },
        {
          icon: 'account_tree',
          label: this.$t('settings.pause.propagate'),
          value: this.$t(this.action.parameters.inherited ? 'common.yes' : 'common.no'),
        },
      ];
    },

    /**
     * Get items for rendering for `assocticket` action type
     *
     * @returns {Object[]}
     */
    getItemsForAssocTicketType() {
      const items = [
        this.outputItem,
        {
          icon: 'assignment',
          label: this.$t('declareTicket.ticketID'),
          value: this.action.parameters.ticket,
        },
      ];

      if (this.action.parameters.ticket_url) {
        items.push({
          icon: 'assignment',
          label: this.$t('declareTicket.ticketURL'),
          value: this.action.parameters.ticket_url,
        });
      }

      if (this.action.parameters.ticket_url_title) {
        items.push({
          icon: 'assignment',
          label: this.$t('declareTicket.ticketURLTitle'),
          value: this.action.parameters.ticket_url_title,
        });
      }

      if (this.action.parameters.ticket_system_name) {
        items.push({
          icon: 'assignment',
          label: this.$t('common.systemName'),
          value: this.action.parameters.ticket_system_name,
        });
      }

      return items;
    },

    /**
     * Get items for rendering for `changeState` action type
     *
     * @returns {Object[]}
     */
    getItemsForChangeStateType() {
      return [
        this.outputItem,

        {
          icon: 'assignment',
          label: this.$t('common.state'),
          value: this.$t(`common.stateTypes.${this.action.parameters.state}`),
        },
      ];
    },

    /**
     * Get items for rendering for `snooze` action type
     *
     * @returns {Object[]}
     */
    getItemsForSnoozeType() {
      return [
        this.outputItem,

        {
          icon: 'alarm_off',
          label: this.$t('common.duration'),
          value: convertDurationToString(this.action.parameters.duration),
        },
      ];
    },

    /**
     * Get items for rendering for `pbehavior` action type
     *
     * @returns {Object[]}
     */
    getItemsForPbehaviorType() {
      const { parameters } = this.action;
      const result = [
        {
          icon: 'short_text',
          label: this.$t('common.name'),
          value: parameters.name,
        },
      ];

      if (parameters.start_on_trigger) {
        result.push(
          {
            icon: 'alarm_on',
            label: this.$t('modals.createPbehavior.steps.general.fields.startOnTrigger'),
            flexClass: 'w-50',
          },
          {
            icon: 'av_timer',
            label: this.$t('common.duration'),
            value: convertDurationToString(parameters.duration),
            flexClass: 'w-50',
          },
        );
      } else {
        result.push(
          {
            icon: 'alarm_on',
            label: this.$t('common.start'),
            value: convertDateToString(parameters.tstart),
            flexClass: 'w-50',
          },
          {
            icon: 'alarm_off',
            label: this.$t('common.stop'),
            value: convertDateToString(parameters.tstop),
            flexClass: 'w-50',
          },
        );
      }

      result.push(
        {
          icon: parameters.type.icon_name,
          label: this.$t('common.type'),
          value: parameters.type.name,
        },
        {
          icon: 'assignment',
          label: this.$t('common.reason'),
          value: parameters.reason.name,
        },
      );

      if (parameters.rrule) {
        const recurrenceRule = rrulestr(parameters.rrule);

        result.push({
          icon: 'calendar_today',
          label: this.$t('pbehavior.rrule'),
          value: recurrenceRule.toText(),
        });
      }

      return result;
    },

    /**
     * Get items for rendering for `webhook` action type
     *
     * @returns {Object[]}
     */
    getItemsForWebhookType() {
      const { parameters = {} } = this.action;
      const { request = {}, declare_ticket: declareTicket = {} } = parameters;
      const headersArray = request.headers ? Object.entries(request.headers) : [];
      const declareTicketArray = Object.entries(declareTicket);

      const result = [
        {
          icon: 'short_text',
          label: this.$t('common.method'),
          value: request.method,
          flexClass: 'w-50',
        },
        {
          icon: 'language',
          label: this.$t('common.url'),
          value: request.url,
          flexClass: 'w-50',
        },
      ];

      if (parameters.multiple_urls) {
        result.push({
          icon: 'link',
          label: this.$t('scenario.allowMultipleUrls'),
          flexClass: 'w-50 ml-auto',
        });
      }

      if (request.auth?.username) {
        result.push({
          icon: 'person',
          label: this.$t('common.username'),
          value: request.auth.username,
          flexClass: 'w-50',
        });
      }

      if (request.timeout) {
        result.push({
          icon: 'av_timer',
          label: this.$t('common.request.timeout'),
          value: convertDurationToString(request.timeout),
        });
      }

      if (request.retry_count) {
        result.push(
          {
            icon: 'av_timer',
            label: this.$t('common.request.repeatRequest'),
          },
          {
            label: this.$t('common.retryCount'),
            value: request.retry_count,
            flexClass: 'w-50',
          },
          {
            label: this.$t('common.retryDelay'),
            value: convertDurationToString(request.retry_delay),
            flexClass: 'w-50',
          },
        );
      }

      if (headersArray.length) {
        const headersItems = headersArray.map(([label, value]) => ({ label, value }));

        result.push(
          { icon: 'list', label: this.$t('common.header', 2) },

          ...headersItems,
        );
      }

      if (declareTicketArray.length) {
        const declareTicketItems = declareTicketArray.map(([label, value]) => ({
          label,
          value: String(value),
        }));

        result.push(
          { icon: 'list', label: this.$t('scenario.declareTicket') },

          ...declareTicketItems,
        );
      }

      if (request.payload) {
        result.push({
          icon: 'assignment',
          label: this.$t('common.payload'),
          subcomponent: 'v-textarea',
          subcomponentProps: {
            modelValue: request.payload,
            rows: 3,
            autoGrow: true,
            readonly: true,
            disabled: true,
          },
        });
      }

      return result;
    },
  },
};
</script>

<style lang="scss" scoped>
.scenario-action-card-general-tab {
  .v-input,
  .compiled-template__wrapper {
    margin-left: 16px;

    textarea {
      margin: 0;
    }
  }
}
</style>
