import { storeToRefs } from 'pinia'

import { useCorporatePatternStore } from '@/stores/modules/pattern-corporate'

/**
 * Corporate pattern store facade.
 *
 * @returns {Object}
 */
export const useCorporatePattern = () => {
  const corporatePatternStore = useCorporatePatternStore()
  const { items, meta, pending, fetchingParams } = storeToRefs(corporatePatternStore)

  return {
    corporatePatterns: items,
    corporatePatternsMeta: meta,
    corporatePatternsPending: pending,
    corporatePatternsFetchingParams: fetchingParams,

    fetchCorporatePatternsList: corporatePatternStore.fetchList,
    fetchCorporatePatternsListWithPreviousParams: corporatePatternStore.fetchListWithPreviousParams,
    createCorporatePattern: corporatePatternStore.create,
    updateCorporatePattern: corporatePatternStore.update,
    removeCorporatePattern: corporatePatternStore.remove,
  }
}
