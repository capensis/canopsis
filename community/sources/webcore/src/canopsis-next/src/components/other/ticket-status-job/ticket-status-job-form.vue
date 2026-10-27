<template>
  <v-row>
    <v-col cols="12" md="4">
      <c-text-field
        :model-value="ruleType"
        :label="t('jobs.ruleType')"
        readonly
        disabled
      />
    </v-col>

    <v-col cols="12" md="4">
      <c-text-field
        :model-value="ruleNameField"
        :label="t('jobs.ruleName')"
        readonly
        disabled
      >
        <template #append-inner>
          <c-help-icon
            :text="t('jobs.ruleNameTooltip')"
            icon="help"
            icon-class="text-grey"
            left
          />
        </template>
      </c-text-field>
    </v-col>

    <v-col cols="12" md="4">
      <c-text-field
        :model-value="ticketSystemNameField"
        :label="t('jobs.ticketSystemName')"
        readonly
        disabled
      />
    </v-col>

    <v-col cols="12">
      <c-text-field
        v-model="ticketIdField"
        :label="t('common.ticket')"
        :disabled="!isNew"
        name="ticket"
        required
      />
    </v-col>

    <v-col cols="12">
      <c-form-block>
        <declare-ticket-rule-check-ticket-status-field
          v-model:form="checkTicketStatusField"
          :with-templates="false"
        />
      </c-form-block>
    </v-col>
  </v-row>
</template>

<script>
import { computed } from 'vue';

import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import DeclareTicketRuleCheckTicketStatusField from '@/components/other/declare-ticket/form/fields/declare-ticket-rule-check-ticket-status-field.vue';

export default {
  components: {
    DeclareTicketRuleCheckTicketStatusField,
  },

  props: {
    isNew: {
      type: Boolean,
      default: false,
    },
  },

  setup() {
    const { t } = useI18n();
    const {
      ruleNameField,
      ruleTypeField,
      ticketSystemNameField,
      ticketIdField,
      checkTicketStatusField,
    } = useFormModels({
      ruleNameField: 'rule_name',
      ruleTypeField: 'rule_type',
      ticketSystemNameField: 'ticket_system_name',
      ticketIdField: 'ticket_id',
      checkTicketStatusField: 'check_ticket_status',
    });
    const ruleType = computed(() => t(`jobs.types.${ruleTypeField.value}`));

    return {
      t,
      ruleType,
      ruleNameField,
      ticketSystemNameField,
      ticketIdField,
      checkTicketStatusField,
    };
  },
};
</script>
