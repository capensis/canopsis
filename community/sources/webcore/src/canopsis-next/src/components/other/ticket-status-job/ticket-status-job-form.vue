<template>
  <c-layout gap="3" column>
    <c-layout gap="3" wrap>
      <c-text-field
        :model-value="ruleType"
        :label="t('jobs.ruleType')"
        readonly
        disabled
      />

      <c-text-field
        :model-value="form.rule_name"
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

      <c-text-field
        :model-value="form.ticket_system_name"
        :label="t('jobs.ticketSystemName')"
        readonly
        disabled
      />
    </c-layout>

    <c-text-field
      :model-value="form.ticket_id"
      :label="$tc('common.ticket')"
      :disabled="!isNew"
      name="ticket"
      required
      @update:model-value="updateField('ticket_id', $event)"
    />

    <c-form-block>
      <declare-ticket-rule-check-ticket-status-field
        :form="form.check_ticket_status"
        :template-vars="templateVars"
        @update:form="updateField('check_ticket_status', $event)"
      />
    </c-form-block>
  </c-layout>
</template>

<script>
import { computed, onMounted } from 'vue';

import { TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';
import { useTemplateVarsList } from '@/hooks/vars/template';

import DeclareTicketRuleCheckTicketStatusField from '@/components/other/declare-ticket/form/fields/declare-ticket-rule-check-ticket-status-field.vue';

export default {
  components: {
    DeclareTicketRuleCheckTicketStatusField,
  },

  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    isNew: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:form'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateField } = useModelField(props, emit);
    const ruleType = computed(() => t(`jobs.types.${props.form.rule_type}`));
    const {
      vars: templateVars,
      fetchList: fetchTemplateVarsList,
    } = useTemplateVarsList({ type: TEMPLATE_TESTING_TEST_TYPES.declareTicketRule });

    onMounted(fetchTemplateVarsList);

    return {
      t,
      ruleType,
      templateVars,
      updateField,
    };
  },
};
</script>
