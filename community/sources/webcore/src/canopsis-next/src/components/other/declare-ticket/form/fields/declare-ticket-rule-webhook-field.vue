<template lang="pug">
  v-layout(column)
    request-form(v-field="value.request", :name="`${name}.request`")
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
  },
};
</script>
