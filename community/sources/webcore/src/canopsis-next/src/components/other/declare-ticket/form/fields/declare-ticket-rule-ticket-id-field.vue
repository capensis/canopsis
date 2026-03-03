<template>
  <v-layout class="gap-3">
    <v-flex xs6>
      <declare-ticket-rule-ticket-id-payload-text-field
        v-if="value.ticket_id.template"
        v-field="value.ticket_id.value"
        :name="name"
        :disabled="disabled"
        :variables="variables"
        :required="required"
      />
      <declare-ticket-rule-ticket-id-text-field
        v-else
        v-field="value.ticket_id.value"
        :name="name"
        :disabled="disabled"
        :required="required"
      />
    </v-flex>
    <c-enabled-field
      v-field="value.ticket_id.template"
      :label="$t('declareTicket.allowTicketIdTemplate')"
    />
    <c-enabled-field
      v-field="value.is_regexp"
      :label="$t('declareTicket.isRegexp')"
      class="ml-4"
    />
  </v-layout>
</template>

<script>
import DeclareTicketRuleTicketIdPayloadTextField from './declare-ticket-rule-ticket-id-payload-text-field.vue';
import DeclareTicketRuleTicketIdTextField from './declare-ticket-rule-ticket-id-text-field.vue';

export default {
  inject: ['$validator'],
  components: {
    DeclareTicketRuleTicketIdTextField,
    DeclareTicketRuleTicketIdPayloadTextField,
  },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'ticket_id',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    variables: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
