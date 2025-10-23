<template>
  <v-tabs
    v-model="activeTab"
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

    <template-testing-test-variables-tab v-if="hasAccess" :disabled="isEmptyVariablesFields" />

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
    <v-tab-item v-if="hasAccess" :disabled="isEmptyVariablesFields">
      <template-testing-test-variables
        :general-form="form"
        :variables-fields="variablesFields"
        :template-vars="templateVars"
        :rule-id="ruleId"
        :type="type"
        :active="isActiveTestingTab"
      />
    </v-tab-item>
  </v-tabs>
</template>

<script>
import {
  computed,
  ref,
  toRef,
  watch,
  onMounted,
} from 'vue';

import { TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

import { useTemplateVarsList } from '@/hooks/vars/template';

import {
  useTestVariablesTabData,
} from '@/components/other/template-testing/test-variables/hooks/template-test-variables-wrapper';

import TemplateTestingTestVariables from '@/components/other/template-testing/test-variables/template-testing-test-variables.vue';
import TemplateTestingTestVariablesTab from '@/components/other/template-testing/test-variables/partials/template-testing-test-variables-tab.vue';

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
    const activeTab = ref(0);

    const type = TEMPLATE_TESTING_TEST_TYPES.declareTicketRule;

    const hasGeneralError = ref(false);
    const hasPatternsError = ref(false);

    const general = ref(null);
    const patterns = ref(null);

    const isActiveTestingTab = computed(() => activeTab.value === 3);

    const {
      vars: templateVars,
      pending: templateVarsPending,
      fetchList: fetchTemplateVarsList,
    } = useTemplateVarsList({
      type,
      form: toRef(props, 'form'),
    });

    const {
      hasAccess,

      items: variablesFields,
      isEmptyItems: isEmptyVariablesFields,
    } = useTestVariablesTabData(props, type, emit);

    watch(() => general.value?.hasAnyError, value => hasGeneralError.value = value);
    watch(() => patterns.value?.hasAnyError, value => hasPatternsError.value = value);

    onMounted(() => {
      fetchTemplateVarsList();
    });

    return {
      activeTab,
      isActiveTestingTab,

      type,

      hasGeneralError,
      hasPatternsError,

      general,
      patterns,

      hasAccess,
      variablesFields,
      isEmptyVariablesFields,

      templateVars,
      templateVarsPending,
    };
  },
};
</script>
