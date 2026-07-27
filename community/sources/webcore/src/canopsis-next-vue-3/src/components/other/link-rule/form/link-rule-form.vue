<template>
  <c-layout class="position-relative" gap="2" column>
    <c-progress-overlay :pending="templateVarsPending" />
    <c-enabled-field v-model="enabledField" with-background />
    <v-tabs v-model="activeTab" slider-color="primary" fixed-tabs>
      <v-tab :value="LINK_RULE_FORM_TABS.general" :class="{ 'error--text': hasGeneralError }">
        {{ $t('common.general') }}
      </v-tab>
      <v-tab :value="LINK_RULE_FORM_TABS.patterns" :class="{ 'error--text': hasPatternsError }">
        {{ $t('common.pattern', 2) }}
      </v-tab>
      <v-tab
        :value="LINK_RULE_FORM_TABS.simple"
        :class="{ 'error--text': hasSimpleError || hasLinksError }"
        :disabled="sourceCodeWasChanged"
      >
        {{ $t('linkRule.simpleMode') }}
      </v-tab>
      <v-tab
        :value="LINK_RULE_FORM_TABS.advanced"
        :class="{ 'error--text': hasAdvancedError || hasLinksError }"
      >
        {{ $t('linkRule.advancedMode') }}
      </v-tab>
      <template-testing-test-variables-tab
        v-if="hasAccess"
        :value="LINK_RULE_FORM_TABS.testing"
        :disabled="isEmptyVariablesFields"
      />
    </v-tabs>

    <v-window v-model="activeTab">
      <v-window-item :value="LINK_RULE_FORM_TABS.general" class="mt-3">
        <link-rule-general-form ref="generalElement" :template-vars="templateVars" class="mt-2" />
      </v-window-item>

      <v-window-item :value="LINK_RULE_FORM_TABS.patterns" class="mt-3">
        <link-rule-patterns-form ref="patternsElement" :is-alarm-type="isAlarmType" />
      </v-window-item>

      <v-window-item :value="LINK_RULE_FORM_TABS.simple" class="mt-3">
        <link-rule-simple-form
          ref="simpleElement"
          :type="typeField"
          :template-vars="templateVars"
          :error-messages="linksErrorMessages"
          @change="resetRequiredRule"
        />
      </v-window-item>
      <v-window-item :value="LINK_RULE_FORM_TABS.advanced" class="mt-3">
        <link-rule-advanced-form
          ref="advancedElement"
          :type="typeField"
          :error-messages="linksErrorMessages"
          @change="resetRequiredRule"
        />
      </v-window-item>

      <v-window-item
        v-if="hasAccess"
        :value="LINK_RULE_FORM_TABS.testing"
        :disabled="isEmptyVariablesFields"
      >
        <template-testing-test-variables
          :general-form="values"
          :variables-fields="variablesFields"
          :template-vars="templateVars"
          :rule-id="ruleId"
          :type="type"
          :active="isActiveTestingTab"
        />
      </v-window-item>
    </v-window>
  </c-layout>
</template>

<script>
import { computed, ref, watch, onBeforeUnmount, onMounted } from 'vue';

import { TEMPLATE_TESTING_TEST_TYPES, LINK_RULE_TYPES } from '@/constants';

import { isDefaultSourceCode } from '@/helpers/entities/link/form';

import { useFormContext, useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useValidator } from '@/hooks/validator/validator';
import { useTemplateVarsList } from '@/hooks/vars/template';
import { useValidationAttachRequired } from '@/hooks/validator/validation-attach-required';
import { useAiChatExpand } from '@/hooks/ai/ai-chat-form';

import { useTestVariablesTabData } from '@/components/other/template-testing/test-variables/hooks/template-test-variables-wrapper';

import TemplateTestingTestVariables from '@/components/other/template-testing/test-variables/template-testing-test-variables.vue';
import TemplateTestingTestVariablesTab from '@/components/other/template-testing/test-variables/partials/template-testing-test-variables-tab.vue';

import LinkRuleGeneralForm from './link-rule-general-form.vue';
import LinkRuleSimpleForm from './link-rule-simple-form.vue';
import LinkRuleAdvancedForm from './link-rule-advanced-form.vue';
import LinkRulePatternsForm from './link-rule-patterns-form.vue';

const LINK_RULE_FORM_TABS = {
  general: 'general',
  patterns: 'patterns',
  simple: 'simple',
  advanced: 'advanced',
  testing: 'testing',
};

export default {
  components: {
    TemplateTestingTestVariables,
    TemplateTestingTestVariablesTab,

    LinkRuleGeneralForm,
    LinkRulePatternsForm,
    LinkRuleSimpleForm,
    LinkRuleAdvancedForm,
  },
  props: {
    ruleId: {
      type: String,
      required: false,
    },
  },

  setup() {
    const { t } = useI18n();
    const { errors } = useValidator();
    const form = useFormContext();

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

    const { enabledField, typeField, sourceCodeField, linksField } = useFormModels({
      enabledField: 'enabled',
      typeField: 'type',
      sourceCodeField: 'source_code',
      linksField: 'links',
    });

    const isActiveTestingTab = computed(() => activeTab.value === LINK_RULE_FORM_TABS.testing);
    const isAlarmType = computed(() => typeField.value === LINK_RULE_TYPES.alarm);
    const hasLinksError = computed(() => errors.has?.('links') || !!form.errors.value?.links);

    const linksErrorMessages = computed(() =>
      hasLinksError.value ? [t('linkRule.linksEmptyError')] : [],
    );

    const { attachRequiredRule, detachRequiredRule, resetRequiredRule } =
      useValidationAttachRequired('links');

    const {
      vars: templateVars,
      pending: templateVarsPending,
      fetchList: fetchTemplateVarsList,
    } = useTemplateVarsList({
      type,
      form: form.values,
    });

    const testingProps = {
      get form() {
        return form.values.value;
      },
    };
    const emitTestingUpdate = (event, value) => {
      if (event === 'input') {
        form.values.value = value;
      }
    };

    const {
      hasAccess,

      items: variablesFields,
      isEmptyItems: isEmptyVariablesFields,
    } = useTestVariablesTabData(testingProps, type, emitTestingUpdate);

    const sourceCodeWasChanged = computed(() => !isDefaultSourceCode(sourceCodeField.value));

    /**
     * Getter function for validation rule that checks if links are required.
     * Returns true if there are links in the form or if source code was modified from default.
     *
     * @returns {boolean} True if links field should be considered valid (has links or custom source code)
     */
    const requiredRuleGetter = () =>
      !!linksField.value.length || !isDefaultSourceCode(sourceCodeField.value);

    watch(
      () => generalElement.value?.hasAnyError,
      (value) => (hasGeneralError.value = value),
    );
    watch(
      () => patternsElement.value?.hasAnyError,
      (value) => (hasPatternsError.value = value),
    );
    watch(
      () => simpleElement.value?.hasAnyError,
      (value) => (hasSimpleError.value = value),
    );
    watch(
      () => advancedElement.value?.hasAnyError,
      (value) => (hasAdvancedError.value = value),
    );

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
      hasLinksError,

      type,
      values: form.values,
      enabledField,
      typeField,

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
