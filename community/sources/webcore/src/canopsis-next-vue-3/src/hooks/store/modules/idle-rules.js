import { storeToRefs } from 'pinia';

import { useIdleRulesStore } from '@/stores/modules/idle-rules';

/**
 * Idle rules store facade.
 *
 * @returns {Object}
 */
export const useIdleRules = () => {
  const idleRulesStore = useIdleRulesStore();
  const {
    items,
    meta,
    pending,
    fetchingParams,
  } = storeToRefs(idleRulesStore);

  return {
    idleRules: items,
    idleRulesMeta: meta,
    idleRulesPending: pending,
    idleRulesFetchingParams: fetchingParams,

    fetchIdleRulesList: idleRulesStore.fetchList,
    fetchIdleRulesListWithPreviousParams: idleRulesStore.fetchListWithPreviousParams,
    fetchIdleRulesListWithoutStore: idleRulesStore.fetchListWithoutStore,
    createIdleRule: idleRulesStore.create,
    updateIdleRule: idleRulesStore.update,
    removeIdleRule: idleRulesStore.remove,
    bulkEnableIdleRules: idleRulesStore.bulkEnable,
    bulkDisableIdleRules: idleRulesStore.bulkDisable,
    bulkRemoveIdleRules: idleRulesStore.bulkRemove,
  };
};
