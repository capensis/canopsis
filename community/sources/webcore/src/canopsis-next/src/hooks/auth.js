import { CRUD_ACTIONS } from '@/constants';

import { checkUserAccess } from '@/helpers/entities/permissions/list';

import { useStoreModuleHooks } from './store';

/**
 * Provides hooks for accessing the Vuex store module specifically for authentication.
 * This function utilizes the `useStoreModuleHooks` to interact with the `auth` namespace of the Vuex store.
 *
 * @returns {Object} An object containing the store, module, and functions to access getters and actions related
 * to authentication.
 */
export const useAuthStoreModule = () => useStoreModuleHooks('auth');

/**
 * Custom hook for accessing and checking the current user's permissions.
 * This hook utilizes the Vuex store module `auth` to retrieve permissions and provides methods to check specific
 * CRUD actions.
 *
 * @returns {Object} An object containing methods to check user permissions for various CRUD actions.
 *
 * @example
 * // Example of using useCurrentUserPermissions in a Vue component
 * import { useCurrentUserPermissions } from `./path/to/useCurrentUserPermissions`;
 *
 * export default {
 *   setup() {
 *     const {
 *       checkAccess,
 *       checkCreateAccess,
 *       checkReadAccess,
 *       checkUpdateAccess,
 *       checkDeleteAccess,
 *     } = useCurrentUserPermissions();
 *
 *     // Check if the user has `can` access for permission ID `123`
 *     const canAccess = checkAccess(`123`);
 *     // Check if the user has `create` access for permission ID `123`
 *     const canCreate = checkCreateAccess(`123`);
 *
 *     return {
 *       canAccess,
 *       canCreate
 *     };
 *   }
 * }
 */
export const useCurrentUserPermissions = () => {
  const { useGetters } = useAuthStoreModule();
  const { currentUserPermissionsById } = useGetters(['currentUserPermissionsById']);

  const checkAccess = (permissionId, action = CRUD_ACTIONS.can) => (
    checkUserAccess(currentUserPermissionsById.value[permissionId], action)
  );

  const checkCreateAccess = permissionId => (
    checkUserAccess(currentUserPermissionsById.value[permissionId], CRUD_ACTIONS.create)
  );

  const checkReadAccess = permissionId => (
    checkUserAccess(currentUserPermissionsById.value[permissionId], CRUD_ACTIONS.read)
  );

  const checkUpdateAccess = permissionId => (
    checkUserAccess(currentUserPermissionsById.value[permissionId], CRUD_ACTIONS.update)
  );

  const checkDeleteAccess = permissionId => (
    checkUserAccess(currentUserPermissionsById.value[permissionId], CRUD_ACTIONS.delete)
  );

  return {
    currentUserPermissionsById,
    checkAccess,
    checkCreateAccess,
    checkReadAccess,
    checkUpdateAccess,
    checkDeleteAccess,
  };
};
