import { useStoreModuleHooks } from '@/hooks/store';

const useResolveRulesStoreModule = () => useStoreModuleHooks('resolveRules');

/**
 * Hook for accessing resolve rules store module actions
 *
 * @returns {Object} An object containing resolve rules actions
 */
export const useResolveRules = () => {
  const { useActions } = useResolveRulesStoreModule();

  const actions = useActions({
    fetchResolveRulesList: 'fetchList',
    fetchResolveRulesListWithPreviousParams: 'fetchListWithPreviousParams',
    fetchResolveRulesListWithoutStore: 'fetchListWithoutStore',
    createResolveRule: 'create',
    updateResolveRule: 'update',
    removeResolveRule: 'remove',
    bulkEnableResolveRules: 'bulkEnable',
    bulkDisableResolveRules: 'bulkDisable',
    bulkRemoveResolveRules: 'bulkRemove',
  });

  return {
    ...actions,
  };
};
