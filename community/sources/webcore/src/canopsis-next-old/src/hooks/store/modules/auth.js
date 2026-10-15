import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates hooks for accessing the auth Vuex store module.
 * Provides access to getters and actions within the 'auth' namespace.
 *
 * @returns {Object} An object containing:
 *   - store: The Vuex store instance
 *   - module: The auth module instance
 *   - useGetters: Function to access module getters
 *   - useActions: Function to access module actions
 */
export const useAuthStoreModule = () => useStoreModuleHooks('auth');
