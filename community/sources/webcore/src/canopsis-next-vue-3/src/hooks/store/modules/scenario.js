import { storeToRefs } from 'pinia';

import { useScenarioStore } from '@/stores/modules/scenario';

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
  const scenarioStore = useScenarioStore();
  const {
    items,
    meta,
    pending,
    fetchingParams,
  } = storeToRefs(scenarioStore);

  return {
    scenariosMeta: meta,
    scenariosPending: pending,
    scenarios: items,
    scenariosFetchingParams: fetchingParams,

    fetchScenariosList: scenarioStore.fetchList,
    fetchScenariosListWithPreviousParams: scenarioStore.fetchListWithPreviousParams,
    fetchScenariosListWithoutStore: scenarioStore.fetchListWithoutStore,
    createScenario: scenarioStore.create,
    updateScenario: scenarioStore.update,
    removeScenario: scenarioStore.remove,
    createTestScenarioExecution: scenarioStore.createTestScenarioExecution,
    fetchTestScenarioExecutionWithoutStore: scenarioStore.fetchTestScenarioExecutionWithoutStore,
    fetchTestScenarioExecutionWebhooksResponse: scenarioStore.fetchTestScenarioExecutionWebhooksResponse,
    bulkEnableScenarios: scenarioStore.bulkEnable,
    bulkDisableScenarios: scenarioStore.bulkDisable,
    bulkRemoveScenarios: scenarioStore.bulkRemove,
  };
};
