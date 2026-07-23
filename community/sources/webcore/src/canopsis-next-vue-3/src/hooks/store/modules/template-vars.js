import { useTemplateVarsStore } from '@/stores/modules/template-vars';

/**
 * Custom hook for accessing and interacting with template variables.
 *
 * @returns {Object} Template vars actions.
 */
export const useTemplateVars = () => {
  const templateVarsStore = useTemplateVarsStore();

  return {
    fetchTemplateVarsList: templateVarsStore.fetchTemplateVarsList,
    fetchEntityServicesTemplateVarsWithoutStore:
      templateVarsStore.fetchEntityServicesTemplateVarsWithoutStore,
    fetchEventFiltersTemplateVarsWithoutStore:
      templateVarsStore.fetchEventFiltersTemplateVarsWithoutStore,
    fetchScenariosTemplateVarsWithoutStore:
      templateVarsStore.fetchScenariosTemplateVarsWithoutStore,
    fetchLinkRulesTemplateVarsWithoutStore:
      templateVarsStore.fetchLinkRulesTemplateVarsWithoutStore,
    fetchWidgetsTemplateVarsWithoutStore: templateVarsStore.fetchWidgetsTemplateVarsWithoutStore,
    fetchDeclareTicketRulesTemplateVarsWithoutStore:
      templateVarsStore.fetchDeclareTicketRulesTemplateVarsWithoutStore,
    fetchDynamicInfosTemplateVarsWithoutStore:
      templateVarsStore.fetchDynamicInfosTemplateVarsWithoutStore,
    fetchInstructionsTemplateVarsWithoutStore:
      templateVarsStore.fetchInstructionsTemplateVarsWithoutStore,
    fetchJobsTemplateVarsWithoutStore: templateVarsStore.fetchJobsTemplateVarsWithoutStore,
    fetchMetaAlarmRulesTemplateVarsWithoutStore:
      templateVarsStore.fetchMetaAlarmRulesTemplateVarsWithoutStore,
    fetchWebhookTokenRulesTemplateVarsWithoutStore:
      templateVarsStore.fetchWebhookTokenRulesTemplateVarsWithoutStore,
  };
};
