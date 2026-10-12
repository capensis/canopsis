<template>
  <v-layout class="position-relative gap-2" column>
    <c-progress-overlay :pending="templateVarsPending" />
    <c-enabled-field v-field="form.enabled" with-background />
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
      <v-tab
        :class="{ 'error--text': hasSimpleError || errors.has('links') }"
        :disabled="sourceCodeWasChanged"
      >
        {{ $t('linkRule.simpleMode') }}
      </v-tab>
      <v-tab :class="{ 'error--text': hasAdvancedError || errors.has('links') }">
        {{ $t('linkRule.advancedMode') }}
      </v-tab>
      <template-testing-test-variables-tab v-if="hasAccess" :disabled="isEmptyVariablesFields" />

      <v-tab-item
        class="mt-3"
        eager
      >
        <link-rule-general-form
          v-field="form"
          ref="generalElement"
          :template-vars="templateVars"
          class="mt-2"
        />
      </v-tab-item>

      <v-tab-item
        class="mt-3"
        eager
      >
        <link-rule-patterns-form
          v-field="form.patterns"
          ref="patternsElement"
          :is-alarm-type="isAlarmType"
        />
      </v-tab-item>

      <v-tab-item
        class="mt-3"
        eager
      >
        <link-rule-simple-form
          v-field="form.links"
          ref="simpleElement"
          :type="form.type"
          :template-vars="templateVars"
          :error-messages="linksErrorMessages"
          @input="resetRequiredRule"
        />
      </v-tab-item>
      <v-tab-item
        class="mt-3"
        eager
      >
        <link-rule-advanced-form
          v-field="form.source_code"
          ref="advancedElement"
          :type="form.type"
          :error-messages="linksErrorMessages"
          @input="resetRequiredRule"
        />
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
  </v-layout>
</template>

<script>
import {
  computed,
  ref,
  toRef,
  watch,
  onBeforeUnmount,
  onMounted,
} from 'vue';

import { TEMPLATE_TESTING_TEST_TYPES, LINK_RULE_TYPES } from '@/constants';

import { isDefaultSourceCode } from '@/helpers/entities/link/form';

import { useI18n } from '@/hooks/i18n';
import { useValidator } from '@/hooks/validator/validator';
import { useTemplateVarsList } from '@/hooks/vars/template';
import { useValidationAttachRequired } from '@/hooks/validator/validation-attach-required';
import { useAiChatExpand } from '@/hooks/ai/ai-chat-form';

import {
  useTestVariablesTabData,
} from '@/components/other/template-testing/test-variables/hooks/template-test-variables-wrapper';

import TemplateTestingTestVariables from '@/components/other/template-testing/test-variables/template-testing-test-variables.vue';
import TemplateTestingTestVariablesTab from '@/components/other/template-testing/test-variables/partials/template-testing-test-variables-tab.vue';

import LinkRuleGeneralForm from './link-rule-general-form.vue';
import LinkRuleSimpleForm from './link-rule-simple-form.vue';
import LinkRuleAdvancedForm from './link-rule-advanced-form.vue';
import LinkRulePatternsForm from './link-rule-patterns-form.vue';

const LINK_RULE_FORM_TABS = {
  general: 0,
  patterns: 1,
  simple: 2,
  advanced: 3,
  testing: 4,
};

export default {
  inject: ['$validator'],
  components: {
    TemplateTestingTestVariables,
    TemplateTestingTestVariablesTab,

    LinkRuleGeneralForm,
    LinkRulePatternsForm,
    LinkRuleSimpleForm,
    LinkRuleAdvancedForm,
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
    const { t } = useI18n();
    const { errors } = useValidator();

    const activeTab = ref(LINK_RULE_FORM_TABS.general);

    const type = TEMPLATE_TESTING_TEST_TYPES.linkRule;

    const hasGeneralError = ref(false);
    const hasPatternsError = ref(false);
    const hasSimpleError = ref(false);
    const hasAdvancedError = ref(false);

    const generalElement = ref(null);
    const patternsElement = ref(null);
    const simpleElement = ref(null);
    const advancedElement = ref(null);

    const isActiveTestingTab = computed(() => activeTab.value === LINK_RULE_FORM_TABS.testing);
    const isAlarmType = computed(() => props.form.type === LINK_RULE_TYPES.alarm);

    const linksErrorMessages = computed(() => (errors.has('links') ? [t('linkRule.linksEmptyError')] : []));

    const {
      attachRequiredRule,
      detachRequiredRule,
      resetRequiredRule,
    } = useValidationAttachRequired('links');

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

    const sourceCodeWasChanged = computed(() => !isDefaultSourceCode(props.form.source_code));

    /**
     * Getter function for validation rule that checks if links are required.
     * Returns true if there are links in the form or if source code was modified from default.
     *
     * @returns {boolean} True if links field should be considered valid (has links or custom source code)
     */
    const requiredRuleGetter = () => !!props.form.links.length || !isDefaultSourceCode(props.form.source_code);

    watch(() => generalElement.value?.hasAnyError, value => hasGeneralError.value = value);
    watch(() => patternsElement.value?.hasAnyError, value => hasPatternsError.value = value);
    watch(() => simpleElement.value?.hasAnyError, value => hasSimpleError.value = value);
    watch(() => advancedElement.value?.hasAnyError, value => hasAdvancedError.value = value);

    useAiChatExpand({ activeTab, neededTab: LINK_RULE_FORM_TABS.general });

    onMounted(() => {
      attachRequiredRule(requiredRuleGetter);
      fetchTemplateVarsList();
    });

    onBeforeUnmount(detachRequiredRule);

    return {
      LINK_RULE_FORM_TABS,

      activeTab,
      isActiveTestingTab,
      isAlarmType,
      linksErrorMessages,

      type,

      hasGeneralError,
      hasPatternsError,
      hasSimpleError,
      hasAdvancedError,

      generalElement,
      patternsElement,
      simpleElement,
      advancedElement,

      hasAccess,
      variablesFields,
      isEmptyVariablesFields,

      templateVars,
      templateVarsPending,

      sourceCodeWasChanged,

      resetRequiredRule,
    };
  },
};
</script>
