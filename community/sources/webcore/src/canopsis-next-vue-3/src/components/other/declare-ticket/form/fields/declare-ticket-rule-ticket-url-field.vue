<template>
  <c-layout gap="3" column>
    <declare-ticket-rule-ticket-url-payload-text-field
      v-if="form.template"
      :value="form.value"
      :name="name"
      :disabled="disabled"
      :variables="variables"
      @update:value="updateField('value', $event)"
    />
    <declare-ticket-rule-ticket-url-text-field
      v-else
      :value="form.value"
      :name="name"
      :disabled="disabled"
      @update:value="updateField('value', $event)"
    />
    <c-enabled-field
      :model-value="form.template"
      :label="$t('declareTicket.allowTicketURLTemplate')"
      no-margin
      @update:model-value="updateField('template', $event)"
    />
  </c-layout>
</template>

<script>
import { formMixin } from '@/mixins/form';

import DeclareTicketRuleTicketUrlTextField from './declare-ticket-rule-ticket-url-text-field.vue';
import DeclareTicketRuleTicketUrlPayloadTextField from './declare-ticket-rule-ticket-url-payload-text-field.vue';

export default {
  components: {
    DeclareTicketRuleTicketUrlTextField,
    DeclareTicketRuleTicketUrlPayloadTextField,
  },
  mixins: [formMixin],
  emits: ['update:form'],
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'ticket_url',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    variables: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
