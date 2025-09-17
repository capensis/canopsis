<template>
  <v-layout column>
    <request-form
      v-field="webhook.request"
      :name="`${name}.request`"
      :url-variables="templateVars[webhookTemplateVarsKey]"
      :headers-variables="templateVars[webhookTemplateVarsKey]"
      :payload-variables="templateVars[webhookTemplateVarsKey]"
    />
    <declare-ticket-rule-ticket-mapping-field
      v-field="webhook"
      :template-vars="templateVars.ticket"
      with-ticket-system-name
    />
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import RequestForm from '@/components/forms/request/request-form.vue';
import DeclareTicketRuleTicketMappingField from '@/components/other/declare-ticket/form/fields/declare-ticket-rule-ticket-mapping-field.vue';

export default {
  inject: ['$validator'],
  components: { DeclareTicketRuleTicketMappingField, RequestForm },
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
    hasPreviousWebhook: {
      type: Boolean,
      default: false,
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props) {
    const webhookTemplateVarsKey = computed(() => (props.hasPreviousWebhook ? 'webhook' : 'first_webhook'));

    return {
      webhookTemplateVarsKey,
    };
  },
};
</script>
