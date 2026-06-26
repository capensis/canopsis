import { storeToRefs } from 'pinia';

import { useTemplateVarsStore } from '@/stores/modules/template-vars';

/**
 * Custom hook for accessing and interacting with the template variables store module.
 *
 * @returns {Object} An object containing mapped state and actions for template vars.
 */
export const useTemplateVars = () => {
  const templateVarsStore = useTemplateVarsStore();
  const { items, pending } = storeToRefs(templateVarsStore);

  return {
    templateVars: items,
    templateVarsPending: pending,

    fetchTemplateVarsList: templateVarsStore.fetchList,
    fetchEntityServicesTemplateVarsWithoutStore: templateVarsStore.fetchEntityServicesVarsWithoutStore,
    fetchEventFiltersTemplateVarsWithoutStore: templateVarsStore.fetchEventFiltersVarsWithoutStore,
    fetchScenariosTemplateVarsWithoutStore: templateVarsStore.fetchScenariosVarsWithoutStore,
    fetchLinkRulesTemplateVarsWithoutStore: templateVarsStore.fetchLinkRulesVarsWithoutStore,
    fetchWidgetsTemplateVarsWithoutStore: templateVarsStore.fetchWidgetsVarsWithoutStore,
    fetchDeclareTicketRulesTemplateVarsWithoutStore: templateVarsStore.fetchDeclareTicketRulesVarsWithoutStore,
    fetchDynamicInfosTemplateVarsWithoutStore: templateVarsStore.fetchDynamicInfosVarsWithoutStore,
    fetchInstructionsTemplateVarsWithoutStore: templateVarsStore.fetchInstructionsVarsWithoutStore,
    fetchJobsTemplateVarsWithoutStore: templateVarsStore.fetchJobsVarsWithoutStore,
    fetchMetaAlarmRulesTemplateVarsWithoutStore: templateVarsStore.fetchMetaAlarmRulesVarsWithoutStore,
    fetchWebhookTokenRulesTemplateVarsWithoutStore: templateVarsStore.fetchWebhookTokenRulesVarsWithoutStore,
  };
};
