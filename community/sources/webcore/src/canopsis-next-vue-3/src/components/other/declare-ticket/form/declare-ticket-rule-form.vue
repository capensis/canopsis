<template>
  <c-layout class="position-relative" gap="3" column>
    <c-enabled-field v-model="enabledField" with-background />

    <c-form-general-patterns-tabs :form="values" :rule-id="ruleId" :type="type">
      <template #general="{ setRef, templateVars }">
        <declare-ticket-rule-general-form :ref="setRef" :template-vars="templateVars" />
      </template>
      <template #patterns="{ setRef }">
        <declare-ticket-rule-patterns-form :ref="setRef" />
      </template>
      <template #test-query>
        <declare-ticket-rule-test-query :form="values" />
      </template>
    </c-form-general-patterns-tabs>
  </c-layout>
</template>

<script>
import { TEMPLATE_TESTING_TEST_TYPES } from '@/constants';
import { useFormContext, useFormModel } from '@/hooks/form';

import DeclareTicketRuleTestQuery from '../partials/declare-ticket-rule-test-query.vue';

import DeclareTicketRuleGeneralForm from './declare-ticket-rule-general-form.vue';
import DeclareTicketRulePatternsForm from './declare-ticket-rule-patterns-form.vue';

export default {
  components: {
    DeclareTicketRuleTestQuery,
    DeclareTicketRulePatternsForm,
    DeclareTicketRuleGeneralForm,
  },
  props: {
    ruleId: {
      type: String,
      required: false,
    },
  },
  setup() {
    const { values } = useFormContext();
    const enabledField = useFormModel('enabled');
    const type = TEMPLATE_TESTING_TEST_TYPES.declareTicketRule;

    return {
      values,
      enabledField,
      type,
    };
  },
};
</script>
