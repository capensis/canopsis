import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Hook for accessing the icon Vuex store module
 * Creates hooks for accessing icon module's getters and actions
 *
 * @returns {Object} An object containing store module utilities
 * @property {Object} store - The Vuex store instance
 * @property {Object} module - The icon Vuex module
 * @property {Function} useGetters - Function to access icon module getters
 * @property {Function} useActions - Function to access icon module actions
 * @example
 * const { useGetters, useActions } = useIconStoreModule();
 * const { icons } = useGetters(['icons']);
 * const { fetchIconsList } = useActions(['fetchIconsList']);
 */
const useIconStoreModule = () => useStoreModuleHooks('icon');

/**
 * Hook for managing icon-related operations and state
 *
 * @returns {Object} An object containing icon-related getters and actions
 * @property {Array} icons - List of icons
 * @property {boolean} iconsPending - Loading state of icons
 * @property {Object} iconsMeta - Metadata for icons
 * @property {Function} fetchIconsList - Fetches the list of icons
 * @property {Function} fetchIconsListWithPreviousParams - Fetches icons list with previous parameters
 * @property {Function} createIcon - Creates a new icon
 * @property {Function} updateIcon - Updates an existing icon
 * @property {Function} removeIcon - Removes an icon
 */
export const useIcon = () => {
  const { useGetters, useActions } = useIconStoreModule();

  const getters = useGetters({
    icons: 'items',
    iconsPending: 'pending',
    iconsMeta: 'meta',
  });

  const actions = useActions({
    fetchIconsList: 'fetchList',
    fetchIconsListWithPreviousParams: 'fetchListWithPreviousParams',
    createIcon: 'create',
    updateIcon: 'update',
    removeIcon: 'remove',
  });

  return {
    ...getters,
    ...actions,
  };
};
