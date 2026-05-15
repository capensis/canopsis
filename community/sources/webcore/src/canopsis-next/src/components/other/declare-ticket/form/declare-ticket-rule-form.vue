<template>
  <v-layout class="gap-3" column>
    <c-enabled-field v-field="form.enabled" with-background />

    <c-form-general-patterns-tabs
      :form="form"
      :rule-id="ruleId"
      :type="type"
    >
      <template #general="{ setRef, templateVars }">
        <declare-ticket-rule-general-form
          v-field="form"
          :ref="setRef"
          :template-vars="templateVars"
        />
      </template>
      <template #patterns="{ setRef }">
        <declare-ticket-rule-patterns-form
          v-field="form.patterns"
          :ref="setRef"
        />
      </template>
      <template #test-query>
        <declare-ticket-rule-test-query :form="form" />
      </template>
    </c-form-general-patterns-tabs>
  </v-layout>
</template>

<script>
import { TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

import DeclareTicketRuleTestQuery from '../partials/declare-ticket-rule-test-query.vue';

import DeclareTicketRuleGeneralForm from './declare-ticket-rule-general-form.vue';
import DeclareTicketRulePatternsForm from './declare-ticket-rule-patterns-form.vue';

export default {
  components: {

    DeclareTicketRuleTestQuery,
    DeclareTicketRulePatternsForm,
    DeclareTicketRuleGeneralForm,
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
    ruleId: {
      type: String,
      required: false,
    },
  },
  setup() {
    const type = TEMPLATE_TESTING_TEST_TYPES.declareTicketRule;

    return {
      type,
    };
  },
};
</script>
