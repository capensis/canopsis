<template lang="pug">
  v-layout(column)
    request-form(v-field="webhook.request", :name="name")
    declare-ticket-rule-ticket-mapping-field(v-field="webhook.declare_ticket")
</template>

<script>
import { formMixin } from '@/mixins/form';

import RequestForm from '@/components/forms/request/request-form.vue';
import DeclareTicketRuleTicketMappingField from '@/components/other/declare-ticket/form/fields/declare-ticket-rule-ticket-mapping-field.vue';

export default {
  inject: ['$validator'],
  components: { DeclareTicketRuleTicketMappingField, RequestForm },
  mixins: [formMixin],
  model: {
    prop: 'webhook',
    event: 'input',
  },
  props: {
    webhook: {
      type: Object,
      required: true,
    },
    name: {
      type: String,
      required: true,
    },
  },
};
</script>
