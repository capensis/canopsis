<template>
  <div>
    <c-form-block-row :label="$t('common.systemName')" :depth="depth">
      <c-name-field
        :model-value="form.ticket_system_name"
        :label="$t('common.systemName')"
        :name="systemNameFieldName"
        autofocus
        @update:model-value="updateField('ticket_system_name', $event)"
      />
    </c-form-block-row>

    <c-form-block-row :label="$t('declareTicket.ticketID')" :depth="depth">
      <declare-ticket-rule-ticket-id-text-field
        :value="form.ticket"
        :name="ticketFieldName"
        required
        @update:value="updateField('ticket', $event)"
      />
    </c-form-block-row>

    <c-form-block-row :label="$t('declareTicket.ticketURL')" :depth="depth">
      <c-layout column>
        <declare-ticket-rule-ticket-url-text-field
          :value="form.ticket_url"
          :name="ticketUrlFieldName"
          @update:value="updateField('ticket_url', $event)"
        />
        <declare-ticket-rule-ticket-url-title-field
          :value="form.ticket_url_title"
          @update:value="updateField('ticket_url_title', $event)"
        />
      </c-layout>
    </c-form-block-row>

    <c-form-block-row :label="$t('common.customField', 2)" :depth="depth">
      <declare-ticket-rule-ticket-custom-fields-field
        :value="form.mapping"
        :name="name"
        @update:value="updateField('mapping', $event)"
      />
    </c-form-block-row>
  </div>
</template>

<script>
import { computed } from 'vue';

import { formMixin } from '@/mixins/form';

import DeclareTicketRuleTicketCustomFieldsField from './fields/declare-ticket-rule-ticket-custom-fields-field.vue';
import DeclareTicketRuleTicketIdTextField from './fields/declare-ticket-rule-ticket-id-text-field.vue';
import DeclareTicketRuleTicketUrlTextField from './fields/declare-ticket-rule-ticket-url-text-field.vue';
import DeclareTicketRuleTicketUrlTitleField from './fields/declare-ticket-rule-ticket-url-title-field.vue';

export default {
  components: {
    DeclareTicketRuleTicketUrlTitleField,
    DeclareTicketRuleTicketCustomFieldsField,
    DeclareTicketRuleTicketUrlTextField,
    DeclareTicketRuleTicketIdTextField,
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
      required: false,
    },
    depth: {
      type: Number,
      default: 0,
    },
  },
  setup(props) {
    const prepareFieldName = fieldName => [props.name, fieldName].filter(Boolean).join('.');

    const systemNameFieldName = computed(() => prepareFieldName('system_name'));
    const ticketFieldName = computed(() => prepareFieldName('ticket'));
    const ticketUrlFieldName = computed(() => prepareFieldName('ticket_url'));

    return {
      systemNameFieldName,
      ticketFieldName,
      ticketUrlFieldName,
    };
  },
};
</script>
