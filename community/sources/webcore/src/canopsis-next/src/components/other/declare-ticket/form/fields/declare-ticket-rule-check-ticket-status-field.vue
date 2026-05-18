<template>
  <div>
    <c-form-block-row :depth="depth" :label="$t('declareTicket.checkTicketStatus')" top-border>
      <c-enabled-field v-field="form.enabled">
        <template #append>
          <c-help-icon
            :text="$t('declareTicket.checkTicketStatusHelpText')"
            icon="help"
            color="grey darken-1"
            left
          />
        </template>
      </c-enabled-field>
    </c-form-block-row>

    <v-expand-transition>
      <div v-if="form.enabled">
        <request-with-token-form
          v-field="form"
          :name="`${name}.request`"
          :disabled="disabled"
          :url-variables="templateVars.ticket"
          :headers-variables="templateVars.ticket"
          :payload-variables="templateVars.ticket"
          :hide-auth="form.reuse_headers_and_auth"
          :hide-headers="form.reuse_headers_and_auth"
          :url-label="$t('declareTicket.ticketStatusEndpoint')"
          :depth="depth + 1"
          hide-repeat
        >
          <template #additional-fields>
            <c-form-block-row
              :depth="depth + 1"
              :label="$t('declareTicket.reuseHeadersAndAuthFromTicketDeclarationRule')"
            >
              <c-enabled-field
                v-field="form.reuse_headers_and_auth"
                :label="$t('declareTicket.reuseHeadersAndAuthFromTicketDeclarationRule')"
                :disabled="disabled"
              />
            </c-form-block-row>
          </template>
        </request-with-token-form>

        <declare-ticket-rule-ticket-status-source-field
          v-field="form.ticket_status"
          :name="`${name}.ticket_status_source`"
          :disabled="disabled"
          :variables="templateVars.ticket_status"
          :depth="depth + 1"
        />

        <declare-ticket-rule-ticket-status-mapping-field
          v-field="form.status_mapping"
          :name="`${name}.status_mapping`"
          :disabled="disabled"
          :depth="depth + 1"
        />
      </div>
    </v-expand-transition>
  </div>
</template>

<script>
import RequestWithTokenForm from '@/components/forms/request/request-with-token-form.vue';

import DeclareTicketRuleTicketStatusSourceField from './declare-ticket-rule-ticket-status-source-field.vue';
import DeclareTicketRuleTicketStatusMappingField from './declare-ticket-rule-ticket-status-mapping-field.vue';

export default {
  components: {
    RequestWithTokenForm,
    DeclareTicketRuleTicketStatusSourceField,
    DeclareTicketRuleTicketStatusMappingField,
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
      default: 'check_ticket_status',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
    depth: {
      type: Number,
      default: 0,
    },
  },
};
</script>
