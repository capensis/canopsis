<template>
  <div>
    <c-form-block-row :label="$t('common.systemName')" :depth="depth">
      <c-name-field
        v-field="form.ticket_system_name"
        :label="$t('common.systemName')"
        :name="systemNameFieldName"
        autofocus
      />
    </c-form-block-row>

    <c-form-block-row :label="$t('declareTicket.ticketID')" :depth="depth">
      <declare-ticket-rule-ticket-id-text-field
        v-field="form.ticket"
        :name="ticketFieldName"
        required
      />
    </c-form-block-row>

    <c-form-block-row :label="$t('declareTicket.ticketURL')" :depth="depth">
      <v-layout column>
        <declare-ticket-rule-ticket-url-text-field
          v-field="form.ticket_url"
          :name="ticketUrlFieldName"
        />
        <declare-ticket-rule-ticket-url-title-field v-field="form.ticket_url_title" />
      </v-layout>
    </c-form-block-row>

    <c-form-block-row :label="$tc('common.customField', 2)" :depth="depth" indented>
      <declare-ticket-rule-ticket-custom-fields-field
        v-field="form.mapping"
        :name="name"
      />
    </c-form-block-row>
  </div>
</template>

<script>
import { computed } from 'vue';

import DeclareTicketRuleTicketIdTextField from './fields/declare-ticket-rule-ticket-id-text-field.vue';
import DeclareTicketRuleTicketUrlTextField from './fields/declare-ticket-rule-ticket-url-text-field.vue';
import DeclareTicketRuleTicketCustomFieldsField from './fields/declare-ticket-rule-ticket-custom-fields-field.vue';
import DeclareTicketRuleTicketUrlTitleField from './fields/declare-ticket-rule-ticket-url-title-field.vue';

export default {
  components: {
    DeclareTicketRuleTicketUrlTitleField,
    DeclareTicketRuleTicketCustomFieldsField,
    DeclareTicketRuleTicketUrlTextField,
    DeclareTicketRuleTicketIdTextField,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
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
