import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates hooks for accessing the permissions Vuex store module.
 * Provides access to getters and actions within the 'permission' namespace.
 *
 * @returns {Object} An object containing:
 *   - store: The Vuex store instance
 *   - module: The permissions module instance
 *   - useGetters: Function to access module getters
 *   - useActions: Function to access module actions
 * @throws {Error} If the permissions module is not registered in the store
 */
const usePermissionsStoreModule = () => useStoreModuleHooks('permission');

/**
 * Hook for accessing permission-related actions from the Vuex store.
 * Provides access to permission actions within the 'permission' namespace.
 *
 * @returns {Object} An object containing:
 *   - fetchPermissionsListWithoutStore: Function to fetch permissions list without storing in Vuex
 */
export const usePermissions = () => {
  const { useActions } = usePermissionsStoreModule();

  const actions = useActions({
    fetchPermissionsListWithoutStore: 'fetchListWithoutStore',
  });

  return {
    ...actions,
  };
};
