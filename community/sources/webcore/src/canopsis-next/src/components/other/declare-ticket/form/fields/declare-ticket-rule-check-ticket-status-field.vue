<template>
  <v-layout column>
    <c-information-block
      :title="$t('declareTicket.checkTicketStatus')"
      :help-text="$t('declareTicket.checkTicketStatusHelpText')"
      help-icon="help"
      help-icon-color="grey darken-1"
    >
      <v-layout class="gap-3" column>
        <c-enabled-field v-field="form.enabled" />
        <v-expand-transition>
          <v-layout v-if="form.enabled" class="gap-3" column>
            <request-with-token-form
              v-field="form"
              :name="`${name}.request`"
              :disabled="disabled"
              :url-variables="templateVars.ticket"
              :headers-variables="templateVars.ticket"
              :payload-variables="templateVars.ticket"
              :hide-auth="form.reuse_headers_and_auth"
              hide-repeat
            >
              <template #additional-fields>
                <c-enabled-field
                  v-field="form.reuse_headers_and_auth"
                  :label="$t('declareTicket.reuseHeadersAndAuthFromTicketDeclarationRule')"
                  :disabled="disabled"
                />
              </template>
            </request-with-token-form>
            <declare-ticket-rule-ticket-status-source-field
              v-field="form.ticket_status"
              :name="`${name}.ticket_status_source`"
              :disabled="disabled"
              :variables="templateVars.ticket"
            />
            <declare-ticket-rule-ticket-status-mapping-field
              v-field="form.status_mapping"
              :name="`${name}.status_mapping`"
              :disabled="disabled"
            />
          </v-layout>
        </v-expand-transition>
      </v-layout>
    </c-information-block>
  </v-layout>
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
  },
};
</script>
