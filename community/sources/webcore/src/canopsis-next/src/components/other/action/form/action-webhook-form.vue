<template>
  <v-layout class="gap-3" column>
    <request-with-token-form
      v-field="webhook"
      :name="`${name}.request`"
      :url-variables="templateVars[webhookTemplateVarsKey]"
      :headers-variables="templateVars[webhookTemplateVarsKey]"
      :payload-variables="templateVars[webhookTemplateVarsKey]"
      class="c-alternative-bg-panel pa-5"
      with-multiple-urls
    />
    <declare-ticket-rule-ticket-mapping-field
      v-field="webhook"
      :variables="templateVars.ticket"
      class="c-alternative-bg-panel pa-5"
      with-ticket-system-name
    />
    <declare-ticket-rule-check-ticket-status-field
      v-if="webhook.declare_ticket?.enabled"
      v-field="webhook.declare_ticket.check_ticket_status"
      :name="`${name}.declare_ticket.check_ticket_status`"
      :template-vars="templateVars"
      class="c-alternative-bg-panel pa-5"
    />
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { useModelField } from '@/hooks/form/model-field';

import RequestWithTokenForm from '@/components/forms/request/request-with-token-form.vue';
import DeclareTicketRuleTicketMappingField from '@/components/other/declare-ticket/form/fields/declare-ticket-rule-ticket-mapping-field.vue';
import DeclareTicketRuleCheckTicketStatusField from '@/components/other/declare-ticket/form/fields/declare-ticket-rule-check-ticket-status-field.vue';

export default {
  components: {
    RequestWithTokenForm,
    DeclareTicketRuleTicketMappingField,
    DeclareTicketRuleCheckTicketStatusField,
  },
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
    hasPrevious: {
      type: Boolean,
      default: false,
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props, { emit }) {
    const { updateField } = useModelField(props, emit);

    const webhookTemplateVarsKey = computed(() => (props.hasPrevious ? 'webhook' : 'first_webhook'));

    /**
     * Updates the multiple_urls field in the webhook form
     *
     * @param {boolean} multiple - Whether multiple URLs are enabled for the webhook
     */
    const updateMultiple = multiple => updateField('multiple_urls', multiple);

    return {
      webhookTemplateVarsKey,
      updateMultiple,
    };
  },
};
</script>
