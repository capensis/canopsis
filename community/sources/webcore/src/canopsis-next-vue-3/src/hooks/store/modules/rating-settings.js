import { storeToRefs } from 'pinia'

import { useRatingSettingsStore } from '@/stores/modules/rating-settings'

/**
 * Rating settings store facade.
 *
 * @returns {Object}
 */
export const useRatingSettings = () => {
  const ratingSettingsStore = useRatingSettingsStore()
  const { items, meta, pending, fetchingParams, updatedAt } = storeToRefs(ratingSettingsStore)

  return {
    ratingSettings: items,
    ratingSettingsMeta: meta,
    ratingSettingsPending: pending,
    ratingSettingsFetchingParams: fetchingParams,
    ratingSettingsUpdatedAt: updatedAt,

    fetchRatingSettingsList: ratingSettingsStore.fetchList,
    fetchRatingSettingsListWithPreviousParams: ratingSettingsStore.fetchListWithPreviousParams,
    fetchRatingSettingsListWithoutStore: ratingSettingsStore.fetchListWithoutStore,
    createRatingSetting: ratingSettingsStore.create,
    updateRatingSetting: ratingSettingsStore.update,
    removeRatingSetting: ratingSettingsStore.remove,
    bulkUpdateRatingSettings: ratingSettingsStore.bulkUpdate,
  }
}
