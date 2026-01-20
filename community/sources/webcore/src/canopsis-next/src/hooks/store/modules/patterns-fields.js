import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates hooks for accessing the patterns-fields Vuex store module.
 * Provides access to actions for fetching pattern fields for different rule types.
 *
 * @returns {Object} An object containing store module utilities:
 * @property {import('vuex').Store} store - The Vuex store instance
 * @property {import('vuex').Module} module - The patterns-fields module instance
 * @property {Function} useGetters - Function to access module getters
 * @property {Function} useActions - Function to access module actions
 */
const usePatternsFieldsStoreModule = () => useStoreModuleHooks('patternsFields');

/**
 * Custom hook for pattern fields operations.
 * Provides convenient access to pattern fields actions for various rule types.
 *
 * @returns {Object} An object containing pattern fields actions:
 * @property {Function} fetchFlappingRulePatternFields - Fetches pattern fields for flapping rules
 * @property {Function} fetchIdleRulePatternFields - Fetches pattern fields for idle rules
 * @property {Function} fetchLinkRulePatternFields - Fetches pattern fields for link rules
 * @property {Function} fetchRulePatternFields - Fetches pattern fields for general rules
 * @property {Function} fetchPbehaviorPatternFields - Fetches pattern fields for pbehaviors
 * @property {Function} fetchAlarmTagPatternFields - Fetches pattern fields for alarm tags
 * @property {Function} fetchWidgetFilterPatternFields - Fetches pattern fields for widget filters
 * @property {Function} fetchEntityservicePatternFields - Fetches pattern fields for entity services
 * @property {Function} fetchStateSettingPatternFields - Fetches pattern fields for state settings
 * @property {Function} fetchEventfilterPatternFields - Fetches pattern fields for event filters
 * @property {Function} fetchScenarioPatternFields - Fetches pattern fields for scenarios
 * @property {Function} fetchMetaalarmrulePatternFields - Fetches pattern fields for meta alarm rules
 * @property {Function} fetchDeclareTicketRulePatternFields - Fetches pattern fields for declare ticket rules
 * @property {Function} fetchInstructionPatternFields - Fetches pattern fields for instructions
 * @property {Function} fetchKpiFilterPatternFields - Fetches pattern fields for KPI filters
 * @property {Function} fetchDynamicInfosPatternFields - Fetches pattern fields for dynamic infos
 *
 * @example
 * // Usage in a component
 * const { fetchFlappingRulePatternFields } = usePatternsFields();
 * const fields = await fetchFlappingRulePatternFields({ params: {} });
 */
export const usePatternsFields = () => {
  const { useActions } = usePatternsFieldsStoreModule();

  const actions = useActions({
    fetchFlappingRulePatternFields: 'fetchFlappingRulePatternFields',
    fetchIdleRulePatternFields: 'fetchIdleRulePatternFields',
    fetchLinkRulePatternFields: 'fetchLinkRulePatternFields',
    fetchRulePatternFields: 'fetchRulePatternFields',
    fetchPbehaviorPatternFields: 'fetchPbehaviorPatternFields',
    fetchAlarmTagPatternFields: 'fetchAlarmTagPatternFields',
    fetchWidgetFilterPatternFields: 'fetchWidgetFilterPatternFields',
    fetchEntityservicePatternFields: 'fetchEntityservicePatternFields',
    fetchStateSettingPatternFields: 'fetchStateSettingPatternFields',
    fetchEventfilterPatternFields: 'fetchEventfilterPatternFields',
    fetchScenarioPatternFields: 'fetchScenarioPatternFields',
    fetchMetaalarmrulePatternFields: 'fetchMetaalarmrulePatternFields',
    fetchDeclareTicketRulePatternFields: 'fetchDeclareTicketRulePatternFields',
    fetchInstructionPatternFields: 'fetchInstructionPatternFields',
    fetchKpiFilterPatternFields: 'fetchKpiFilterPatternFields',
    fetchDynamicInfosPatternFields: 'fetchDynamicInfosPatternFields',
  });

  return {
    ...actions,
  };
};
