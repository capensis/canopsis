import { useStoreModuleHooks } from '@/hooks/store';

const useIdleRulesStoreModule = () => useStoreModuleHooks('idleRules');

/**
 * Hook for accessing idle rules store module actions
 *
 * @returns {Object} An object containing idle rules actions
 */
export const useIdleRules = () => {
  const { useActions } = useIdleRulesStoreModule();

  const actions = useActions({
    fetchIdleRulesList: 'fetchList',
    fetchIdleRulesListWithPreviousParams: 'fetchListWithPreviousParams',
    fetchIdleRulesListWithoutStore: 'fetchListWithoutStore',
    createIdleRule: 'create',
    updateIdleRule: 'update',
    removeIdleRule: 'remove',
    bulkEnableIdleRules: 'bulkEnable',
    bulkDisableIdleRules: 'bulkDisable',
    bulkRemoveIdleRules: 'bulkRemove',
  });

  return {
    ...actions,
  };
};
