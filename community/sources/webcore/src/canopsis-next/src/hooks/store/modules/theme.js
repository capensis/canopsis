import { useStoreModuleHooks } from '@/hooks/store';

const useThemeStoreModule = () => useStoreModuleHooks('theme');

/**
 * Hook for managing theme-related operations and state
 *
 * @returns {Object} An object containing theme-related getters and actions
 * @property {Array} themes - List of themes
 * @property {boolean} themesPending - Loading state of themes
 * @property {Object} themesMeta - Metadata for themes
 * @property {Function} fetchThemesList - Fetches the list of themes
 * @property {Function} fetchThemesListWithPreviousParams - Fetches themes list with previous parameters
 * @property {Function} fetchThemesListWithoutStore - Fetches themes list without updating store
 * @property {Function} createTheme - Creates a new theme
 * @property {Function} updateTheme - Updates an existing theme
 * @property {Function} removeTheme - Removes a theme
 * @property {Function} bulkRemoveThemes - Removes multiple themes
 */
export const useTheme = () => {
  const { useGetters, useActions } = useThemeStoreModule();

  const getters = useGetters({
    themes: 'items',
    themesPending: 'pending',
    themesMeta: 'meta',
  });

  const actions = useActions({
    fetchThemesList: 'fetchList',
    fetchThemesListWithPreviousParams: 'fetchListWithPreviousParams',
    fetchThemesListWithoutStore: 'fetchListWithoutStore',
    createTheme: 'create',
    updateTheme: 'update',
    removeTheme: 'remove',
    bulkRemoveThemes: 'bulkRemove',
  });

  return {
    ...getters,
    ...actions,
  };
};
