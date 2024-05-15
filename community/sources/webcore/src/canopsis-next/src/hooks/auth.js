import { computed } from 'vue';

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

/**
 * A Vue composition function that checks if the current user has access based on a specified permission.
 * This function utilizes the `useCurrentUserPermissions` hook to retrieve the necessary checkAccess method,
 * and then uses Vue's computed property to reactively handle the permission check.
 *
 * @param {string} permission - The permission ID to check access against.
 * @returns {Object} An object containing a reactive `hasAccess` property that indicates whether the user has the
 * specified permission.
 *
 * @example
 * // Example of using useCanPermission in a Vue component
 * import { useCanPermission } from `./path/to/useCanPermission`;
 *
 * export default {
 *   setup() {
 *     const { hasAccess } = useCanPermission(`123`);
 *     return { hasAccess };
 *   }
 * }
 */
export const useCanPermission = (permission) => {
  const { checkAccess } = useCurrentUserPermissions();

  const hasAccess = computed(() => checkAccess(permission));

  return { hasAccess };
};

/**
 * Provides reactive Vue computed properties to check CRUD permissions for a specific permission ID.
 * This function utilizes the `useCurrentUserPermissions` hook to access check functions for create, read, update, and
 * delete permissions.
 *
 * @param {string} permission - The permission ID to check access against.
 * @returns {Object} An object containing Vue computed properties for each CRUD operation's access status.
 * @example
 * // Example of using useCRUDPermissions in a Vue component
 * import { useCRUDPermissions } from `./path/to/useCRUDPermissions`;
 *
 * export default {
 *   setup() {
 *     const permissionId = `123`;
 *     const {
 *       hasCreateAccess,
 *       hasReadAccess,
 *       hasUpdateAccess,
 *       hasDeleteAccess,
 *     } = useCRUDPermissions(permissionId);
 *
 *     return {
 *       hasCreateAccess,
 *       hasReadAccess,
 *       hasUpdateAccess,
 *       hasDeleteAccess,
 *     };
 *   }
 * }
 */
export const useCRUDPermissions = (permission) => {
  const {
    checkCreateAccess,
    checkReadAccess,
    checkUpdateAccess,
    checkDeleteAccess,
  } = useCurrentUserPermissions();

  const hasCreateAccess = computed(() => checkCreateAccess(permission));
  const hasReadAccess = computed(() => checkReadAccess(permission));
  const hasUpdateAccess = computed(() => checkUpdateAccess(permission));
  const hasDeleteAccess = computed(() => checkDeleteAccess(permission));

  return {
    hasCreateAccess,
    hasReadAccess,
    hasUpdateAccess,
    hasDeleteAccess,
  };
};
