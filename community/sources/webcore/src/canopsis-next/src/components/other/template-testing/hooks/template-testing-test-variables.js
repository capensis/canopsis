import { get, isEmpty } from 'lodash';
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
 * @returns {Object} Validation form state and methods
 */
export const useTemplateTestingValidationForm = (props) => {
  const validationForm = ref([]);
  const selectedTest = ref({});
  const lastRunVariables = ref({});

  const { setFormErrors } = useValidationFormErrors(validationForm);

  // Watch for changes in general form and update validation form accordingly
  watch(() => props.generalForm, (newForm) => {
    const newValidationForm = formToTemplateTestingTestValidateForm(newForm, props.type);

    const { added, removed } = getChangesForValidateForm(newValidationForm, validationForm.value);

    validationForm.value = validationForm.value
      .filter(item => !removed.some(removedItem => removedItem.key === item.key));
    validationForm.value.push(...added);

    validationForm.value.forEach((_, index) => {
      if (newValidationForm[index]) {
        validationForm.value[index] = {
          ...validationForm.value[index],
          index: newValidationForm[index].index,
        };
      }
    });
  }, { immediate: true });

  const getValidationFormData = () => formToTemplateTestingTestValidate(validationForm.value);

  const updateSelectedTest = (newSelectedTest) => {
    selectedTest.value = newSelectedTest;
    validationForm.value = templateTestingTestValidateToForm(validationForm.value, newSelectedTest.data);
  };

  const resetLastRunVariables = () => {
    lastRunVariables.value = Object.fromEntries(
      Object.keys(lastRunVariables.value).map(key => [key, Symbol('lastRunVariables')]),
    );
  };

  const setLastRunVariables = (variables) => {
    lastRunVariables.value = { ...variables };
  };

  return {
    validationForm,
    selectedTest,
    lastRunVariables,
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

  const { createTemplateTest, updateTemplateTest } = useTemplateTest();

  const save = async (selectedTest, getValidationFormData) => {
    const data = {
      ...selectedTest,
      data: getValidationFormData(),
    };

    if (props.ruleId) {
      await updateTemplateTest({ data, id: selectedTest._id });

      popups.success({ text: t('templateTesting.testSaved') });

      return;
    }

    templateTestRequestData.value = data;

    popups.success({ text: t('templateTesting.testWillSaveAfterFormSaving') });
  };

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

  const getTemplateTestRequestData = () => templateTestRequestData.value;

  const createTemplateTestFromRequestData = async () => {
    if (!templateTestRequestData.value) {
      return null;
    }

    const newTest = await createTemplateTest({ data: templateTestRequestData.value });
    templateTestRequestData.value = null;

    return newTest;
  };

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
 * @returns {Object} Test runner state and methods
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
          const key = field.formKey ?? field.key;

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
 * Hook for managing after submit observer logic
 *
 * @param {Object} params - Parameters
 * @param {Function} params.getTemplateTestRequestData - Function to get template test request data
 * @param {Object} props - Component props
 * @returns {Object} After submit methods
 */
export const useTemplateTestAfterSubmit = ({ getTemplateTestRequestData }, props) => {
  const afterSubmitObserver = inject('$afterSubmitObserver');

  const { t } = useI18n();
  const popups = usePopups();

  const { createTemplateTest } = useTemplateTest();

  const createAfterSubmitHandler = () => {
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

    return afterSubmit;
  };

  const setupAfterSubmitObserver = () => {
    const afterSubmit = createAfterSubmitHandler();

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

    return afterSubmit;
  };

  return {
    createAfterSubmitHandler,
    setupAfterSubmitObserver,
  };
};
