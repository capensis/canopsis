<template>
  <c-form-block-row :label="$t('declareTicket.ticketStatusSourceField')" :depth="1" top-border>
    <v-col cols="6">
      <c-payload-text-field
        v-if="form.template"
        :value="form.value"
        :label="$t('declareTicket.ticketStatusSourceField')"
        :name="name"
        :variables="variables"
        :disabled="disabled"
        required
        @update:value="updateField('value', $event)"
      />
      <c-name-field
        v-else
        :model-value="form.value"
        :label="$t('declareTicket.ticketStatusSourceField')"
        :name="name"
        :disabled="disabled"
        required
        @update:model-value="updateField('value', $event)"
      />
    </v-col>
    <c-enabled-field
      :model-value="form.template"
      :label="$t('declareTicket.allowTicketStatusSourceTemplate')"
      @update:model-value="updateField('template', $event)"
    />
  </c-form-block-row>
</template>

<script>
import { formMixin } from '@/mixins/form';

export default {
  mixins: [formMixin],
  emits: ['update:form'],
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'ticket_status_source',
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
