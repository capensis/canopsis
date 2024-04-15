<template>
  <v-layout>
    <v-flex xs6>
      <c-enabled-field
        v-field="form.template"
        :label="$t('declareTicket.allowTicketURLTemplate')"
      />
    </v-flex>
    <v-flex xs6>
      <declare-ticket-rule-ticket-url-payload-text-field
        v-if="form.template"
        v-field="form.value"
        :name="name"
        :disabled="disabled"
        :variables="variables"
      />
      <declare-ticket-rule-ticket-url-text-field
        v-else
        v-field="form.value"
        :name="name"
        :disabled="disabled"
      />
    </v-flex>
  </v-layout>
</template>

<script>
import DeclareTicketRuleTicketUrlTextField from './declare-ticket-rule-ticket-url-text-field.vue';
import DeclareTicketRuleTicketUrlPayloadTextField from './declare-ticket-rule-ticket-url-payload-text-field.vue';

export default {
  components: {
    DeclareTicketRuleTicketUrlTextField,
    DeclareTicketRuleTicketUrlPayloadTextField,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'ticket_url',
    },
    disabled: {
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
