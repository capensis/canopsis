import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates hooks for accessing the role Vuex store module.
 * Provides a wrapper around useStoreModuleHooks specifically for the role module namespace.
 *
 * @returns {Object} An object containing store module utilities:
 * @returns {Object} store - The Vuex store instance
 * @returns {Object} module - The role Vuex module instance
 * @returns {Function} useGetters - Function to access role module getters
 * @returns {Function} useActions - Function to access role module actions
 */
const useRoleStoreModule = () => useStoreModuleHooks('role');

/**
 * Creates hooks for accessing the role Vuex store module's getters and actions.
 * This composable function provides a centralized way to interact with role-related store operations.
 *
 * @returns {Object} An object containing role-related getters and actions:
 * @returns {Array} roles - List of all roles
 * @returns {Function} getRoleById - Function to get a role by its ID
 * @returns {boolean} rolesPending - Loading state of role operations
 * @returns {Object} rolesMeta - Meta information about roles
 * @returns {Function} fetchRolesListWithoutStore - Fetch roles without storing them
 * @returns {Function} fetchRolesList - Fetch and store roles
 * @returns {Function} removeRole - Remove a role
 * @returns {Function} createRole - Create a new role
 * @returns {Function} updateRole - Update an existing role
 */
export const useRole = () => {
  const { useGetters, useActions } = useRoleStoreModule();

  const getters = useGetters({
    roles: 'items',
    getRoleById: 'getItemById',
    rolesPending: 'pending',
    rolesMeta: 'meta',
  });

  const actions = useActions({
    fetchRolesListWithoutStore: 'fetchListWithoutStore',
    fetchRolesList: 'fetchList',
    removeRole: 'remove',
    createRole: 'create',
    updateRole: 'update',
  });

  return {
    ...getters,
    ...actions,
  };
};
