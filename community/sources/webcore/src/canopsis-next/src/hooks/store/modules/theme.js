import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates hooks for accessing the theme Vuex store module
 *
 * @returns {Object} Store module hooks for theme namespace
 */
const useThemeStoreModule = () => useStoreModuleHooks('theme');

/**
 * Hook for managing theme-related operations and state
 *
 * @returns {Object} An object containing theme getters, actions, and methods
 */
export const useTheme = () => {
  const { useGetters, useActions } = useThemeStoreModule();

  const getters = useGetters({
    themes: 'items',
    themesMeta: 'meta',
    themesPending: 'pending',
  });

  const actions = useActions({
    fetchThemesList: 'fetchList',
    fetchThemesListWithPreviousParams: 'fetchListWithPreviousParams',
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
