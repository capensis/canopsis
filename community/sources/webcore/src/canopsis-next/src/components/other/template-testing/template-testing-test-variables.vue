<template>
  <v-layout class="gap-4" column>
    <template-testing-test-variables-form
      v-model="validationForm"
      :selected-test="selectedTest"
      :variables-fields="variablesFields"
      :type="type"
      :rule-id="ruleId"
      :test-result="testResult"
      :last-run-variables="lastRunVariables"
      :template-vars="templateVars"
      @input="validationFormInput"
      @update:selected-test="updateSelectedTest"
    />
    <c-alert :value="!isGeneralFormValid" type="error">
      {{ $t('templateTesting.mainFormHasErrors') }}
    </c-alert>
    <v-layout class="gap-2" justify-end>
      <v-btn
        class="template-testing-test-variables__secondary-btn"
        color="secondary"
        outlined
        @click="saveAsNew"
      >
        {{ $t('templateTesting.saveTestAsNew') }}
      </v-btn>
      <v-btn
        :disabled="!selectedTest"
        color="secondary"
        @click="save"
      >
        {{ $t('templateTesting.saveTest') }}
      </v-btn>
      <v-btn
        :loading="running"
        :disabled="running"
        color="primary"
        @click="runTest"
      >
        {{ $t('templateTesting.runTest') }}
      </v-btn>
    </v-layout>
  </v-layout>
</template>

<script>
import { ref } from 'vue';

import {
  useTemplateTestingValidationForm,
  useTemplateTestActions,
  useTemplateTestRunner,
  useTemplateTestAfterSubmit,
} from './hooks/template-testing-test-variables';
import TemplateTestingTestVariablesForm from './template-testing-test-variables-form.vue';

export default {
  inject: ['$validator', '$afterSubmitObserver'],
  components: {
    TemplateTestingTestVariablesForm,
  },
  model: {
    prop: 'generalForm',
    event: 'input',
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
  },
  setup(props) {
    const isGeneralFormValid = ref(true);

    // Use validation form hook
    const {
      validationForm,
      selectedTest,
      lastRunVariables,
      setFormErrors,
      getValidationFormData,
      updateSelectedTest: updateSelectedTestHook,
      resetLastRunVariables,
      setLastRunVariables,
    } = useTemplateTestingValidationForm(props);

    // Use test actions hook
    const {
      save: saveAction,
      saveAsNew: saveAsNewAction,
      getTemplateTestRequestData,
      clearTemplateTestRequestData,
    } = useTemplateTestActions(props);

    // Use test runner hook
    const {
      testResult,
      createRunTestHandler,
    } = useTemplateTestRunner(props);

    // Create run test handler with callbacks
    const { running, runTest } = createRunTestHandler(
      getValidationFormData,
      setFormErrors,
      (isValid) => { isGeneralFormValid.value = isValid; },
      setLastRunVariables,
    );

    // Use after submit hook
    const { setupAfterSubmitObserver } = useTemplateTestAfterSubmit({
      getTemplateTestRequestData,
    }, props);

    // Setup after submit observer
    setupAfterSubmitObserver();

    // Create wrapped save methods
    const save = () => saveAction(selectedTest.value, getValidationFormData);

    const saveAsNew = () => saveAsNewAction(
      getValidationFormData,
      (newTest) => { selectedTest.value = newTest; },
    );

    const validationFormInput = () => {
      resetLastRunVariables();
    };

    const updateSelectedTest = (newSelectedTest) => {
      updateSelectedTestHook(newSelectedTest);

      if (getTemplateTestRequestData()) {
        clearTemplateTestRequestData();
      }
    };

    return {
      selectedTest,
      isGeneralFormValid,
      validationForm,
      testResult,
      lastRunVariables,
      running,
      saveAsNew,
      save,
      runTest,
      validationFormInput,
      updateSelectedTest,
    };
  },
};
</script>

<style lang="scss" scoped>
.variables-item {
  display: flex;
  align-items: center;
  gap: 4px;

  &__name {
    font-weight: 500;
  }

  &__separator {
    color: rgba(0, 0, 0, 0.6);
  }

  &__value {
    flex: 1;
    min-width: 0;
  }

  &__empty-object {
    color: rgba(0, 0, 0, 0.6);
  }
}

.theme--dark {
  --secondaryBtnColor: #C7C7C7;

  .variables-item {
    &__separator {
      color: rgba(255, 255, 255, 0.7);
    }

    &__empty-object {
      color: rgba(255, 255, 255, 0.7);
    }
  }

  .template-testing-test-variables__secondary-btn {
    color: var(--secondaryBtnColor) !important;
    border-color: var(--secondaryBtnColor) !important;
  }
}
</style>
