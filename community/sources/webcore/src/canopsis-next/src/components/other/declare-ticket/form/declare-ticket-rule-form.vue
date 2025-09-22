<template>
  <v-tabs
    slider-color="primary"
    fixed-tabs
  >
    <v-tab :class="{ 'error--text': hasGeneralError }">
      {{ $t('common.general') }}
    </v-tab>
    <v-tab :class="{ 'error--text': hasPatternsError }">
      {{ $tc('common.pattern') }}
    </v-tab>
    <v-tab>{{ $t('common.testQuery') }}</v-tab>

    <template-testing-test-variables-tab :disabled="isEmptyVariablesFields" />

    <v-tab-item eager>
      <declare-ticket-rule-general-form
        v-field="form"
        ref="general"
        :template-vars="templateVars"
        class="mt-2"
      />
    </v-tab-item>
    <v-tab-item eager>
      <declare-ticket-rule-patterns-form
        v-field="form.patterns"
        ref="patterns"
        class="mt-2"
      />
    </v-tab-item>
    <v-tab-item>
      <v-layout>
        <v-flex
          offset-xs1
          xs10
        >
          <declare-ticket-rule-test-query
            :form="form"
            class="mt-2"
          />
        </v-flex>
      </v-layout>
    </v-tab-item>
    <v-tab-item :disabled="isEmptyVariablesFields">
      <template-testing-test-variables
        :general-form="form"
        :variables-fields="variablesFields"
        :template-vars="templateVars"
        :rule-id="ruleId"
        :type="type"
      />
    </v-tab-item>
  </v-tabs>
</template>

<script>
import { ref, toRef, watch, onMounted } from 'vue';

import { TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

import { useTemplateVarsList } from '@/hooks/vars/template';

import { useTestVariablesFields } from '@/components/other/template-testing/hooks/template-test-variables-wrapper';

import TemplateTestingTestVariables from '@/components/other/template-testing/template-testing-test-variables.vue';
import TemplateTestingTestVariablesTab from '@/components/other/template-testing/template-testing-test-variables-tab.vue';

import DeclareTicketRuleTestQuery from '../partials/declare-ticket-rule-test-query.vue';

import DeclareTicketRuleGeneralForm from './declare-ticket-rule-general-form.vue';
import DeclareTicketRulePatternsForm from './declare-ticket-rule-patterns-form.vue';

export default {
  components: {
    TemplateTestingTestVariables,
    TemplateTestingTestVariablesTab,

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
  setup(props, { emit }) {
    const type = TEMPLATE_TESTING_TEST_TYPES.declareTicketRule;

    const hasGeneralError = ref(false);
    const hasPatternsError = ref(false);

    const general = ref(null);
    const patterns = ref(null);

    const {
      vars: templateVars,
      pending: templateVarsPending,
      fetchList: fetchTemplateVarsList,
    } = useTemplateVarsList({
      type,
      form: toRef(props, 'form'),
    });

    const {
      items: variablesFields,
      isEmptyItems: isEmptyVariablesFields,
    } = useTestVariablesFields(props, type, emit);

    watch(() => general.value?.hasAnyError, value => hasGeneralError.value = value);
    watch(() => patterns.value?.hasAnyError, value => hasPatternsError.value = value);

    onMounted(() => {
      fetchTemplateVarsList();
    });

    return {
      type,

      hasGeneralError,
      hasPatternsError,

      general,
      patterns,

      variablesFields,
      isEmptyVariablesFields,

      templateVars,
      templateVarsPending,
    };
  },
};
</script>
