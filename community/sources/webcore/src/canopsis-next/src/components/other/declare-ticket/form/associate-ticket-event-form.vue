<template>
  <c-information-block :title="$t('declareTicket.ticketUrlAndId')">
    <c-name-field
      v-field="form.system_name"
      :label="$t('common.systemName')"
      :name="systemNameFieldName"
    />
    <v-layout>
      <v-flex
        class="mr-3"
        xs6
      >
        <declare-ticket-rule-ticket-id-text-field
          v-field="form.ticket_id"
          :name="ticketIdFieldName"
          required
        />
      </v-flex>
      <v-flex xs6>
        <declare-ticket-rule-ticket-url-text-field
          v-field="form.ticket_url"
          :name="ticketUrlFieldName"
        />
      </v-flex>
    </v-layout>
    <declare-ticket-rule-ticket-custom-fields-field
      v-field="form.mapping"
      :name="name"
    />
  </c-information-block>
</template>

<script>
import DeclareTicketRuleTicketIdTextField from './fields/declare-ticket-rule-ticket-id-text-field.vue';
import DeclareTicketRuleTicketUrlTextField from './fields/declare-ticket-rule-ticket-url-text-field.vue';
import DeclareTicketRuleTicketCustomFieldsField from './fields/declare-ticket-rule-ticket-custom-fields-field.vue';

export default {
  components: {
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
  },
  computed: {
    systemNameFieldName() {
      return this.prepareFieldName('system_name');
    },

    ticketIdFieldName() {
      return this.prepareFieldName('ticket_id');
    },

    ticketUrlFieldName() {
      return this.prepareFieldName('ticket_url');
    },
  },
  methods: {
    prepareFieldName(name) {
      return [this.name, name].filter(Boolean).join('.');
    },
  },
};
</script>
