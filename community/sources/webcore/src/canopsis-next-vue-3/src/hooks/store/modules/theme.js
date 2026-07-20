import { storeToRefs } from 'pinia'

import { useThemeStore } from '@/stores/modules/theme'

/**
 * Theme store facade.
 *
 * @returns {Object}
 */
export const useTheme = () => {
  const themeStore = useThemeStore()
  const { items, meta, pending, fetchingParams } = storeToRefs(themeStore)

  return {
    themes: items,
    themesMeta: meta,
    themesPending: pending,
    themesFetchingParams: fetchingParams,

    fetchThemesList: themeStore.fetchList,
    fetchThemesListWithPreviousParams: themeStore.fetchListWithPreviousParams,
    fetchThemesListWithoutStore: themeStore.fetchListWithoutStore,
    createTheme: themeStore.create,
    updateTheme: themeStore.update,
    removeTheme: themeStore.remove,
    bulkRemoveThemes: themeStore.bulkRemove,
  }
}
