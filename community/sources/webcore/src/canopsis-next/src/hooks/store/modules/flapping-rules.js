import { useStoreModuleHooks } from '@/hooks/store';

const useFlappingRulesStoreModule = () => useStoreModuleHooks('flappingRules');

/**
 * Hook for accessing flapping rules store module actions
 *
 * @returns {Object} An object containing flapping rules actions
 */
export const useFlappingRules = () => {
  const { useActions } = useFlappingRulesStoreModule();

  const actions = useActions({
    fetchFlappingRulesList: 'fetchList',
    fetchFlappingRulesListWithPreviousParams: 'fetchListWithPreviousParams',
    fetchFlappingRulesListWithoutStore: 'fetchListWithoutStore',
    createFlappingRule: 'create',
    updateFlappingRule: 'update',
    removeFlappingRule: 'remove',
    bulkEnableFlappingRules: 'bulkEnable',
    bulkDisableFlappingRules: 'bulkDisable',
    bulkRemoveFlappingRules: 'bulkRemove',
  });

  return {
    ...actions,
  };
};
