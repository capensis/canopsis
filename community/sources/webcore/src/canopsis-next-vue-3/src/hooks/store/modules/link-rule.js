import { storeToRefs } from 'pinia'

import { useLinkRuleStore } from '@/stores/modules/link-rule'

/**
 * Link rule store facade.
 *
 * @returns {Object}
 */
export const useLinkRule = () => {
  const linkRuleStore = useLinkRuleStore()
  const { items, meta, pending, fetchingParams } = storeToRefs(linkRuleStore)

  return {
    linkRules: items,
    linkRulesMeta: meta,
    linkRulesPending: pending,
    linkRulesFetchingParams: fetchingParams,

    fetchLinkRulesList: linkRuleStore.fetchList,
    fetchLinkRulesListWithPreviousParams: linkRuleStore.fetchListWithPreviousParams,
    fetchLinkRulesListWithoutStore: linkRuleStore.fetchListWithoutStore,
    fetchLinkCategoriesWithoutStore: linkRuleStore.fetchLinkCategoriesWithoutStore,
    createLinkRule: linkRuleStore.create,
    updateLinkRule: linkRuleStore.update,
    removeLinkRule: linkRuleStore.remove,
    bulkRemoveLinkRules: linkRuleStore.bulkRemove,
    bulkEnableLinkRules: linkRuleStore.bulkEnable,
    bulkDisableLinkRules: linkRuleStore.bulkDisable,
  }
}
