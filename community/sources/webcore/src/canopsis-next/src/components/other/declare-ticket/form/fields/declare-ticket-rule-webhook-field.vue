<template lang="pug">
  v-layout(column)
    request-form(
      v-field="value.request",
      :name="`${name}.request`",
      :url-variables="payloadVariables",
      :headers-variables="payloadVariables",
      :payload-variables="payloadVariables"
    )
    declare-ticket-rule-ticket-mapping-field.mb-2(
      v-field="value.declare_ticket",
      :name="`${name}.declare_ticket`",
      :is-declare-ticket-exist="isDeclareTicketExist"
    )
    c-workflow-field(
      v-field="value.stop_on_fail",
      :disabled="disabled",
      :label="$t('declareTicket.workflowIfStepFails')",
      :continue-label="$t('declareTicket.continueWithNextStep')"
    )
</template>

<script>
import {
  ALARM_PAYLOADS_VARIABLES,
  ENTITY_PAYLOAD_VARIABLES,
  DECLARE_TICKET_PAYLOAD_ADDITIONAL_DATA_VARIABLES,
  DECLARE_TICKET_PAYLOAD_PREVIOUS_STEP_VARIABLES,
} from '@/constants';

import RequestForm from '@/components/forms/request/request-form.vue';

import DeclareTicketRuleTicketMappingField from './declare-ticket-rule-ticket-mapping-field.vue';

export default {
  components: { DeclareTicketRuleTicketMappingField, RequestForm },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
    name: {
      type: String,
      default: 'webhook',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    isDeclareTicketExist: {
      type: Boolean,
      default: false,
    },
    hasPrevious: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    alarmPayloadVariables() {
      return [{
        value: ALARM_PAYLOADS_VARIABLES.alarms,
        enumerable: true,
        variables: [
          {
            value: ALARM_PAYLOADS_VARIABLES.component,
            text: this.$t('common.component'),
          },
          {
            value: ALARM_PAYLOADS_VARIABLES.resource,
            text: this.$t('common.resource'),
          },
          {
            value: ALARM_PAYLOADS_VARIABLES.stateMessage,
            text: this.$t('alarm.stateMessage'),
          },
          {
            value: ALARM_PAYLOADS_VARIABLES.stateValue,
            text: this.$t('alarm.stateValue'),
          },
          {
            value: ALARM_PAYLOADS_VARIABLES.statusValue,
            text: this.$t('common.status'),
          },
          {
            value: ALARM_PAYLOADS_VARIABLES.ticketAuthor,
            text: this.$t('alarm.ticketAuthor'),
          },
          {
            value: ALARM_PAYLOADS_VARIABLES.ticketValue,
            text: this.$t('alarm.ticketId'),
          },
          {
            value: ALARM_PAYLOADS_VARIABLES.ticketMessage,
            text: this.$t('alarm.ticketMessage'),
          },
          {
            value: ALARM_PAYLOADS_VARIABLES.ackAuthor,
            text: this.$t('alarm.ackAuthor'),
          },
          {
            value: ALARM_PAYLOADS_VARIABLES.ackMessage,
            text: this.$t('alarm.ackMessage'),
          },
          {
            value: ALARM_PAYLOADS_VARIABLES.lastCommentAuthor,
            text: this.$t('alarm.lastCommentAuthor'),
          },
          {
            value: ALARM_PAYLOADS_VARIABLES.lastCommentMessage,
            text: this.$t('alarm.lastCommentMessage'),
          },
        ],
      }];
    },

    entityPayloadVariables() {
      return [{
        value: ENTITY_PAYLOAD_VARIABLES.entity,
        enumerable: true,
        variables: [{
          value: ENTITY_PAYLOAD_VARIABLES.infosValue,
          text: this.$t('common.infos'),
        }],
      }];
    },

    payloadVariablesFromPreviousStep() {
      return [
        {
          value: DECLARE_TICKET_PAYLOAD_PREVIOUS_STEP_VARIABLES.header,
          text: this.$t('declareTicket.headerFieldFromPreviousSteps'),
        },
        {
          value: DECLARE_TICKET_PAYLOAD_PREVIOUS_STEP_VARIABLES.response,
          text: this.$t('declareTicket.responseFieldFromPreviousSteps'),
        },
        {
          value: DECLARE_TICKET_PAYLOAD_PREVIOUS_STEP_VARIABLES.headerByStep,
          text: this.$t('declareTicket.headerFieldFromStep'),
        },
      ];
    },

    additionalDataVariables() {
      return [
        {
          value: DECLARE_TICKET_PAYLOAD_ADDITIONAL_DATA_VARIABLES.author,
          text: this.$t('common.username'),
        },
        {
          value: DECLARE_TICKET_PAYLOAD_ADDITIONAL_DATA_VARIABLES.user,
          text: this.$t('declareTicket.userId'),
        },
        {
          value: DECLARE_TICKET_PAYLOAD_ADDITIONAL_DATA_VARIABLES.alarmChangeType,
          text: this.$tc('common.trigger'),
        },
        {
          value: DECLARE_TICKET_PAYLOAD_ADDITIONAL_DATA_VARIABLES.initiator,
          text: this.$t('declareTicket.actionInitiator'),
        },
        {
          value: DECLARE_TICKET_PAYLOAD_ADDITIONAL_DATA_VARIABLES.output,
          text: this.$t('declareTicket.triggerEventMessage'),
        },
      ];
    },

    payloadVariables() {
      const variables = [...this.alarmPayloadVariables, ...this.entityPayloadVariables];

      if (this.hasPrevious) {
        variables.push(
          ...this.payloadVariablesFromPreviousStep,
          ...this.additionalDataVariables,
        );
      }

      return variables;
    },
  },
};
</script>
