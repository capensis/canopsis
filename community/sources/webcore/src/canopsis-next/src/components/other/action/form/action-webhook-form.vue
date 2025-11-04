<template>
  <v-layout column>
    <request-with-token-form
      v-field="webhook"
      :name="`${name}.request`"
      :url-variables="templateVars[webhookTemplateVarsKey]"
      :headers-variables="templateVars[webhookTemplateVarsKey]"
      :payload-variables="templateVars[webhookTemplateVarsKey]"
      with-multiple-urls
    />
    <declare-ticket-rule-ticket-mapping-field
      v-field="webhook"
      :variables="templateVars.ticket"
      with-ticket-system-name
    />
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { useModelField } from '@/hooks/form/model-field';

import RequestWithTokenForm from '@/components/forms/request/request-with-token-form.vue';
import DeclareTicketRuleTicketMappingField from '@/components/other/declare-ticket/form/fields/declare-ticket-rule-ticket-mapping-field.vue';

export default {
  components: { DeclareTicketRuleTicketMappingField, RequestWithTokenForm },
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
