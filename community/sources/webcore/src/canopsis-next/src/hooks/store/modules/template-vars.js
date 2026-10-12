import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Hook to access the Vuex store module for template variables.
 * This function utilizes the `useStoreModuleHooks` to provide a modular way to interact with
 * the `template/vars` Vuex module.
 *
 * @returns {Object} An object containing the store, module, and functions to access getters and actions for
 * the `template/vars` module.
 */
const useTemplateVarsStoreModule = () => useStoreModuleHooks('template/vars');

/**
 * Custom hook for accessing and interacting with the template variables Vuex store module.
 * This hook provides a convenient way to access the module's getters and actions for both
 * generic template variables and entity-specific template variables.
 *
 * @returns {Object} An object containing the mapped getters and actions for the template vars module.
 *
 * @property {Function} templateVars - Getter for accessing the template variables items.
 * @property {Function} templateVarsPending - Getter for checking if the template variables data is pending.
 * @property {Function} fetchTemplateVarsList - Action for fetching the generic list of template variables.
 * @property {Function} fetchEntityServicesTemplateVarsWithoutStore - Action for fetching entity services
 *                                                            template variables.
 * @property {Function} fetchEventFiltersTemplateVarsWithoutStore - Action for fetching event filters
 *                                                            template variables.
 * @property {Function} fetchScenariosTemplateVarsWithoutStore - Action for fetching scenarios template variables.
 * @property {Function} fetchLinkRulesTemplateVarsWithoutStore - Action for fetching link rules template variables.
 * @property {Function} fetchWidgetsTemplateVarsWithoutStore - Action for fetching widgets template variables.
 * @property {Function} fetchDeclareTicketRulesTemplateVarsWithoutStore - Action for fetching declare ticket rules
 *                                                            template variables.
 * @property {Function} fetchDynamicInfosTemplateVarsWithoutStore - Action for fetching dynamic infos
 *                                                            template variables.
 * @property {Function} fetchInstructionsTemplateVarsWithoutStore - Action for fetching instructions template variables.
 * @property {Function} fetchJobsTemplateVarsWithoutStore - Action for fetching jobs template variables.
 * @property {Function} fetchMetaAlarmRulesTemplateVarsWithoutStore - Action for fetching meta alarm rules
 *                                                            template variables.
 */
export const useTemplateVars = () => {
  const { useGetters, useActions } = useTemplateVarsStoreModule();

  const getters = useGetters({
    templateVars: 'items',
    templateVarsPending: 'pending',
  });

  const actions = useActions({
    fetchTemplateVarsList: 'fetchList',
    fetchEntityServicesTemplateVarsWithoutStore: 'fetchEntityServicesVarsWithoutStore',
    fetchEventFiltersTemplateVarsWithoutStore: 'fetchEventFiltersVarsWithoutStore',
    fetchScenariosTemplateVarsWithoutStore: 'fetchScenariosVarsWithoutStore',
    fetchLinkRulesTemplateVarsWithoutStore: 'fetchLinkRulesVarsWithoutStore',
    fetchWidgetsTemplateVarsWithoutStore: 'fetchWidgetsVarsWithoutStore',
    fetchDeclareTicketRulesTemplateVarsWithoutStore: 'fetchDeclareTicketRulesVarsWithoutStore',
    fetchDynamicInfosTemplateVarsWithoutStore: 'fetchDynamicInfosVarsWithoutStore',
    fetchInstructionsTemplateVarsWithoutStore: 'fetchInstructionsVarsWithoutStore',
    fetchJobsTemplateVarsWithoutStore: 'fetchJobsVarsWithoutStore',
    fetchMetaAlarmRulesTemplateVarsWithoutStore: 'fetchMetaAlarmRulesVarsWithoutStore',
    fetchWebhookTokenRulesTemplateVarsWithoutStore: 'fetchWebhookTokenRulesVarsWithoutStore',
  });

  return {
    ...getters,
    ...actions,
  };
};
