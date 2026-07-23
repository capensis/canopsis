import { computed, ref, unref } from 'vue';

import { LINK_RULE_TYPES, TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

import { usePendingHandler } from '@/hooks/query/pending';
import { useTemplateVars } from '@/hooks/store/modules/template-vars';

import { useVarsPrepare } from './shared';

/**
 * Hook for managing template variables list based on test type.
 *
 * @param {Object} options - Configuration options
 * @param {import('vue').Ref<string>|string} [options.type] - The type of template testing
 * @param {import('vue').Ref<Object>|Object} [options.form] - The form object
 * @returns {Object} Hook return object
 */
export const useTemplateVarsList = ({ type, form = {} } = {}) => {
  const vars = ref({});

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
    fetchWebhookTokenRulesTemplateVarsWithoutStore,
  } = useTemplateVars();

  const fetchHandler = computed(
    () =>
      ({
        [TEMPLATE_TESTING_TEST_TYPES.eventFilter]: fetchEventFiltersTemplateVarsWithoutStore,
        [TEMPLATE_TESTING_TEST_TYPES.linkRule]: async () => {
          const { alarm: alarmData, entity: entityData } =
            await fetchLinkRulesTemplateVarsWithoutStore();

          return unref(form).type === LINK_RULE_TYPES.alarm ? alarmData : entityData;
        },
        [TEMPLATE_TESTING_TEST_TYPES.scenario]: fetchScenariosTemplateVarsWithoutStore,
        [TEMPLATE_TESTING_TEST_TYPES.widget]: fetchWidgetsTemplateVarsWithoutStore,
        [TEMPLATE_TESTING_TEST_TYPES.declareTicketRule]:
          fetchDeclareTicketRulesTemplateVarsWithoutStore,
        [TEMPLATE_TESTING_TEST_TYPES.dynamicInfo]: fetchDynamicInfosTemplateVarsWithoutStore,
        [TEMPLATE_TESTING_TEST_TYPES.instruction]: fetchInstructionsTemplateVarsWithoutStore,
        [TEMPLATE_TESTING_TEST_TYPES.job]: fetchJobsTemplateVarsWithoutStore,
        [TEMPLATE_TESTING_TEST_TYPES.metaAlarmRule]: fetchMetaAlarmRulesTemplateVarsWithoutStore,
        [TEMPLATE_TESTING_TEST_TYPES.externalAuthToken]:
          fetchWebhookTokenRulesTemplateVarsWithoutStore,
      })[unref(type)] ?? fetchEntityServicesTemplateVarsWithoutStore,
  );

  const { prepare: prepareVars } = useVarsPrepare('templateTesting.templateVars');

  const { pending, handler: fetchList } = usePendingHandler(async () => {
    if (!fetchHandler.value) {
      return;
    }

    const data = await fetchHandler.value();

    vars.value = prepareVars(data);
  });

  return {
    vars,
    pending,
    fetchList,
  };
};
