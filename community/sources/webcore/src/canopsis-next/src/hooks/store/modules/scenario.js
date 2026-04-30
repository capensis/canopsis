import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates hooks for accessing the scenario Vuex store module.
 * Provides access to getters and actions for managing scenario operations.
 *
 * @returns {Object} An object containing store module utilities:
 * @property {import('vuex').Store} store - The Vuex store instance
 * @property {import('vuex').Module} module - The scenario module instance
 * @property {Function} useGetters - Function to access module getters
 * @property {Function} useActions - Function to access module actions
 */
const useScenarioStoreModule = () => useStoreModuleHooks('scenario');

/**
 * Custom hook for scenario operations.
 * Provides convenient access to scenario getters and actions.
 *
 * @returns {Object} An object containing scenario getters and actions:
 * @property {Object} scenariosMeta - Getter for scenario metadata
 * @property {boolean} scenariosPending - Getter for scenario loading state
 * @property {Array} scenarios - Getter for scenario items
 * @property {Function} fetchScenariosList - Action to fetch scenarios list
 * @property {Function} createScenario - Action to create a scenario
 * @property {Function} updateScenario - Action to update a scenario
 * @property {Function} removeScenario - Action to remove a scenario
 * @property {Function} createTestScenarioExecution - Action to create test scenario execution
 * @property {Function} fetchTestScenarioExecutionWithoutStore - Action to fetch test scenario execution without store
 * @property {Function} fetchTestScenarioExecutionWebhooksResponse
 *                      - Action to fetch test scenario execution webhooks response
 * @property {Function} bulkEnableScenarios - Action to bulk enable scenarios
 * @property {Function} bulkDisableScenarios - Action to bulk disable scenarios
 * @property {Function} bulkRemoveScenarios - Action to bulk remove scenarios
 * @property {Function} fetchScenariosListWithPreviousParams - Action to fetch scenarios list with previous params
 * @property {Function} fetchScenariosListWithoutStore - Action to fetch scenarios list without store
 *
 * @example
 * // Usage in a component
 * const { scenarios, scenariosPending, fetchScenariosList } = useScenario();
 * await fetchScenariosList({ params: { page: 1, limit: 10 } });
 */
export const useScenario = () => {
  const { useGetters, useActions } = useScenarioStoreModule();

  const getters = useGetters({
    scenariosMeta: 'meta',
    scenariosPending: 'pending',
    scenarios: 'items',
  });

  const actions = useActions({
    fetchScenariosList: 'fetchList',
    fetchScenariosListWithPreviousParams: 'fetchListWithPreviousParams',
    fetchScenariosListWithoutStore: 'fetchListWithoutStore',
    createScenario: 'create',
    updateScenario: 'update',
    removeScenario: 'remove',
    createTestScenarioExecution: 'createTestScenarioExecution',
    fetchTestScenarioExecutionWithoutStore: 'fetchTestScenarioExecutionWithoutStore',
    fetchTestScenarioExecutionWebhooksResponse: 'fetchTestScenarioExecutionWebhooksResponse',
    bulkEnableScenarios: 'bulkEnable',
    bulkDisableScenarios: 'bulkDisable',
    bulkRemoveScenarios: 'bulkRemove',
  });

  return {
    ...getters,
    ...actions,
  };
};
