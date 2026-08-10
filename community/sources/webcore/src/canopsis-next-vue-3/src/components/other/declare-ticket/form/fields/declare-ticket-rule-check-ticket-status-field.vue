<template>
  <div>
    <c-form-block-row
      :depth="depth"
      :label="$t('declareTicket.checkTicketStatus')"
      :top-border="topBorder"
    >
      <c-enabled-field
        :model-value="form.enabled"
        @update:model-value="updateField('enabled', $event)"
      >
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
          :form="form"
          :name="`${name}.request`"
          :disabled="disabled"
          :url-variables="ticketTemplateVars"
          :headers-variables="ticketTemplateVars"
          :payload-variables="ticketTemplateVars"
          :hide-auth="form.reuse_headers_and_auth"
          :hide-headers="form.reuse_headers_and_auth"
          :url-label="$t('declareTicket.ticketStatusEndpoint')"
          :depth="depth + 1"
          hide-repeat
          @update:form="updateModel"
        >
          <template #additional-fields>
            <c-form-block-row
              :depth="depth + 1"
              :label="$t('declareTicket.reuseHeadersAndAuthFromTicketDeclarationRule')"
            >
              <c-enabled-field
                :model-value="form.reuse_headers_and_auth"
                :label="$t('declareTicket.reuseHeadersAndAuthFromTicketDeclarationRule')"
                :disabled="disabled"
                @update:model-value="updateField('reuse_headers_and_auth', $event)"
              />
            </c-form-block-row>
          </template>
        </request-with-token-form>

        <declare-ticket-rule-ticket-status-source-field
          :form="form.ticket_status"
          :name="`${name}.ticket_status_source`"
          :disabled="disabled"
          :variables="ticketStatusTemplateVars"
          :depth="depth + 1"
          @update:form="updateField('ticket_status', $event)"
        />

        <declare-ticket-rule-ticket-status-mapping-field
          :value="form.status_mapping"
          :name="`${name}.status_mapping`"
          :disabled="disabled"
          :depth="depth + 1"
          @update:value="updateField('status_mapping', $event)"
        />
      </div>
    </v-expand-transition>
  </div>
</template>

<script>
import { computed } from 'vue';

import RequestWithTokenForm from '@/components/forms/request/request-with-token-form.vue';
import { formMixin } from '@/mixins/form';

import DeclareTicketRuleTicketStatusMappingField from './declare-ticket-rule-ticket-status-mapping-field.vue';
import DeclareTicketRuleTicketStatusSourceField from './declare-ticket-rule-ticket-status-source-field.vue';

export default {
  components: {
    RequestWithTokenForm,
    DeclareTicketRuleTicketStatusSourceField,
    DeclareTicketRuleTicketStatusMappingField,
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
    topBorder: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const ticketTemplateVars = computed(() => props.templateVars.ticket ?? []);
    const ticketStatusTemplateVars = computed(() => props.templateVars.ticket_status ?? []);

    return {
      ticketTemplateVars,
      ticketStatusTemplateVars,
    };
  },
};
</script>
