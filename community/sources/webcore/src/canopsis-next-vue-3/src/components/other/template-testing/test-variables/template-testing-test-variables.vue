<template>
  <c-layout class="template-testing-test-variables" gap="4" column>
    <template-testing-test-variables-form
      :form="validationForm"
      :selected-test="selectedTest"
      :variables-fields="variablesFields"
      :type="type"
      :rule-id="ruleId"
      :test-result="testResult"
      :last-run-variables="lastRunVariables"
      :template-vars="templateVars"
      @update:form="updateValidationForm"
      @update:selected-test="updateSelectedTest"
    />

    <c-alert v-if="!isGeneralFormValid" type="error">
      {{ t('templateTesting.mainFormHasErrors') }}
    </c-alert>

    <c-layout gap="2" justify="end">
      <v-btn
        :disabled="isDisabledSaveAsNewButton"
        class="template-testing-test-variables__secondary-btn"
        color="secondary"
        variant="outlined"
        @click="saveAsNew"
      >
        {{ t('templateTesting.saveTestAsNew') }}
      </v-btn>

      <c-tooltip :disabled="!saveButtonTooltip" location="top">
        <template #activator="{ props: activatorProps }">
          <span v-bind="activatorProps">
            <v-btn :disabled="isDisabledSaveButton" color="secondary" @click="save">
              {{ t('templateTesting.saveTest') }}
            </v-btn>
          </span>
        </template>

        <span>{{ saveButtonTooltip }}</span>
      </c-tooltip>

      <v-btn
        :loading="running"
        :disabled="running"
        color="primary"
        @click="runTest"
      >
        {{ t('templateTesting.runTest') }}
      </v-btn>
    </c-layout>
  </c-layout>
</template>

<script>
import { computed, provide, ref } from 'vue';

import { TEMPLATE_STUDIO_FORMS_CONTEXT } from '@/constants';

import { useForm, useFormContext } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import TemplateTestingTestVariablesForm from './form/template-testing-test-variables-form.vue';
import {
  useTemplateTestingValidationForm,
  useTemplateTestActions,
  useTemplateTestRunner,
  useTemplateTestAfterSubmit,
} from './hooks/template-testing-test-variables';

export default {
  components: {
    TemplateTestingTestVariablesForm,
  },

  props: {
    generalForm: {
      type: Object,
      default: () => ({}),
    },
    variablesFields: {
      type: Array,
      default: () => [],
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
    ruleId: {
      type: String,
      required: false,
    },
    type: {
      type: Number,
      required: false,
    },
    active: {
      type: Boolean,
      default: false,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const mainFormContext = useFormContext();
    const testFormContext = useForm();
    const testDataFormContext = useForm();

    provide(TEMPLATE_STUDIO_FORMS_CONTEXT, {
      mainFormContext,
      testFormContext,
      testDataFormContext,
    });

    const isGeneralFormValid = ref(true);

    const {
      validationForm,
      selectedTest,
      lastRunVariables,
      hasSelectedTest,
      isEmptyValidationForm,
      isEqualValidationForm,
      getValidationFormData,
      updateSelectedTest: updateSelectedTestHook,
      resetLastRunVariables,
      setLastRunVariables,
    } = useTemplateTestingValidationForm(props);

    const {
      save: saveAction,
      saveAsNew: saveAsNewAction,
      getTemplateTestRequestData,
      clearTemplateTestRequestData,
    } = useTemplateTestActions(props, {
      testFormContext,
      testDataFormContext,
    });

    const { testResult, createRunTestHandler } = useTemplateTestRunner(props, {
      mainFormContext,
      testDataFormContext,
    });

    const { running, runTest } = createRunTestHandler(
      getValidationFormData,
      (isValid) => {
        isGeneralFormValid.value = isValid;
      },
      setLastRunVariables,
    );

    const isDisabledSaveAsNewButton = computed(() => isEmptyValidationForm.value);

    const isDisabledSaveButton = computed(
      () => !hasSelectedTest.value || isEmptyValidationForm.value || isEqualValidationForm.value,
    );

    const saveButtonTooltip = computed(() => hasSelectedTest.value && isEqualValidationForm.value ? t('templateTesting.testSaved') : '',
    );
    const hasAnyError = computed(
      () =>
        !!Object.keys(testFormContext.errors.value).length ||
        !!Object.keys(testDataFormContext.errors.value).length,
    );

    const save = () => saveAction(selectedTest.value, getValidationFormData);

    const saveAsNew = () => saveAsNewAction(getValidationFormData, (newTest) => {
        selectedTest.value = newTest;
      });

    const updateValidationForm = (value) => {
      validationForm.value = value;
      resetLastRunVariables();
    };

    const updateSelectedTest = (newSelectedTest) => {
      updateSelectedTestHook(newSelectedTest);

      if (getTemplateTestRequestData()) {
        clearTemplateTestRequestData();
      }
    };

    useTemplateTestAfterSubmit(
      {
        getTemplateTestRequestData,
      },
      props,
    );

    return {
      t,
      selectedTest,
      isGeneralFormValid,
      validationForm,
      isDisabledSaveAsNewButton,
      isDisabledSaveButton,
      saveButtonTooltip,
      testResult,
      lastRunVariables,
      hasAnyError,
      running,
      saveAsNew,
      save,
      runTest,
      updateValidationForm,
      updateSelectedTest,
    };
  },
};
</script>
