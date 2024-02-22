<template>
  <v-layout
    class="scenario-action-card-general-tab"
    wrap
  >
    <v-flex
      v-for="(item, index) in items"
      :key="index"
      :class="item.flexClass || 'xs12'"
      class="mt-1"
    >
      <scenario-info-item
        :icon="item.icon"
        :label="item.label"
        :value="item.value"
      />
      <component
        v-if="item.subcomponent"
        :is="item.subcomponent"
        v-bind="item.subcomponentProps"
      />
    </v-flex>
  </v-layout>
</template>

<script>
import { rrulestr } from 'rrule';

import { ACTION_TYPES } from '@/constants';

import { convertDurationToString } from '@/helpers/date/duration';

import ScenarioInfoItem from './scenario-info-item.vue';
import ScenarioActionCardCompiledTemplate from './scenario-action-card-compiled-template.vue';

export default {
  components: { ScenarioInfoItem, ScenarioActionCardCompiledTemplate },
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
      const getItems = {
        [ACTION_TYPES.ack]: this.getItemsDefault,
        [ACTION_TYPES.ackremove]: this.getItemsDefault,
        [ACTION_TYPES.cancel]: this.getItemsDefault,
        [ACTION_TYPES.assocticket]: this.getItemsForAssocTicketType,
        [ACTION_TYPES.changeState]: this.getItemsForChangeStateType,
        [ACTION_TYPES.snooze]: this.getItemsForSnoozeType,
        [ACTION_TYPES.pbehavior]: this.getItemsForPbehaviorType,
        [ACTION_TYPES.webhook]: this.getItemsForWebhookType,
      }[this.action.type];

      const items = getItems();

      const dropScenarioItemField = this.action.drop_scenario_if_not_matched
        ? { icon: 'stop', label: this.$t('common.stop') }
        : { icon: 'trending_flat', label: this.$t('scenario.remainingAction') };

      items.push(
        {
          icon: 'message',
          label: this.$tc('common.comment'),
          subcomponent: 'scenario-action-card-compiled-template',
          subcomponentProps: {
            template: this.action.comment,
          },
        },
        {
          icon: 'bolt',
          label: this.$t('common.emitTrigger'),
          flexClass: 'xs6',
        },
        {
          ...dropScenarioItemField,
          flexClass: 'xs6',
        },
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
     * Get items for rendering for `assocticket` action type
     *
     * @returns {Object[]}
     */
    getItemsForAssocTicketType() {
      const items = [
        this.outputItem,
        {
          icon: 'assignment',
          label: this.$tc('declareTicket.ticketID'),
          value: this.action.parameters.ticket,
        },
      ];

      if (this.action.parameters.ticket_url) {
        items.push({
          icon: 'assignment',
          label: this.$tc('declareTicket.ticketURL'),
          value: this.action.parameters.ticket_url,
        });
      }

      if (this.action.parameters.ticket_system_name) {
        items.push({
          icon: 'assignment',
          label: this.$tc('common.systemName'),
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
      const { filters } = this.$options;
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
            flexClass: 'xs6',
          },
          {
            icon: 'av_timer',
            label: this.$t('common.duration'),
            value: filters.duration(parameters.duration),
            flexClass: 'xs6',
          },
        );
      } else {
        result.push(
          {
            icon: 'alarm_on',
            label: this.$t('common.start'),
            value: filters.date(parameters.tstart),
            flexClass: 'xs6',
          },
          {
            icon: 'alarm_off',
            label: this.$t('common.stop'),
            value: filters.date(parameters.tstop),
            flexClass: 'xs6',
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
          flexClass: 'xs6',
        },
        {
          icon: 'language',
          label: this.$t('common.url'),
          value: request.url,
          flexClass: 'xs6',
        },
      ];

      if (request.auth) {
        result.push(
          {
            icon: 'person',
            label: this.$t('common.username'),
            value: request.auth.username,
            flexClass: 'xs6',
          },
          {
            icon: 'lock',
            label: this.$t('common.password'),
            value: request.auth.password,
            flexClass: 'xs6',
          },
        );
      }

      if (request.timeout) {
        result.push(
          {
            icon: 'av_timer',
            label: this.$t('common.request.timeout'),
            value: convertDurationToString(request.timeout),
          },
        );
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
            flexClass: 'xs6',
          },
          {
            label: this.$t('common.retryDelay'),
            value: convertDurationToString(request.retry_delay),
            flexClass: 'xs6',
          },
        );
      }

      if (headersArray.length) {
        const headersItems = headersArray.map(([label, value]) => ({ label, value }));

        result.push(
          { icon: 'list', label: this.$tc('common.header', 2) },

          ...headersItems,
        );
      }

      if (declareTicketArray.length) {
        const declareTicketItems = declareTicketArray.map(([label, value]) => ({ label, value: String(value) }));

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
            value: request.payload,
            rows: 3,
            autoGrow: true,
            box: true,
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
.scenario-action-card-general-tab ::v-deep {
  .v-input, .compiled-template__wrapper {
    margin-left: 16px;

    textarea {
      margin: 0;
    }
  }
}
</style>
