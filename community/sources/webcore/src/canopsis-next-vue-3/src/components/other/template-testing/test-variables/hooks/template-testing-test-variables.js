import { get, isEmpty, isEqual } from 'lodash';
import { computed, onBeforeUnmount, ref, unref, watch } from 'vue';

import { MODALS, TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

import { formToDeclareTicketRule } from '@/helpers/entities/declare-ticket/rule/form';
import { formToDynamicInfo } from '@/helpers/entities/dynamic-info/rule/form';
import { formToEventFilter } from '@/helpers/entities/event-filter/rule/form';
import { formToExternalAuthToken } from '@/helpers/entities/external-auth-token/form';
import { formToLinkRule } from '@/helpers/entities/link/form';
import { formToMetaAlarmRule } from '@/helpers/entities/meta-alarm/rule/form';
import { formToRemediationInstructionRequest } from '@/helpers/entities/remediation/instruction/form';
import { formToRemediationJob } from '@/helpers/entities/remediation/job/form';
import { formToScenario } from '@/helpers/entities/scenario/form';
import { formToService } from '@/helpers/entities/service/form';
import {
  formToTemplateTestingTestValidateForm,
  getChangesForValidateForm,
  formToTemplateTestingTestValidate,
  templateTestingTestValidateToForm,
} from '@/helpers/entities/template-testing-test/form';
import { formToWidget } from '@/helpers/entities/widget/form';
import { uid } from '@/helpers/uid';

import { getServerFieldErrors } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { usePendingHandler } from '@/hooks/query/pending';
import { useTemplateTesting } from '@/hooks/store/modules/template-testing';

/**
 * Extracts template validation errors belonging to the rule form.
 *
 * @param {*} errors
 * @param {Object} [form={}] Rule form containing stable keys for array items.
 * @returns {Object<string, string|string[]>}
 */
export const getTemplateRuleValidationErrors = (errors, form = {}) =>
  Object.entries(getServerFieldErrors(errors)).reduce((result, [key, value]) => {
    if (!key.startsWith('rule.')) {
      return result;
    }

    const field = key
      .slice('rule.'.length)
      .split('.')
      .reduce(
        ({ path, value: formValue }, segment) => {
          if (Array.isArray(formValue) && /^\d+$/.test(segment)) {
            const item = formValue[Number(segment)];

            return {
              path: [...path, item?.key ?? segment],
              value: item,
            };
          }

          return {
            path: [...path, segment],
            value: formValue?.[segment],
          };
        },
        { path: [], value: unref(form) },
      )
      .path.join('.');

    return {
      ...result,
      [field]: value,
    };
  }, {});

/**
 * Maps template validation result errors to registered template field names.
 *
 * `resultKey` is the API response key for fields whose form path differs from
 * the result path. Only configured template fields are returned.
 *
 * @param {Object} [result={}] Template validation result.
 * @param {TemplateTestingTestValidateForm} [fields=[]] Configured template fields.
 * @returns {Object<string, string>}
 */
export const getTemplateFieldsValidationErrors = (result = {}, fields = []) =>
  fields.reduce((errors, field) => {
    const resultKey = field.resultKey ?? field.key;
    const error = result[resultKey]?.err;

    if (!error) {
      return errors;
    }

    return {
      ...errors,
      [`templateTesting.${resultKey}`]: `${error.line}|${error.message}`,
    };
  }, {});

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
   * Resets last run variables with a non-string value so every template is marked as not run.
   */
  const resetLastRunVariables = () => {
    lastRunVariables.value = Object.fromEntries(
      Object.keys(lastRunVariables.value).map(key => [key, null]),
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
export const useTemplateTestActions = (props, { testFormContext, testDataFormContext }) => {
  const templateTestRequestData = ref(null);

  const { t } = useI18n();
  const popups = usePopups();
  const modals = useModals();
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

    const isValid = await Promise.all([
      testFormContext.validateForSubmit(),
      testDataFormContext.validateForSubmit(),
    ]).then(results => results.every(Boolean));

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
    const isValid = await testDataFormContext.validateForSubmit();

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
export const useTemplateTestRunner = (props, { mainFormContext, testDataFormContext }) => {
  const testResult = ref({});

  const {
    validateEntityServicesVariables,
    validateEventFilterRulesVariables,
    validateLinkRulesVariables,
    validateScenariosVariables,
    validateWidgetsVariables,
    validateDeclareTicketRulesVariables,
    validateDynamicInfosVariables,
    validateInstructionsVariables,
    validateJobsVariables,
    validateMetaAlarmRulesVariables,
    validateWebhookTokenRulesVariables,
  } = useTemplateTesting();

  const validateHandler = computed(
    () =>
      ({
        [TEMPLATE_TESTING_TEST_TYPES.eventFilter]: validateEventFilterRulesVariables,
        [TEMPLATE_TESTING_TEST_TYPES.linkRule]: validateLinkRulesVariables,
        [TEMPLATE_TESTING_TEST_TYPES.scenario]: validateScenariosVariables,
        [TEMPLATE_TESTING_TEST_TYPES.widget]: validateWidgetsVariables,
        [TEMPLATE_TESTING_TEST_TYPES.declareTicketRule]: validateDeclareTicketRulesVariables,
        [TEMPLATE_TESTING_TEST_TYPES.dynamicInfo]: validateDynamicInfosVariables,
        [TEMPLATE_TESTING_TEST_TYPES.instruction]: validateInstructionsVariables,
        [TEMPLATE_TESTING_TEST_TYPES.job]: validateJobsVariables,
        [TEMPLATE_TESTING_TEST_TYPES.metaAlarmRule]: validateMetaAlarmRulesVariables,
        [TEMPLATE_TESTING_TEST_TYPES.externalAuthToken]: validateWebhookTokenRulesVariables,
      })[props.type] ?? validateEntityServicesVariables,
  );

  const formToRequest = computed(
    () =>
      ({
        [TEMPLATE_TESTING_TEST_TYPES.eventFilter]: formToEventFilter,
        [TEMPLATE_TESTING_TEST_TYPES.linkRule]: formToLinkRule,
        [TEMPLATE_TESTING_TEST_TYPES.scenario]: formToScenario,
        [TEMPLATE_TESTING_TEST_TYPES.widget]: formToWidget,
        [TEMPLATE_TESTING_TEST_TYPES.declareTicketRule]: formToDeclareTicketRule,
        [TEMPLATE_TESTING_TEST_TYPES.dynamicInfo]: formToDynamicInfo,
        [TEMPLATE_TESTING_TEST_TYPES.instruction]: formToRemediationInstructionRequest,
        [TEMPLATE_TESTING_TEST_TYPES.job]: formToRemediationJob,
        [TEMPLATE_TESTING_TEST_TYPES.metaAlarmRule]: formToMetaAlarmRule,
        [TEMPLATE_TESTING_TEST_TYPES.externalAuthToken]: formToExternalAuthToken,
      })[props.type] ?? formToService,
  );

  const createRunTestHandler = (
    getValidationFormData,
    onValidationStateChange,
    onLastRunVariablesChange,
  ) => {
    const { pending: running, handler: runTest } = usePendingHandler(async () => {
      unref(testDataFormContext)?.resetErrors?.();

      const [isMainFormValid, isTestDataFormValid] = await Promise.all([
        unref(mainFormContext)?.validateForSubmit?.() ?? true,
        unref(testDataFormContext)?.validateForSubmit?.() ?? true,
      ]);
      const isValid = isMainFormValid && isTestDataFormValid;

      if (onValidationStateChange) {
        onValidationStateChange(isMainFormValid);
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

        unref(testDataFormContext)?.setFormErrors?.(
          getTemplateFieldsValidationErrors(testResult.value, props.variablesFields),
        );

        const newLastRunVariables = props.variablesFields.reduce((acc, field) => {
          const { key } = field;

          acc[key] = get(props.generalForm, key);

          return acc;
        }, {});

        if (onLastRunVariablesChange) {
          onLastRunVariablesChange(newLastRunVariables);
        }
      } catch (err) {
        const mainFormErrors = getTemplateRuleValidationErrors(err, props.generalForm);

        if (!isEmpty(mainFormErrors) && onValidationStateChange) {
          onValidationStateChange(false);
        }

        unref(mainFormContext)?.setFormErrors?.(mainFormErrors);
        unref(testDataFormContext)?.setFormErrors?.(err);
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
