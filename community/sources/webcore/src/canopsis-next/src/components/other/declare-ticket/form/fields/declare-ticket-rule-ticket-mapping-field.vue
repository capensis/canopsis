<template>
  <v-layout column>
    <c-information-block
      :title="$t('declareTicket.ticketUrlAndId')"
      :help-text="ticketUrlHelpText"
      help-icon="help"
      help-icon-color="grey darken-1"
    >
      <c-alert
        v-if="isDeclareTicketExist"
        type="info"
      >
        {{ $t('declareTicket.webhookTicketDeclarationExist') }}
      </c-alert>
      <v-layout>
        <v-flex xs6>
          <c-enabled-field
            v-field="form.declare_ticket.enabled"
            :disabled="isDeclareTicketExist"
          />
        </v-flex>
        <v-flex
          v-if="form.declare_ticket.enabled"
          xs6
        >
          <c-enabled-field
            v-field="form.declare_ticket.is_regexp"
            :label="$t('declareTicket.isRegexp')"
          />
        </v-flex>
      </v-layout>
      <template v-if="form.declare_ticket.enabled">
        <c-enabled-field
          v-if="!hideEmptyResponse"
          v-field="form.declare_ticket.empty_response"
          :label="$t('declareTicket.emptyResponse')"
        />
        <declare-ticket-rule-ticket-id-field
          v-field="form.declare_ticket.ticket_id"
          :disabled="disabled"
          :name="ticketIdFieldName"
          :required="ticketIdRequired"
          :variables="variables"
        />
        <declare-ticket-rule-ticket-url-field
          v-field="form.declare_ticket.ticket_url"
          :disabled="disabled"
          :name="ticketUrlFieldName"
          :variables="variables"
        />
        <v-flex offset-xs6>
          <declare-ticket-rule-ticket-url-title-field v-field="form.declare_ticket.ticket_url_title" />
          <v-text-field
            v-if="withTicketSystemName"
            v-field="form.ticket_system_name"
            :label="$t('declareTicket.ticketSystemName')"
          />
        </v-flex>
        <declare-ticket-rule-ticket-custom-fields-field
          v-field="form.declare_ticket.mapping"
          :name="name"
          :disabled="disabled"
        />
      </template>
    </c-information-block>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

import DeclareTicketRuleTicketIdField from './declare-ticket-rule-ticket-id-field.vue';
import DeclareTicketRuleTicketCustomFieldsField from './declare-ticket-rule-ticket-custom-fields-field.vue';
import DeclareTicketRuleTicketUrlField from './declare-ticket-rule-ticket-url-field.vue';
import DeclareTicketRuleTicketUrlTitleField from './declare-ticket-rule-ticket-url-title-field.vue';

export default {
  inject: ['$validator'],
  components: {
    DeclareTicketRuleTicketUrlField,
    DeclareTicketRuleTicketCustomFieldsField,
    DeclareTicketRuleTicketIdField,
    DeclareTicketRuleTicketUrlTitleField,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    name: {
      type: String,
      default: 'declare_ticket',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    isDeclareTicketExist: {
      type: Boolean,
      default: false,
    },
    hideEmptyResponse: {
      type: Boolean,
      default: false,
    },
    ticketIdRequired: {
      type: Boolean,
      default: false,
    },
    onlyOneTicketId: {
      type: Boolean,
      default: false,
    },
    withTicketSystemName: {
      type: Boolean,
      default: false,
    },
    variables: {
      type: Array,
      default: () => [],
    },
  },
  setup(props) {
    const { t } = useI18n();

    const ticketIdFieldName = computed(() => `${props.name}.ticket_id`);

    const ticketUrlFieldName = computed(() => `${props.name}.ticket_url`);

    const ticketUrlHelpText = computed(() => [
      t('declareTicket.ticketUrlAndIdHelpText'),
      props.onlyOneTicketId && t('declareTicket.dataFromOneStepAttention'),
    ].filter(Boolean).join('\n'));

    return {
      ticketIdFieldName,
      ticketUrlFieldName,
      ticketUrlHelpText,
    };
  },
};
</script>
