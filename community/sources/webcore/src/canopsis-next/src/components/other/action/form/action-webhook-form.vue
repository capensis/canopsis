<template>
  <v-layout column>
    <request-form
      v-field="webhook.request"
      :name="`${name}.request`"
      :headers-variables="payloadVariables"
      :payload-variables="payloadVariables"
      :url-variables="payloadVariables"
    />
    <declare-ticket-rule-ticket-mapping-field v-field="webhook" with-ticket-system-name />
  </v-layout>
</template>

<script>
import { formMixin } from '@/mixins/form';
import { payloadVariablesMixin } from '@/mixins/payload/variables';

import RequestForm from '@/components/forms/request/request-form.vue';
import DeclareTicketRuleTicketMappingField from '@/components/other/declare-ticket/form/fields/declare-ticket-rule-ticket-mapping-field.vue';

export default {
  inject: ['$validator'],
  components: { DeclareTicketRuleTicketMappingField, RequestForm },
  mixins: [
    formMixin,
    payloadVariablesMixin,
  ],
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
  computed: {
    payloadVariables() {
      const variables = [
        ...this.alarmPayloadVariables,
      ];

      if (this.hasPrevious) {
        variables.push(...this.payloadVariablesFromPreviousStep);
      }

      variables.push(...this.additionalDataVariables);

      return variables;
    },
  },
};
</script>
