import { storeToRefs } from 'pinia';

import { useMetaAlarmRuleStore } from '@/stores/modules/meta-alarm-rule';

/**
 * Meta alarm rule store hook.
 *
 * @returns {Object}
 */
export const useMetaAlarmRule = () => {
  const metaAlarmRuleStore = useMetaAlarmRuleStore();
  const { items, meta, pending, fetchingParams } = storeToRefs(metaAlarmRuleStore);

  return {
    metaAlarmRules: items,
    metaAlarmRulesMeta: meta,
    metaAlarmRulesPending: pending,
    metaAlarmRulesFetchingParams: fetchingParams,

    fetchMetaAlarmRulesList: metaAlarmRuleStore.fetchList,
    fetchMetaAlarmRulesListWithPreviousParams: metaAlarmRuleStore.fetchListWithPreviousParams,
    fetchMetaAlarmRulesListWithoutStore: metaAlarmRuleStore.fetchListWithoutStore,
    createMetaAlarmRule: metaAlarmRuleStore.create,
    updateMetaAlarmRule: metaAlarmRuleStore.update,
    removeMetaAlarmRule: metaAlarmRuleStore.remove,
    bulkEnableMetaAlarmRules: metaAlarmRuleStore.bulkEnable,
    bulkDisableMetaAlarmRules: metaAlarmRuleStore.bulkDisable,
    bulkRemoveMetaAlarmRules: metaAlarmRuleStore.bulkRemove,
    checkPatternsAlarmsCount: metaAlarmRuleStore.checkPatternsAlarmsCount,
  };
};
