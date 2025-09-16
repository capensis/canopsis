import { isEqual, get, mapValues } from 'lodash';
import { computed, ref, unref, watch } from 'vue';

import { TEMPLATE_TESTING_TEST_TYPES, LINK_RULE_TYPES } from '@/constants';

import { templateVarsToVariables } from '@/helpers/variables';
import { convertRuleToTemplateTestingTestValidateForm } from '@/helpers/entities/template-testing-validate/form';

import { useI18n } from '@/hooks/i18n';
import { useTemplateVars } from '@/hooks/store/modules/template-vars';
import { usePendingHandler } from '@/hooks/query/pending';
import { useModelField } from '@/hooks/form/model-field';

/**
 * Hook for managing template variables list based on test type
 *
 * @param {Object} options - Configuration options
 * @param {import('vue').Ref<string>|string} [options.type] - The type of template testing
 *                                                            (eventFilter, linkRule, scenario, etc.)
 * @param {import('vue').Ref<Object>|Object} [options.form] - The form object
 *
 * @returns {Object} Hook return object
 * @returns {import('vue').Ref<Array>} returns.templateVars - Reactive array of template variables
 * @returns {import('vue').Ref<boolean>} returns.pending - Loading state indicator
 * @returns {Function} returns.fetchList - Function to fetch template variables list
 */
export const useTemplateVarsList = ({ type, form } = {}) => {
  const { t } = useI18n();

  const templateVars = ref({});

  const {
    fetchEntityServicesTemplateVarsWithoutStore,
    fetchEventFiltersTemplateVarsWithoutStore,
    fetchLinkRulesTemplateVarsWithoutStore,
    fetchScenariosTemplateVarsWithoutStore,
    fetchWidgetsTemplateVarsWithoutStore,
    fetchDeclareTicketRulesTemplateVarsWithoutStore,
    fetchDynamicInfosTemplateVarsWithoutStore,
    fetchInstructionsTemplateVarsWithoutStore,
    fetchJobsTemplateVarsWithoutStore,
    fetchMetaAlarmRulesTemplateVarsWithoutStore,
  } = useTemplateVars();

  const fetchHandler = computed(() => ({
    [TEMPLATE_TESTING_TEST_TYPES.eventFilter]: fetchEventFiltersTemplateVarsWithoutStore,
    [TEMPLATE_TESTING_TEST_TYPES.linkRule]: async () => {
      const { alarm: alarmData, entity: entityData } = await fetchLinkRulesTemplateVarsWithoutStore();

      return unref(form).type === LINK_RULE_TYPES.alarm ? alarmData : entityData;
    },
    [TEMPLATE_TESTING_TEST_TYPES.scenario]: fetchScenariosTemplateVarsWithoutStore,
    [TEMPLATE_TESTING_TEST_TYPES.widget]: fetchWidgetsTemplateVarsWithoutStore,
    [TEMPLATE_TESTING_TEST_TYPES.declareTicketRule]: fetchDeclareTicketRulesTemplateVarsWithoutStore,
    [TEMPLATE_TESTING_TEST_TYPES.dynamicInfo]: fetchDynamicInfosTemplateVarsWithoutStore,
    [TEMPLATE_TESTING_TEST_TYPES.instruction]: fetchInstructionsTemplateVarsWithoutStore,
    [TEMPLATE_TESTING_TEST_TYPES.job]: fetchJobsTemplateVarsWithoutStore,
    [TEMPLATE_TESTING_TEST_TYPES.metaAlarmRule]: fetchMetaAlarmRulesTemplateVarsWithoutStore,
  }[unref(type)] ?? fetchEntityServicesTemplateVarsWithoutStore));

  const addTranslationToTemplateVars = (items = []) => items.map((item) => {
    const result = {
      ...item,
      text: t(`templateTesting.templateVars.${item.text}`),
    };

    if (item.variables?.length) {
      result.variables = addTranslationToTemplateVars(item.variables);
    }

    return result;
  });

  const {
    pending,
    handler: fetchList,
  } = usePendingHandler(async () => {
    if (!fetchHandler.value) {
      return;
    }

    const data = await fetchHandler.value();

    templateVars.value = mapValues(templateVarsToVariables(data), addTranslationToTemplateVars);
  });

  return {
    templateVars,
    pending,
    fetchList,
  };
};

/**
 * Hook for managing test variables fields with form integration
 *
 * @param {Object} [props={}] - Component props containing form data
 * @param {Object} props.form - Form object containing test variables data
 * @param {import('vue').Ref<string>|string} type - The type of template testing
 * @param {Function} emit - Vue emit function for component communication
 * @returns {Object} Hook return object
 * @returns {import('vue').Ref<Array>} returns.items - Reactive array of test variable field items with update handlers
 */
export const useTestVariablesFields = (props = {}, type, emit) => {
  const items = ref([]);

  const { updateField } = useModelField(props, emit);

  const isEmptyItems = computed(() => !items.value.length);

  watch(() => props.form, (newForm) => {
    const newItems = convertRuleToTemplateTestingTestValidateForm(newForm, unref(type));

    if (isEqual(items.value, newItems)) {
      return;
    }

    items.value = newItems.map(newItem => ({
      ...newItem,
      value: get(newForm, newItem.key),
      updateField: value => updateField(newItem.key, value),
    }));
  }, { immediate: true });

  return {
    items,
    isEmptyItems,
  };
};
