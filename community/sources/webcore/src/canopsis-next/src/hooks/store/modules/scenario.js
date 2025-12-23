import { useStoreModuleHooks } from '@/hooks/store';

const useScenarioStoreModule = () => useStoreModuleHooks('scenario');

/**
 * Hook for accessing scenario store module actions
 *
 * @returns {Object} An object containing scenario actions
 */
export const useScenario = () => {
  const { useActions } = useScenarioStoreModule();

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
    ...actions,
  };
};
