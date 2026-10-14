<template>
  <c-layout gap="3" column>
    <declare-ticket-rule-ticket-id-payload-text-field
      v-if="value.ticket_id.template"
      :value="value.ticket_id.value"
      :name="name"
      :disabled="disabled"
      :variables="variables"
      :required="required"
      @update:value="updateField('ticket_id.value', $event)"
    />
    <declare-ticket-rule-ticket-id-text-field
      v-else
      :value="value.ticket_id.value"
      :name="name"
      :disabled="disabled"
      :required="required"
      @update:value="updateField('ticket_id.value', $event)"
    />
    <c-enabled-field
      :model-value="value.ticket_id.template"
      :label="$t('declareTicket.allowTicketIdTemplate')"
      hide-details
      no-margin
      @update:model-value="updateField('ticket_id.template', $event)"
    />
    <c-enabled-field
      :model-value="value.is_regexp"
      :label="$t('declareTicket.isRegexp')"
      @update:model-value="updateField('is_regexp', $event)"
    />
  </c-layout>
</template>

<script>
import { formMixin } from '@/mixins/form';

import DeclareTicketRuleTicketIdPayloadTextField from './declare-ticket-rule-ticket-id-payload-text-field.vue';
import DeclareTicketRuleTicketIdTextField from './declare-ticket-rule-ticket-id-text-field.vue';

export default {
  components: {
    DeclareTicketRuleTicketIdTextField,
    DeclareTicketRuleTicketIdPayloadTextField,
  },
  mixins: [formMixin],
  emits: ['update:value'],
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'ticket_id',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
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
