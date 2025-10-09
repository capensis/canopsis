import { get, isEmpty, isEqual } from 'lodash';
import {
  computed,
  inject,
  onBeforeUnmount,
  onMounted,
  ref,
  watch,
} from 'vue';

import { MODALS, TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

import { uid } from '@/helpers/uid';
import { formToService } from '@/helpers/entities/service/form';
import { formToEventFilter } from '@/helpers/entities/event-filter/rule/form';
import { formToLinkRule } from '@/helpers/entities/link/form';
import { formToScenario } from '@/helpers/entities/scenario/form';
import { formToWidget } from '@/helpers/entities/widget/form';
import { formToDeclareTicketRule } from '@/helpers/entities/declare-ticket/rule/form';
import { formToDynamicInfo } from '@/helpers/entities/dynamic-info/rule/form';
import { formToRemediationInstructionRequest } from '@/helpers/entities/remediation/instruction/form';
import { formToRemediationJob } from '@/helpers/entities/remediation/job/form';
import { formToMetaAlarmRule } from '@/helpers/entities/meta-alarm/rule/form';
import {
  formToTemplateTestingTestValidateForm,
  getChangesForValidateForm,
  formToTemplateTestingTestValidate,
  templateTestingTestValidateToForm,
} from '@/helpers/entities/template-testing-test/form';

import { useI18n } from '@/hooks/i18n';
import { usePopups } from '@/hooks/popups';
import { useModals } from '@/hooks/modals';
import { usePendingHandler } from '@/hooks/query/pending';
import { useValidator } from '@/hooks/validator/validator';
import { useTemplateTest } from '@/hooks/store/modules/template-test';
import { useTemplateValidation } from '@/hooks/store/modules/template-validation';
import { useValidationFormErrors } from '@/hooks/validator/validation-form-errors';

/**
 * Hook for managing template testing validation form
 *
 * @param {Object} props - Component props
 * @param {Object} props.generalForm - The general form data
 * @param {string} props.type - The type of template testing
 * @returns {Object} Validation form state and methods
 * @returns {import('vue').Ref<Array>} returns.validationForm - Reactive validation form array
 * @returns {import('vue').Ref<Object>} returns.selectedTest - Reactive selected test object
 * @returns {import('vue').Ref<Object>} returns.lastRunVariables - Reactive last run variables object
 * @returns {Function} returns.setFormErrors - Function to set form validation errors
 * @returns {Function} returns.getValidationFormData - Function to get validation form data
 * @returns {Function} returns.updateSelectedTest - Function to update selected test
 * @returns {Function} returns.resetLastRunVariables - Function to reset last run variables
 * @returns {Function} returns.setLastRunVariables - Function to set last run variables
 */
export const useTemplateTestingValidationForm = (props) => {
  const originalValidationForm = ref([]);
  const validationForm = ref([]);
  const selectedTest = ref({});
  const lastRunVariables = ref({});

  const { t } = useI18n();
  const modals = useModals();

  const hasSelectedTest = computed(() => !isEmpty(selectedTest.value));
  const isEmptyValidationForm = computed(() => validationForm.value.every(item => !item.value));
  const isEqualValidationForm = computed(() => isEqual(originalValidationForm.value, validationForm.value));

  const { setFormErrors } = useValidationFormErrors(validationForm);

  watch(() => props.active, (newActive) => {
    if (!newActive) {
      return;
    }

    const newValidationForm = formToTemplateTestingTestValidateForm(props.generalForm, props.type);

    const { added, removed } = getChangesForValidateForm(newValidationForm, validationForm.value);

    validationForm.value = validationForm.value
      .filter(item => !removed.some(removedItem => removedItem.key === item.key));
    validationForm.value.push(...added);

    validationForm.value.forEach((_, index) => {
      if (newValidationForm[index]) {
        validationForm.value[index] = {
          ...validationForm.value[index],
          index: newValidationForm[index].index,
          params: newValidationForm[index].params,
        };
      }
    });
  }, { immediate: true });

  /**
   * Gets the validation form data in the proper format for API requests
   *
   * @returns {Object} The formatted validation form data
   */
  const getValidationFormData = () => formToTemplateTestingTestValidate(validationForm.value);

  /**
   * Sets the selected test data
   *
   * @param {Object} newSelectedTest - The new selected test object to set
   */
  const setSelectedTest = newSelectedTest => selectedTest.value = newSelectedTest;

  /**
   * Resets the selected test to a shallow copy of the current selected test
   * This preserves the current test data while creating a new object reference
   */
  const resetSelectedTest = () => setSelectedTest({ ...(selectedTest.value || {}) });

  /**
   * Applies a new selected test and updates the validation form with the test data
   *
   * @param {Object} newSelectedTest - The new selected test object
   * @param {Object} newSelectedTest.data - The test data to be used for validation form
   */
  const applyNewSelectedTest = (newSelectedTest) => {
    setSelectedTest(newSelectedTest);

    originalValidationForm.value = templateTestingTestValidateToForm(validationForm.value, newSelectedTest.data);
    validationForm.value = originalValidationForm.value;
  };

  /**
   * Updates the selected test and syncs the validation form with test data
   *
   * @param {Object} newSelectedTest - The new selected test object
   * @param {Object} newSelectedTest.data - The test data
   */
  const updateSelectedTest = (newSelectedTest) => {
    if (!newSelectedTest) {
      setSelectedTest({});

      return;
    }

    if (isEmptyValidationForm.value || isEqualValidationForm.value) {
      applyNewSelectedTest(newSelectedTest);

      return;
    }

    modals.show({
      name: MODALS.confirmation,
      config: {
        text: t('templateTesting.lostCurrentTest.text'),
        action: () => applyNewSelectedTest(newSelectedTest),
        actionText: t('templateTesting.lostCurrentTest.actionText'),
        cancel: resetSelectedTest,
        cancelText: t('templateTesting.lostCurrentTest.cancelText'),
        cancelOutlined: true,
      },
    });
  };

  /**
   * Resets last run variables by replacing all values with unique symbols
   * This forces Vue reactivity to detect changes when the same values are set again
   */
  const resetLastRunVariables = () => {
    lastRunVariables.value = Object.fromEntries(
      Object.keys(lastRunVariables.value).map(key => [key, Symbol('lastRunVariables')]),
    );
  };

  /**
   * Sets the last run variables to the provided values
   *
   * @param {Object} variables - The variables object to set
   */
  const setLastRunVariables = (variables) => {
    lastRunVariables.value = { ...variables };
  };

  return {
    validationForm,
    selectedTest,
    lastRunVariables,

    hasSelectedTest,
    isEmptyValidationForm,
    isEqualValidationForm,

    setFormErrors,
    getValidationFormData,
    updateSelectedTest,
    resetLastRunVariables,
    setLastRunVariables,
  };
};

/**
 * Hook for managing template test save actions
 *
 * @param {Object} props - Component props
 * @param {string} props.ruleId - The rule ID if editing an existing rule
 * @param {string} props.type - The type of template testing
 * @returns {Object} Test actions state and methods
 * @returns {import('vue').Ref<Object|null>} returns.templateTestRequestData - Reactive template test request data
 * @returns {Function} returns.save - Function to save a test
 * @returns {Function} returns.saveAsNew - Function to save a test as new
 * @returns {Function} returns.getTemplateTestRequestData - Function to get template test request data
 * @returns {Function} returns.createTemplateTestFromRequestData - Function to create template test from request data
 * @returns {Function} returns.clearTemplateTestRequestData - Function to clear template test request data
 */
export const useTemplateTestActions = (props) => {
  const templateTestRequestData = ref(null);

  const { t } = useI18n();
  const popups = usePopups();
  const modals = useModals();
  const validator = useValidator();

  const { createTemplateTest, updateTemplateTest } = useTemplateTest();

  /**
   * Saves a template test either by updating existing or preparing for later creation
   *
   * @param {Object} selectedTest - The selected test object
   * @param {string} selectedTest._id - The test ID
   * @param {Function} getValidationFormData - Function to get validation form data
   * @returns {Promise<void>}
   */
  const save = async (selectedTest, getValidationFormData) => {
    const data = {
      ...selectedTest,
      data: getValidationFormData(),
    };

    const [isNameValid, isDataValid] = await Promise.all([
      validator.validateAll('test-name'),
      validator.validateAll('test-data'),
    ]);

    if (!isNameValid || !isDataValid) {
      return;
    }

    if (props.ruleId) {
      data.rule = props.ruleId;

      await updateTemplateTest({ data, id: selectedTest._id });

      popups.success({ text: t('templateTesting.testSaved') });

      return;
    }

    templateTestRequestData.value = data;

    popups.success({ text: t('templateTesting.testWillSaveAfterFormSaving') });
  };

  /**
   * Saves a template test as a new test with a custom name
   *
   * @param {Function} getValidationFormData - Function to get validation form data
   * @param {Function} [onTestCreated] - Optional callback when test is created
   * @returns {Promise<void>}
   */
  const saveAsNew = async (getValidationFormData, onTestCreated) => {
    const isValid = await validator.validateAll('test-data');

    if (!isValid) {
      return;
    }

    const modalConfig = {
      title: t('templateTesting.saveTestAsNew'),
      field: {
        text: '',
        name: 'name',
        label: t('common.testName'),
        validationRules: {
          required: true,
        },
      },
      action: async (newTestName) => {
        const data = {
          data: getValidationFormData(),
          description: '',
          name: newTestName,
          type: props.type,
        };

        if (props.ruleId) {
          data.rule = props.ruleId;

          const newTest = await createTemplateTest({ data });

          popups.success({ text: t('templateTesting.testSaved') });

          if (onTestCreated) {
            onTestCreated(newTest);
          }

          return;
        }

        templateTestRequestData.value = data;

        const newTestData = {
          ...data,
          _id: uid('test'),
        };

        if (onTestCreated) {
          onTestCreated(newTestData);
        }

        popups.success({ text: t('templateTesting.testWillSaveAfterFormSaving') });
      },
    };

    if (!props.ruleId) {
      modalConfig.alert = { text: t('templateTesting.isNewForm') };
    }

    modals.show({
      name: MODALS.textFieldEditor,
      config: modalConfig,
    });
  };

  /**
   * Gets the current template test request data
   *
   * @returns {Object|null} The template test request data or null
   */
  const getTemplateTestRequestData = () => templateTestRequestData.value;

  /**
   * Creates a template test from stored request data and clears it
   *
   * @returns {Promise<Object|null>} The created test or null if no request data
   */
  const createTemplateTestFromRequestData = async () => {
    if (!templateTestRequestData.value) {
      return null;
    }

    const newTest = await createTemplateTest({ data: templateTestRequestData.value });
    templateTestRequestData.value = null;

    return newTest;
  };

  /**
   * Clears the template test request data
   */
  const clearTemplateTestRequestData = () => {
    templateTestRequestData.value = null;
  };

  return {
    templateTestRequestData,
    save,
    saveAsNew,
    getTemplateTestRequestData,
    createTemplateTestFromRequestData,
    clearTemplateTestRequestData,
  };
};

/**
 * Hook for managing template test execution
 *
 * @param {Object} props - Component props
 * @param {string} props.type - The type of template testing
 * @param {Object} props.generalForm - The general form data
 * @param {Array} props.variablesFields - Array of variable fields
 * @returns {Object} Test runner state and methods
 * @returns {import('vue').Ref<Object>} returns.testResult - Reactive test result object
 * @returns {Function} returns.createRunTestHandler - Function to create a test run handler
 */
export const useTemplateTestRunner = (props) => {
  const testResult = ref({});

  const validator = useValidator();

  const {
    validateEntityServices,
    validateEventFilters,
    validateScenarios,
    validateLinkRules,
    validateWidgets,
    validateDeclareTicketRules,
    validateDynamicInfos,
    validateInstructions,
    validateJobs,
    validateMetaAlarmRules,
  } = useTemplateValidation();

  const validateHandler = computed(() => ({
    [TEMPLATE_TESTING_TEST_TYPES.eventFilter]: validateEventFilters,
    [TEMPLATE_TESTING_TEST_TYPES.linkRule]: validateLinkRules,
    [TEMPLATE_TESTING_TEST_TYPES.scenario]: validateScenarios,
    [TEMPLATE_TESTING_TEST_TYPES.widget]: validateWidgets,
    [TEMPLATE_TESTING_TEST_TYPES.declareTicketRule]: validateDeclareTicketRules,
    [TEMPLATE_TESTING_TEST_TYPES.dynamicInfo]: validateDynamicInfos,
    [TEMPLATE_TESTING_TEST_TYPES.instruction]: validateInstructions,
    [TEMPLATE_TESTING_TEST_TYPES.job]: validateJobs,
    [TEMPLATE_TESTING_TEST_TYPES.metaAlarmRule]: validateMetaAlarmRules,
  })[props.type] ?? validateEntityServices);

  const formToRequest = computed(() => ({
    [TEMPLATE_TESTING_TEST_TYPES.eventFilter]: formToEventFilter,
    [TEMPLATE_TESTING_TEST_TYPES.linkRule]: formToLinkRule,
    [TEMPLATE_TESTING_TEST_TYPES.scenario]: formToScenario,
    [TEMPLATE_TESTING_TEST_TYPES.widget]: formToWidget,
    [TEMPLATE_TESTING_TEST_TYPES.declareTicketRule]: formToDeclareTicketRule,
    [TEMPLATE_TESTING_TEST_TYPES.dynamicInfo]: formToDynamicInfo,
    [TEMPLATE_TESTING_TEST_TYPES.instruction]: formToRemediationInstructionRequest,
    [TEMPLATE_TESTING_TEST_TYPES.job]: formToRemediationJob,
    [TEMPLATE_TESTING_TEST_TYPES.metaAlarmRule]: formToMetaAlarmRule,
  })[props.type] ?? formToService);

  /**
   * Creates a test run handler with validation and error handling
   *
   * @param {Function} getValidationFormData - Function to get validation form data
   * @param {Function} setFormErrors - Function to set form validation errors
   * @param {Function} [onValidationStateChange] - Optional callback for validation state changes
   * @param {Function} [onLastRunVariablesChange] - Optional callback for last run variables changes
   * @returns {Object} Test run handler object
   * @returns {import('vue').Ref<boolean>} returns.running - Reactive running state
   * @returns {Function} returns.runTest - Function to run the test
   */
  const createRunTestHandler = (
    getValidationFormData,
    setFormErrors,
    onValidationStateChange,
    onLastRunVariablesChange,
  ) => {
    const { pending: running, handler: runTest } = usePendingHandler(async () => {
      validator.errors.clear('test-data');
      validator.errors.clear('test-result');

      const [isParentValid, isValid] = await Promise.all([
        validator.validateAll(),
        validator.validateAll('test-data'),
      ]);

      if (onValidationStateChange) {
        onValidationStateChange(isParentValid);
      }

      if (!isParentValid || !isValid) {
        return;
      }

      try {
        testResult.value = await validateHandler.value({
          data: {
            rule: formToRequest.value(props.generalForm),
            testdata: getValidationFormData(),
          },
        });

        Object.entries(testResult.value).forEach(([key, { err }]) => {
          if (!err) {
            return;
          }

          validator.errors.add({
            field: key,
            scope: 'test-result',
            msg: `${err.line}|${err.message}`,
          });
        });

        const newLastRunVariables = props.variablesFields.reduce((acc, field) => {
          const { key } = field;

          acc[key] = get(props.generalForm, key);

          return acc;
        }, {});

        if (onLastRunVariablesChange) {
          onLastRunVariablesChange(newLastRunVariables);
        }
      } catch (err) {
        const { mainFormErrors, testDataErrors } = Object.entries(err).reduce((acc, [key, value]) => {
          if (key.startsWith('rule.')) {
            acc.mainFormErrors[key.replace('rule.', '')] = value;
          } else {
            acc.testDataErrors[key] = value;
          }

          return acc;
        }, { mainFormErrors: {}, testDataErrors: {} });

        if (!isEmpty(mainFormErrors) && onValidationStateChange) {
          onValidationStateChange(false);
        }

        setFormErrors(mainFormErrors);
        setFormErrors(testDataErrors, 'test-data');
      }
    });

    return { running, runTest };
  };

  return {
    testResult,
    createRunTestHandler,
  };
};

/**
 * Hook for handling post-submission logic when creating template tests
 * Manages the afterSubmit observer registration and handles test creation after form submission
 *
 * @param {Object} config - Configuration object
 * @param {Function} config.getTemplateTestRequestData - Function that returns template test request data
 * @param {Object} props - Component props
 * @param {string} [props.ruleId] - Optional rule ID that determines if observer should be registered
 * @returns {Object} An object containing the afterSubmitObserver
 * @property {Object} afterSubmitObserver - The observer instance for post-submission handling
 */
export const useTemplateTestAfterSubmit = ({ getTemplateTestRequestData }, props) => {
  const afterSubmitObserver = inject('$afterSubmitObserver');

  const { t } = useI18n();
  const popups = usePopups();

  const { createTemplateTest } = useTemplateTest();

  /**
   * Handles post-submission logic for creating template tests
   *
   * @param {Object} [rule] - The rule object if available
   * @param {string} [rule._id] - The rule ID
   * @returns {Promise<void>}
   */
  const afterSubmit = async (rule) => {
    /**
     * Capture the translation string before the component is destroyed
     * vue-i18n provides a proxy object that is destroyed when the component is destroyed
     * so we need to capture the translation string before the component is destroyed
     */
    const testSavedMessage = t('templateTesting.testSaved');

    const templateTestRequestData = getTemplateTestRequestData();

    if (!templateTestRequestData) {
      return;
    }

    if (rule) {
      templateTestRequestData.rule = rule._id;
    }

    await createTemplateTest({ data: templateTestRequestData });

    popups.success({ text: testSavedMessage });
  };

  onMounted(() => {
    if (props.ruleId) {
      return;
    }

    afterSubmitObserver.register(afterSubmit);
  });

  /**
   * Unregister the afterSubmit observer after a timeout to avoid race condition
   */
  onBeforeUnmount(() => setTimeout(() => afterSubmitObserver.unregister(afterSubmit), 0));

  return {
    afterSubmitObserver,
  };
};
