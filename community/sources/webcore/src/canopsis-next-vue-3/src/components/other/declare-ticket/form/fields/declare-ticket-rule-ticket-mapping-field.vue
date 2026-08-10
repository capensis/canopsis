<template>
  <div>
    <c-form-block-row :label="$t('declareTicket.ticketUrlAndId')" :depth="depth" top-border>
      <c-enabled-field
        :model-value="form.declare_ticket.enabled"
        :disabled="isDeclareTicketExist"
        @update:model-value="updateField('declare_ticket.enabled', $event)"
      >
        <template #append>
          <c-help-icon
            :text="ticketUrlHelpText"
            icon="help"
            color="grey darken-1"
            left
          />
        </template>
      </c-enabled-field>

      <c-alert v-if="isDeclareTicketExist" type="info">
        {{ $t('declareTicket.webhookTicketDeclarationExist') }}
      </c-alert>
    </c-form-block-row>

    <v-expand-transition>
      <div v-if="form.declare_ticket.enabled">
        <c-form-block-row
          v-if="!hideEmptyResponse"
          :depth="depth + 1"
          :label="$t('declareTicket.emptyResponse')"
        >
          <c-enabled-field
            :model-value="form.declare_ticket.empty_response"
            :label="$t('declareTicket.emptyResponse')"
            @update:model-value="updateField('declare_ticket.empty_response', $event)"
          />
        </c-form-block-row>

        <c-form-block-row :depth="depth + 1" :label="$t('declareTicket.ticketID')">
          <declare-ticket-rule-ticket-id-field
            :value="form.declare_ticket"
            :disabled="disabled"
            :name="ticketIdFieldName"
            :required="ticketIdRequired"
            :variables="variables"
            @update:value="updateField('declare_ticket', $event)"
          />
        </c-form-block-row>

        <c-form-block-row :depth="depth + 1" :label="$t('declareTicket.ticketURL')">
          <declare-ticket-rule-ticket-url-field
            :form="form.declare_ticket.ticket_url"
            :disabled="disabled"
            :name="ticketUrlFieldName"
            :variables="variables"
            @update:form="updateField('declare_ticket.ticket_url', $event)"
          />

          <declare-ticket-rule-ticket-url-title-field
            :value="form.declare_ticket.ticket_url_title"
            @update:value="updateField('declare_ticket.ticket_url_title', $event)"
          />

          <c-text-field
            v-if="withTicketSystemName"
            :model-value="form.ticket_system_name"
            :label="$t('declareTicket.ticketSystemName')"
            @update:model-value="updateField('ticket_system_name', $event)"
          />
        </c-form-block-row>

        <c-form-block-row :depth="depth + 1" :label="$t('common.customField', 2)" indented>
          <declare-ticket-rule-ticket-custom-fields-field
            :value="form.declare_ticket.mapping"
            :name="name"
            :disabled="disabled"
            @update:value="updateField('declare_ticket.mapping', $event)"
          />
        </c-form-block-row>
      </div>
    </v-expand-transition>
  </div>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';
import { formMixin } from '@/mixins/form';

import DeclareTicketRuleTicketCustomFieldsField from './declare-ticket-rule-ticket-custom-fields-field.vue';
import DeclareTicketRuleTicketIdField from './declare-ticket-rule-ticket-id-field.vue';
import DeclareTicketRuleTicketUrlField from './declare-ticket-rule-ticket-url-field.vue';
import DeclareTicketRuleTicketUrlTitleField from './declare-ticket-rule-ticket-url-title-field.vue';

export default {
  components: {
    DeclareTicketRuleTicketUrlField,
    DeclareTicketRuleTicketCustomFieldsField,
    DeclareTicketRuleTicketIdField,
    DeclareTicketRuleTicketUrlTitleField,
  },
  mixins: [formMixin],
  emits: ['update:form'],
  props: {
    form: {
      type: Object,
      required: true,
    },
    name: {
      type: String,
      default: 'declare_ticket',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    isDeclareTicketExist: {
      type: Boolean,
      default: false,
    },
    hideEmptyResponse: {
      type: Boolean,
      default: false,
    },
    ticketIdRequired: {
      type: Boolean,
      default: false,
    },
    onlyOneTicketId: {
      type: Boolean,
      default: false,
    },
    withTicketSystemName: {
      type: Boolean,
      default: false,
    },
    variables: {
      type: Array,
      default: () => [],
    },
    depth: {
      type: Number,
      default: 0,
    },
  },
  setup(props) {
    const { t } = useI18n();

    const ticketIdFieldName = computed(() => `${props.name}.ticket_id`);

    const ticketUrlFieldName = computed(() => `${props.name}.ticket_url`);

    const ticketUrlHelpText = computed(() =>
      [
        t('declareTicket.ticketUrlAndIdHelpText'),
        props.onlyOneTicketId && t('declareTicket.dataFromOneStepAttention'),
      ]
        .filter(Boolean)
        .join('\n'),
    );

    return {
      ticketIdFieldName,
      ticketUrlFieldName,
      ticketUrlHelpText,
    };
  },
};
</script>
