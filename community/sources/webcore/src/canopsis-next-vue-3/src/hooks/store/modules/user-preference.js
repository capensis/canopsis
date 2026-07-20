import { storeToRefs } from 'pinia'

import { useUserPreferenceStore } from '@/stores/modules/user-preference'

/**
 * User preference store facade.
 *
 * @returns {Object}
 */
export const useUserPreference = () => {
  const userPreferenceStore = useUserPreferenceStore()
  const { items, itemsByWidget, pending, getItemByWidgetId } = storeToRefs(userPreferenceStore)

  return {
    userPreferences: items,
    userPreferencesByWidget: itemsByWidget,
    userPreferencesPending: pending,
    getUserPreferenceByWidgetId: getItemByWidgetId,

    fetchUserPreference: userPreferenceStore.fetchItem,
    fetchUserPreferenceWithoutStore: userPreferenceStore.fetchItemWithoutStore,
    updateUserPreference: userPreferenceStore.update,
    updateLocalUserPreference: userPreferenceStore.updateLocal,
  }
}
