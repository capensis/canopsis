import { get, isEmpty, isEqual } from 'lodash';
import { computed, onBeforeUnmount, ref, watch } from 'vue';

import { MODALS, TEMPLATE_TESTING_TEST_TYPES } from '@/constants';
import { formToRemediationJob } from '@/helpers/entities/remediation/job/form';
import {
  formToTemplateTestingTestValidateForm,
  getChangesForValidateForm,
  formToTemplateTestingTestValidate,
  templateTestingTestValidateToForm,
} from '@/helpers/entities/template-testing-test/form';
import { uid } from '@/helpers/uid';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { usePendingHandler } from '@/hooks/query/pending';
import { useTemplateTesting } from '@/hooks/store/modules/template-testing';
import { useValidationFormErrors } from '@/hooks/validator/validation-form-errors';
import { useValidator } from '@/hooks/validator/validator';

/**
 * Hook for managing template testing validation form
 *
 * @param {Object} props - Component props
 * @param {Object} props.generalForm - The general form data
 * @param {string} props.type - The type of template testing
 * @returns {Object} Validation form state and methods
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
  const isEqualValidationForm = computed(() =>
    isEqual(originalValidationForm.value, validationForm.value),
  );

  const { setFormErrors } = useValidationFormErrors(validationForm);

  watch(
    () => props.active,
    (newActive) => {
      if (!newActive) {
        return;
      }

      const newValidationForm = formToTemplateTestingTestValidateForm(
        props.generalForm,
        props.type,
      );

      const { added, removed } = getChangesForValidateForm(newValidationForm, validationForm.value);

      validationForm.value = validationForm.value.filter(
        item => !removed.some(removedItem => removedItem.key === item.key),
      );
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
    },
    { immediate: true },
  );

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
  const setSelectedTest = newSelectedTest => (selectedTest.value = newSelectedTest);

  /**
   * Resets the selected test to a shallow copy of the current selected test
   */
  const resetSelectedTest = () => setSelectedTest({ ...selectedTest.value });

  /**
   * Applies a new selected test and updates the validation form with the test data
   *
   * @param {Object} newSelectedTest - The new selected test object
   */
  const applyNewSelectedTest = (newSelectedTest) => {
    setSelectedTest(newSelectedTest);

    originalValidationForm.value = templateTestingTestValidateToForm(
      validationForm.value,
      newSelectedTest.data,
    );
    validationForm.value = originalValidationForm.value;
  };

  /**
   * Updates the selected test and syncs the validation form with test data
   *
   * @param {Object} newSelectedTest - The new selected test object
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
 * @returns {Object} Test actions state and methods
 */
export const useTemplateTestActions = (props) => {
  const templateTestRequestData = ref(null);

  const { t } = useI18n();
  const popups = usePopups();
  const modals = useModals();
  const validator = useValidator();

  const { createTemplateTestingTest, updateTemplateTestingTest } = useTemplateTesting();

  /**
   * Saves a template test either by updating existing or preparing for later creation
   *
   * @param {Object} selectedTest - The selected test object
   * @param {Function} getValidationFormData - Function to get validation form data
   * @returns {Promise<void>}
   */
  const save = async (selectedTest, getValidationFormData) => {
    const data = {
      ...selectedTest,
      data: getValidationFormData(),
    };

    const isValid = await validator.validateAll();

    if (!isValid) {
      return;
    }

    if (props.ruleId) {
      data.rule = props.ruleId;

      await updateTemplateTestingTest({ data, id: selectedTest._id });

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
    const isValid = await validator.validateAll();

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

          const newTest = await createTemplateTestingTest({ data });

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

  const getTemplateTestRequestData = () => templateTestRequestData.value;

  const clearTemplateTestRequestData = () => {
    templateTestRequestData.value = null;
  };

  return {
    templateTestRequestData,
    save,
    saveAsNew,
    getTemplateTestRequestData,
    clearTemplateTestRequestData,
  };
};

/**
 * Hook for managing template test execution
 *
 * @param {Object} props - Component props
 * @returns {Object} Test runner state and methods
 */
export const useTemplateTestRunner = (props) => {
  const testResult = ref({});

  const validator = useValidator();

  const { validateJobsVariables } = useTemplateTesting();

  const validateHandler = computed(
    () =>
      ({
        [TEMPLATE_TESTING_TEST_TYPES.job]: validateJobsVariables,
      })[props.type],
  );

  const formToRequest = computed(
    () =>
      ({
        [TEMPLATE_TESTING_TEST_TYPES.job]: formToRemediationJob,
      })[props.type],
  );

  const createRunTestHandler = (
    getValidationFormData,
    setFormErrors,
    onValidationStateChange,
    onLastRunVariablesChange,
  ) => {
    const { pending: running, handler: runTest } = usePendingHandler(async () => {
      validator.errors.clear();

      const isValid = await validator.validateAll();

      if (onValidationStateChange) {
        onValidationStateChange(isValid);
      }

      if (!isValid || !validateHandler.value || !formToRequest.value) {
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
        const { mainFormErrors, testDataErrors } = Object.entries(err).reduce(
          (acc, [key, value]) => {
            if (key.startsWith('rule.')) {
              acc.mainFormErrors[key.replace('rule.', '')] = value;
            } else {
              acc.testDataErrors[key] = value;
            }

            return acc;
          },
          { mainFormErrors: {}, testDataErrors: {} },
        );

        if (!isEmpty(mainFormErrors) && onValidationStateChange) {
          onValidationStateChange(false);
        }

        setFormErrors(mainFormErrors);
        setFormErrors(testDataErrors, 'testdata');
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
 *
 * @returns {Object}
 */
export const useTemplateTestAfterSubmit = () => {
  onBeforeUnmount(() => {});

  return {};
};
