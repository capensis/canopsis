<template>
  <c-layout gap="3" column>
    <template-testing-test-variables-form-provider :form-context="testFormContext">
      <template-testing-test-field
        :value="selectedTest"
        :rule-id="preparedRuleId"
        name="test.name"
        return-object
        required
        @update:value="updateSelectedTest"
      />
    </template-testing-test-variables-form-provider>

    <c-alert v-if="!ruleId" type="warning">
      {{ t('templateTesting.isNewForm') }}
    </c-alert>

    <template-testing-test-variables-form-provider :form-context="testDataFormContext">
      <template-testing-test-variables-test-data
        :form="form"
        :type="type"
        @update:form="updateForm"
      />

      <c-layout gap="3" column>
        <template-testing-test-variables-template-item
          v-for="field in variablesFields"
          :key="field.key"
          :template="field.value"
          :name="field.resultKey ?? field.key"
          :title="t(field.textKey, field.textArgs)"
          :textarea="field.textarea"
          :json="field.json"
          :result="testResult[field.resultKey ?? field.key]"
          :variables="templateVars[field.templateVarsKey]"
          :last-run-value="lastRunVariables[field.key]"
          @update:template="field.updateField"
        />
      </c-layout>
    </template-testing-test-variables-form-provider>
  </c-layout>
</template>

<script>
import { computed, inject } from 'vue';

import { TEMPLATE_STUDIO_FORMS_CONTEXT } from '@/constants';

import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';

import TemplateTestingTestField from '../../form/fields/template-testing-test-field.vue';
import TemplateTestingTestVariablesFormProvider from '../partials/template-testing-test-variables-form-provider.vue';
import TemplateTestingTestVariablesTemplateItem from '../partials/template-testing-test-variables-template-item.vue';
import TemplateTestingTestVariablesTestData from '../partials/template-testing-test-variables-test-data.vue';

export default {
  components: {
    TemplateTestingTestField,
    TemplateTestingTestVariablesFormProvider,
    TemplateTestingTestVariablesTestData,
    TemplateTestingTestVariablesTemplateItem,
  },

  props: {
    form: {
      type: Array,
      default: () => [],
    },
    selectedTest: {
      type: Object,
      default: () => ({}),
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
    variablesFields: {
      type: Array,
      default: () => [],
    },
    type: {
      type: Number,
      required: false,
    },
    ruleId: {
      type: String,
      required: false,
    },
    testResult: {
      type: Object,
      default: () => ({}),
    },
    lastRunVariables: {
      type: Object,
      default: () => ({}),
    },
  },

  emits: ['update:form', 'update:selected-test'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateModel } = useModelField(props, emit);
    const { testFormContext, testDataFormContext } = inject(TEMPLATE_STUDIO_FORMS_CONTEXT);

    const preparedRuleId = computed(() => props.ruleId ?? 'fake-id');

    /**
     * Updates the selected test by emitting the new selected test value
     *
     * @param {Object} newSelectedTest - The new selected test object
     */
    const updateSelectedTest = newSelectedTest => emit('update:selected-test', newSelectedTest);

    return {
      t,
      preparedRuleId,
      testFormContext,
      testDataFormContext,
      updateForm: updateModel,
      updateSelectedTest,
    };
  },
};
</script>
