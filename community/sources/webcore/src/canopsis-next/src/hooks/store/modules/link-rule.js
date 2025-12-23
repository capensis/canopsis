import { useStoreModuleHooks } from '@/hooks/store';

const useLinkRuleStoreModule = () => useStoreModuleHooks('linkRule');

/**
 * Hook for accessing link rules store module actions
 *
 * @returns {Object} An object containing link rules actions
 */
export const useLinkRule = () => {
  const { useActions } = useLinkRuleStoreModule();

  const actions = useActions({
    fetchLinkRulesList: 'fetchList',
    fetchLinkRulesListWithPreviousParams: 'fetchListWithPreviousParams',
    fetchLinkRulesListWithoutStore: 'fetchListWithoutStore',
    createLinkRule: 'create',
    updateLinkRule: 'update',
    removeLinkRule: 'remove',
    bulkRemoveLinkRules: 'bulkRemove',
    fetchLinkCategoriesWithoutStore: 'fetchLinkCategoriesWithoutStore',
    bulkEnableLinkRules: 'bulkEnable',
    bulkDisableLinkRules: 'bulkDisable',
  });

  return {
    ...actions,
  };
};
