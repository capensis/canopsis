import { storeToRefs } from 'pinia';

import { useFlappingRulesStore } from '@/stores/modules/flapping-rules';

/**
 * Flapping rules store facade.
 *
 * @returns {Object}
 */
export const useFlappingRules = () => {
  const flappingRulesStore = useFlappingRulesStore();
  const {
    items,
    meta,
    pending,
    fetchingParams,
  } = storeToRefs(flappingRulesStore);

  return {
    flappingRules: items,
    flappingRulesMeta: meta,
    flappingRulesPending: pending,
    flappingRulesFetchingParams: fetchingParams,

    fetchFlappingRulesList: flappingRulesStore.fetchList,
    fetchFlappingRulesListWithPreviousParams: flappingRulesStore.fetchListWithPreviousParams,
    fetchFlappingRulesListWithoutStore: flappingRulesStore.fetchListWithoutStore,
    createFlappingRule: flappingRulesStore.create,
    updateFlappingRule: flappingRulesStore.update,
    removeFlappingRule: flappingRulesStore.remove,
    bulkEnableFlappingRules: flappingRulesStore.bulkEnable,
    bulkDisableFlappingRules: flappingRulesStore.bulkDisable,
    bulkRemoveFlappingRules: flappingRulesStore.bulkRemove,
  };
};
