<template lang="pug">
  c-information-block(:title="$t('declareTicket.ticketUrlAndId')")
    c-name-field(v-field="form.system_name", :label="$t('common.systemName')", name="system_name")
    v-layout(row)
      v-flex.mr-3(xs6)
        declare-ticket-rule-ticket-id-field(v-field="form.ticket_id", required)
      v-flex(xs6)
        declare-ticket-rule-ticket-url-field(v-field="form.ticket_url")
    declare-ticket-rule-ticket-custom-fields-field(v-field="form.mapping")
    c-description-field(v-field="form.output", :label="$tc('common.comment')")
</template>

<script>
import DeclareTicketRuleTicketIdField from './fields/declare-ticket-rule-ticket-id-field.vue';
import DeclareTicketRuleTicketUrlField from './fields/declare-ticket-rule-ticket-url-field.vue';
import DeclareTicketRuleTicketCustomFieldsField from './fields/declare-ticket-rule-ticket-custom-fields-field.vue';

export default {
  components: {
    DeclareTicketRuleTicketCustomFieldsField,
    DeclareTicketRuleTicketUrlField,
    DeclareTicketRuleTicketIdField,
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
  },
};
</script>
