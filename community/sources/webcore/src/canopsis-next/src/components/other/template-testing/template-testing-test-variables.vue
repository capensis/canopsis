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
import { get, mapValues, isEmpty } from 'lodash';
import {
  computed,
  ref,
  watch,
  inject,
  onMounted,
  onBeforeUnmount,
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
    const afterSubmitObserver = inject('$afterSubmitObserver');

    let templateTestRequestData = {};

    const { t } = useI18n();
    const popups = usePopups();

    const selectedTest = ref({});
    const validationForm = ref([]);
    const testResult = ref({});
    const lastRunVariables = ref({});
    const isGeneralFormValid = ref(true);

    const validator = useValidator();
    const { setFormErrors } = useValidationFormErrors(validationForm);

    const { createTemplateTest, updateTemplateTest } = useTemplateTest();

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

    watch(() => props.generalForm, (newForm) => {
      const newValidationForm = formToTemplateTestingTestValidateForm(newForm, props.type);

      const { added, removed } = getChangesForValidateForm(newValidationForm, validationForm.value);

      validationForm.value = validationForm.value
        .filter(item => !removed.some(removedItem => removedItem.key === item.key));
      validationForm.value.push(...added);

      validationForm.value.forEach((_, index) => validationForm.value[index].index = newValidationForm[index]?.index);
    }, { immediate: true });

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

    const modals = useModals();

    const save = async () => {
      const data = {
        ...selectedTest.value,
        data: formToTemplateTestingTestValidate(validationForm.value),
      };

      if (props.ruleId) {
        await updateTemplateTest({ data, id: selectedTest.value._id });

        popups.success({ text: t('templateTesting.testSaved') });

        return;
      }

      templateTestRequestData = data;

      popups.success({ text: t('templateTesting.testWillSaveAfterFormSaving') });
    };

    const saveAsNew = async () => {
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
            data: formToTemplateTestingTestValidate(validationForm.value),
            description: '',
            name: newTestName,
            type: props.type,
          };

          if (props.ruleId) {
            data.rule = props.ruleId;

            const newTest = await createTemplateTest({ data });

            popups.success({ text: t('templateTesting.testSaved') });

            selectedTest.value = newTest;

            return;
          }

          templateTestRequestData = data;

          selectedTest.value = {
            ...data,
            _id: uid('test'),
          };

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

    const { pending: running, handler: runTest } = usePendingHandler(async () => {
      validator.errors.clear('test-data');
      validator.errors.clear('test-result');

      const [isParentValid, isValid] = await Promise.all([
        validator.validateAll(),
        validator.validateAll('test-data'),
      ]);

      isGeneralFormValid.value = isParentValid;

      if (!isParentValid || !isValid) {
        return;
      }

      try {
        testResult.value = await validateHandler.value({
          data: {
            rule: formToRequest.value(props.generalForm),
            testdata: formToTemplateTestingTestValidate(validationForm.value),
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

        lastRunVariables.value = props.variablesFields.reduce((acc, field) => {
          const key = field.formKey ?? field.key;

          acc[key] = get(props.generalForm, key);

          return acc;
        }, {});
      } catch (err) {
        const { mainFormErrors, testDataErrors } = Object.entries(err).reduce((acc, [key, value]) => {
          if (key.startsWith('rule.')) {
            acc.mainFormErrors[key.replace('rule.', '')] = value;
          } else {
            acc.testDataErrors[key] = value;
          }

          return acc;
        }, { mainFormErrors: {}, testDataErrors: {} });

        if (!isEmpty(mainFormErrors)) {
          isGeneralFormValid.value = false;
        }

        setFormErrors(mainFormErrors);
        setFormErrors(testDataErrors, 'test-data');
      }
    });

    const validationFormInput = () => {
      lastRunVariables.value = mapValues(lastRunVariables.value, () => Symbol('lastRunVariables'));
    };

    const afterSubmit = async (rule) => {
      /**
       * Capture the translation string before the component is destroyed
       * vue-i18n provides a proxy object that is destroyed when the component is destroyed
       * so we need to capture the translation string before the component is destroyed
       */
      const testSavedMessage = t('templateTesting.testSaved');

      if (rule) {
        templateTestRequestData.rule = rule._id;
      }

      await createTemplateTest({ data: templateTestRequestData });

      popups.success({ text: testSavedMessage });
    };

    const updateSelectedTest = (newSelectedTest) => {
      selectedTest.value = newSelectedTest;

      if (templateTestRequestData) {
        templateTestRequestData = null;
      }

      validationForm.value = templateTestingTestValidateToForm(validationForm.value, newSelectedTest.data);
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
