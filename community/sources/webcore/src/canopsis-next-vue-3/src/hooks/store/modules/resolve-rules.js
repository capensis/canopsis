import { storeToRefs } from 'pinia';

import { useResolveRulesStore } from '@/stores/modules/resolve-rules';

/**
 * Resolve rules store facade.
 *
 * @returns {Object}
 */
export const useResolveRules = () => {
  const resolveRulesStore = useResolveRulesStore();
  const {
    items,
    meta,
    pending,
    fetchingParams,
  } = storeToRefs(resolveRulesStore);

  return {
    resolveRules: items,
    resolveRulesMeta: meta,
    resolveRulesPending: pending,
    resolveRulesFetchingParams: fetchingParams,

    fetchResolveRulesList: resolveRulesStore.fetchList,
    fetchResolveRulesListWithPreviousParams: resolveRulesStore.fetchListWithPreviousParams,
    fetchResolveRulesListWithoutStore: resolveRulesStore.fetchListWithoutStore,
    createResolveRule: resolveRulesStore.create,
    updateResolveRule: resolveRulesStore.update,
    removeResolveRule: resolveRulesStore.remove,
    bulkEnableResolveRules: resolveRulesStore.bulkEnable,
    bulkDisableResolveRules: resolveRulesStore.bulkDisable,
    bulkRemoveResolveRules: resolveRulesStore.bulkRemove,
  };
};
