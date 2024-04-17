<template>
  <v-layout>
    <v-flex xs6>
      <c-enabled-field
        v-field="value.template"
        :label="$t('declareTicket.allowTicketIdTemplate')"
      />
    </v-flex>
    <declare-ticket-rule-ticket-id-payload-text-field
      v-if="value.template"
      v-field="value.value"
      :name="name"
      :disabled="disabled"
      :variables="variables"
      :required="required"
    />
    <declare-ticket-rule-ticket-id-text-field
      v-else
      v-field="value.value"
      :name="name"
      :disabled="disabled"
      :required="required"
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
