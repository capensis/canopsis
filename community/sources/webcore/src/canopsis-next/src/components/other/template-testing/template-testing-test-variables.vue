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
} from './test-variables/hooks/template-testing-test-variables';
import TemplateTestingTestVariablesForm from './test-variables/form/template-testing-test-variables-form.vue';

/**
 * Template Testing Test Variables Component
 *
 * Component for managing template testing with test variables, validation,
 * and execution. Provides UI for running tests, saving tests, and viewing results.
 *
 * @vue/component
 */
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
    /**
     * The general form data containing template configuration
     * @type {Object}
     */
    generalForm: {
      type: Object,
      default: () => ({}),
    },
    /**
     * Array of variable fields for the template
     * @type {Array}
     */
    variablesFields: {
      type: Array,
      default: () => [],
    },
    /**
     * Template variables object
     * @type {Object}
     */
    templateVars: {
      type: Object,
      default: () => ({}),
    },
    /**
     * Rule ID for existing rules (optional)
     * @type {string}
     */
    ruleId: {
      type: String,
      required: false,
    },
    type: {
      type: Number,
      required: false,
    },
  },
  /**
   * Vue Composition API setup function
   *
   * @param {Object} props - Component props
   * @param {Object} props.generalForm - The general form data
   * @param {Array} props.variablesFields - Array of variable fields
   * @param {Object} props.templateVars - Template variables object
   * @param {string} props.ruleId - Rule ID for existing rules
   * @param {number} props.type - Type of template testing
   * @returns {Object} Setup return object with reactive data and methods
   */
  setup(props) {
    const isGeneralFormValid = ref(true);

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

    const {
      save: saveAction,
      saveAsNew: saveAsNewAction,
      getTemplateTestRequestData,
      clearTemplateTestRequestData,
    } = useTemplateTestActions(props);

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

    useTemplateTestAfterSubmit({ getTemplateTestRequestData }, props);

    /**
     * Saves the current selected test
     */
    const save = () => saveAction(selectedTest.value, getValidationFormData);

    /**
     * Saves the current test data as a new test with user-provided name
     */
    const saveAsNew = () => saveAsNewAction(
      getValidationFormData,
      (newTest) => { selectedTest.value = newTest; },
    );

    /**
     * Handles validation form input changes and resets last run variables
     */
    const validationFormInput = () => {
      resetLastRunVariables();
    };

    /**
     * Updates the selected test and clears pending request data if exists
     *
     * @param {Object} newSelectedTest - The new selected test object
     */
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
