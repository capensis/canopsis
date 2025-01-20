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
        :url-variables="payloadVariables"
      />
    </template>
    <request-form
      v-field="form.request"
      :name="requestFormName"
      :headers-variables="payloadVariables"
      :payload-variables="payloadVariables"
      hide-url
    />
    <declare-ticket-rule-ticket-mapping-field
      v-field="form"
      :name="`${name}.declare_ticket`"
      :is-declare-ticket-exist="isDeclareTicketExist"
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
import { confirmableFormMixinCreator } from '@/mixins/confirmable-form';
import { payloadVariablesMixin } from '@/mixins/payload/variables';

import RequestForm from '@/components/forms/request/request-form.vue';
import RequestUrlField from '@/components/forms/request/fields/request-url-field.vue';

import DeclareTicketRuleTicketMappingField from './declare-ticket-rule-ticket-mapping-field.vue';

export default {
  components: { RequestUrlField, DeclareTicketRuleTicketMappingField, RequestForm },
  mixins: [
    payloadVariablesMixin,
    confirmableFormMixinCreator({
      field: 'form',
      method: 'removeWebhook',
      cloning: true,
    }),
  ],
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
  },
  computed: {
    requestFormName() {
      return `${this.name}.request`;
    },
  },
  methods: {
    removeWebhook() {
      this.$emit('remove');
    },
  },
};
</script>
