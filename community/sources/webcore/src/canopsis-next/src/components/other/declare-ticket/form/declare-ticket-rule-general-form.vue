<template>
  <c-layout gap="5" column>
    <v-row>
      <v-col cols="6">
        <c-name-field v-model="nameField" autofocus required />
      </v-col>
      <v-col cols="6">
        <c-name-field v-model="systemNameField" :label="t('common.systemName')" name="system_name" />
      </v-col>
    </v-row>

    <c-enabled-field
      v-model="emitTriggerSuccessField"
      :label="t('common.emitTriggerSuccess')"
      name="emit_trigger_success"
      hide-details
      no-margin
    />
    <c-enabled-field
      v-model="emitTriggerFailField"
      :label="t('common.emitTriggerFail')"
      name="emit_trigger_fail"
      hide-details
      no-margin
    />

    <declare-ticket-rule-webhooks-field
      v-model:webhooks="webhooksField"
      :template-vars="templateVars"
    />
  </c-layout>
</template>

<script>
import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useValidationHeader } from '@/hooks/validator/validation-header';

import DeclareTicketRuleWebhooksField from './fields/declare-ticket-rule-webhooks-field.vue';

export default {
  components: { DeclareTicketRuleWebhooksField },
  props: {
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },

  setup() {
    const { t } = useI18n();
    const { hasAnyError } = useValidationHeader();
    const {
      nameField,
      systemNameField,
      emitTriggerSuccessField,
      emitTriggerFailField,
      webhooksField,
    } = useFormModels({
      nameField: 'name',
      systemNameField: 'system_name',
      emitTriggerSuccessField: 'emit_trigger_success',
      emitTriggerFailField: 'emit_trigger_fail',
      webhooksField: 'webhooks',
    });

    return {
      /**
       * It's using in the parent component to display the validation header color for tabs.
       */
      hasAnyError,
      t,
      nameField,
      systemNameField,
      emitTriggerSuccessField,
      emitTriggerFailField,
      webhooksField,
    };
  },
};
</script>
