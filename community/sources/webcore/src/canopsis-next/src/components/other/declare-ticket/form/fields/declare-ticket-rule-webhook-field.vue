<template>
  <c-card-iterator-item
    :item-number="webhookNumber"
    @remove="removeWebhook"
  >
    <template #header="">
      <request-url-field
        v-field="form.request"
        :help-text="$t('common.request.urlHelp')"
        :name="requestFormName"
        :disabled="disabled"
        :url-variables="templateVars[webhookTemplateVarsKey]"
      />
    </template>
    <request-with-token-form
      v-field="form"
      :name="requestFormName"
      :headers-variables="templateVars[webhookTemplateVarsKey]"
      :payload-variables="templateVars[webhookTemplateVarsKey]"
      hide-url
    />
    <declare-ticket-rule-ticket-mapping-field
      v-field="form"
      :name="`${name}.declare_ticket`"
      :is-declare-ticket-exist="isDeclareTicketExist"
      :variables="templateVars.ticket"
      class="mb-2"
      hide-empty-response
      ticket-id-required
      only-one-ticket-id
    />
    <c-workflow-field
      v-field="form.stop_on_fail"
      :disabled="disabled"
      :label="$t('declareTicket.workflowIfStepFails')"
      :continue-label="$t('declareTicket.continueWithNextStep')"
    />
  </c-card-iterator-item>
</template>

<script>
import { computed, toRef } from 'vue';

import { useConfirmableForm } from '@/hooks/confirmable-form';

import RequestWithTokenForm from '@/components/forms/request/request-with-token-form.vue';
import RequestUrlField from '@/components/forms/request/fields/request-url-field.vue';

import DeclareTicketRuleTicketMappingField from './declare-ticket-rule-ticket-mapping-field.vue';

export default {
  components: { RequestUrlField, DeclareTicketRuleTicketMappingField, RequestWithTokenForm },
  model: {
    prop: 'form',
    event: 'input',
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
  setup(props, { emit }) {
    const requestFormName = computed(() => `${props.name}.request`);
    const webhookTemplateVarsKey = computed(() => (props.hasPrevious ? 'webhook' : 'first_webhook'));

    const { confirmAction: removeWebhook } = useConfirmableForm({
      form: toRef(props, 'form'),
      action: () => emit('remove'),
      cloning: true,
    });

    return {
      requestFormName,
      webhookTemplateVarsKey,

      removeWebhook,
    };
  },
};
</script>
