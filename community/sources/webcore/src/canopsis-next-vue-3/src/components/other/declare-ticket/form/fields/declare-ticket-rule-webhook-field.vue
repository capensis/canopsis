<template>
  <c-card-iterator-item :item-number="webhookNumber" @remove="removeWebhook">
    <template #header="">
      <request-url-field
        :method="form.request.method"
        :url="form.request.url"
        :help-text="$t('common.request.urlHelp')"
        :name="requestFormName"
        :disabled="disabled"
        :url-variables="templateVars[webhookTemplateVarsKey]"
        @update:method="updateField('request.method', $event)"
        @update:url="updateField('request.url', $event)"
      />
    </template>
    <c-form-block class="mt-3">
      <request-with-token-form
        :form="form"
        :name="requestFormName"
        :headers-variables="templateVars[webhookTemplateVarsKey]"
        :payload-variables="templateVars[webhookTemplateVarsKey]"
        hide-url
        @update:form="$emit('update:form', $event)"
      />

      <declare-ticket-rule-ticket-mapping-field
        :form="form"
        :name="`${name}.declare_ticket`"
        :is-declare-ticket-exist="isDeclareTicketExist"
        :variables="templateVars.ticket"
        hide-empty-response
        ticket-id-required
        only-one-ticket-id
        @update:form="$emit('update:form', $event)"
      />

      <declare-ticket-rule-check-ticket-status-field
        v-if="form.declare_ticket.enabled"
        :form="form.declare_ticket.check_ticket_status"
        :name="`${name}.declare_ticket.check_ticket_status`"
        :template-vars="templateVars"
        :disabled="disabled"
        top-border
        @update:form="updateField('declare_ticket.check_ticket_status', $event)"
      />

      <c-form-block-row :label="$t('declareTicket.workflowIfStepFails')" indented top-border>
        <c-workflow-field
          :model-value="form.stop_on_fail"
          :disabled="disabled"
          :continue-label="$t('declareTicket.continueWithNextStep')"
          @update:model-value="updateField('stop_on_fail', $event)"
        />
      </c-form-block-row>
    </c-form-block>
  </c-card-iterator-item>
</template>

<script>
import { computed, toRef } from 'vue';

import { useConfirmableForm } from '@/hooks/confirmable-form';
import { useModelField } from '@/hooks/form/model-field';

import RequestUrlField from '@/components/forms/request/fields/request-url-field.vue';
import RequestWithTokenForm from '@/components/forms/request/request-with-token-form.vue';

import DeclareTicketRuleCheckTicketStatusField from './declare-ticket-rule-check-ticket-status-field.vue';
import DeclareTicketRuleTicketMappingField from './declare-ticket-rule-ticket-mapping-field.vue';

export default {
  components: {
    RequestUrlField,
    RequestWithTokenForm,
    DeclareTicketRuleTicketMappingField,
    DeclareTicketRuleCheckTicketStatusField,
  },
  props: {
    form: {
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
    webhookNumber: {
      type: Number,
      required: false,
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

  emits: ['update:form', 'remove'],

  setup(props, { emit }) {
    const { updateField } = useModelField(props, emit);
    const requestFormName = computed(() => `${props.name}.request`);
    const webhookTemplateVarsKey = computed(() =>
      props.hasPrevious ? 'webhook' : 'first_webhook',
    );

    const { confirmAction: removeWebhook } = useConfirmableForm({
      form: toRef(props, 'form'),
      action: () => emit('remove'),
    });

    return {
      requestFormName,
      webhookTemplateVarsKey,
      updateField,

      removeWebhook,
    };
  },
};
</script>
