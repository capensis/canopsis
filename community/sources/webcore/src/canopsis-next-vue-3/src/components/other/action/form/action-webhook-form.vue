<template>
  <div>
    <request-with-token-form
      :form="modelValue"
      :name="requestFormName"
      :url-variables="templateVars[webhookTemplateVarsKey]"
      :headers-variables="templateVars[webhookTemplateVarsKey]"
      :payload-variables="templateVars[webhookTemplateVarsKey]"
      :depth="depth"
      with-multiple-urls
      @update:form="$emit('update:modelValue', $event)"
    />

    <declare-ticket-rule-ticket-mapping-field
      :form="modelValue"
      :name="`${name}.declare_ticket`"
      :depth="depth"
      :variables="templateVars.ticket"
      hide-empty-response
      ticket-id-required
      only-one-ticket-id
      with-ticket-system-name
      @update:form="$emit('update:modelValue', $event)"
    />

    <declare-ticket-rule-check-ticket-status-field
      v-if="modelValue.declare_ticket?.enabled"
      :form="modelValue.declare_ticket.check_ticket_status"
      :name="`${name}.declare_ticket.check_ticket_status`"
      :template-vars="templateVars"
      :depth="depth"
      top-border
      @update:form="updateField('declare_ticket.check_ticket_status', $event)"
    />
  </div>
</template>

<script>
import { computed } from 'vue';

import { useModelField } from '@/hooks/form/model-field';

import RequestWithTokenForm from '@/components/forms/request/request-with-token-form.vue';
import DeclareTicketRuleCheckTicketStatusField from '@/components/other/declare-ticket/form/fields/declare-ticket-rule-check-ticket-status-field.vue';
import DeclareTicketRuleTicketMappingField from '@/components/other/declare-ticket/form/fields/declare-ticket-rule-ticket-mapping-field.vue';

export default {
  components: {
    RequestWithTokenForm,
    DeclareTicketRuleTicketMappingField,
    DeclareTicketRuleCheckTicketStatusField,
  },

  props: {
    modelValue: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'parameters',
    },
    hasPrevious: {
      type: Boolean,
      default: false,
    },
    depth: {
      type: Number,
      default: 0,
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { updateField } = useModelField(props, emit);
    const requestFormName = computed(() => `${props.name}.request`);
    const webhookTemplateVarsKey = computed(() =>
      props.hasPrevious ? 'webhook' : 'first_webhook',
    );

    return {
      requestFormName,
      webhookTemplateVarsKey,
      updateField,
    };
  },
};
</script>
