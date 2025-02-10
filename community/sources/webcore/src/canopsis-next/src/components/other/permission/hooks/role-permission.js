import {
  keyBy,
  omit,
  isEqual,
  filter,
  cloneDeep,
  sortBy,
} from 'lodash';
import { computed, ref, unref, set } from 'vue';

import { API_USER_PERMISSIONS_ROOT_GROUPS, MAX_LIMIT, MODALS, ROLE_TYPES } from '@/constants';

import { permissionsToTreeview } from '@/helpers/entities/permissions/list';
import { formToRolePermissions, roleToPermissionForm } from '@/helpers/entities/role/form';
import { mapIds } from '@/helpers/array';

import { useModals } from '@/hooks/modals';
import { useRole } from '@/hooks/store/modules/role';
import { usePendingHandler } from '@/hooks/query/pending';
import { usePermissions } from '@/hooks/store/modules/permissions';

/**
 * Hook for managing role permissions fetching and manipulation.
 * Provides functionality to fetch, update, and track changes to role permissions.
 *
 * @param {Object} [options = {}] - Configuration options for the hook
 * @param {Ref<string>} [options.activeTab] - Reference to the currently active permissions tab
 *
 * @returns {Object} An object containing:
 *   - pending {import('vue').Ref<boolean>} - Indicates if any async operation is in progress
 *   - roles {import('vue').ComputedRef<Array>} - Computed array of roles filtered by active tab
 *   - treeviewPermissions {import('vue').ComputedRef<Array>} - Computed hierarchical permissions structure
 *   - hasChanges {import('vue').ComputedRef<boolean>} - Indicates if there are unsaved permission changes
 *   - resetRolesById {Function} - Resets roles to their original state
 *   - updateRoles {Function} - Updates modified roles' permissions in bulk
 *   - changeRole {Function} - Updates permissions for a specific role
 *   - fetchList {Function} - Fetches roles and permissions data
 */
export const useRolePermissionFetching = ({ activeTab } = {}) => {
  const originalRolesById = ref({});
  const rolesById = ref({});
  const changedRoles = ref({});
  const permissions = ref([]);

  const { fetchPermissionsListWithoutStore } = usePermissions();
  const { fetchRolesListWithoutStore, bulkUpdateRolePermissions } = useRole();

  /**
   * Resets the roles state by creating a deep clone of the original roles and clearing changed roles
   */
  const resetRolesById = () => {
    rolesById.value = cloneDeep(originalRolesById.value);
    changedRoles.value = {};
  };

  const { pending, handler: fetchList } = usePendingHandler(async () => {
    const [rolesResponse, permissionsResponse] = await Promise.all([
      fetchRolesListWithoutStore({ params: { limit: MAX_LIMIT, with_flags: true } }),
      fetchPermissionsListWithoutStore({ params: { limit: MAX_LIMIT } }),
    ]);

    originalRolesById.value = keyBy(Object.values(rolesResponse.data).map(roleToPermissionForm), '_id');
    permissions.value = permissionsResponse.data;
    resetRolesById();
  });

  const treeviewPermissions = computed(() => (
    sortBy(Object.values(permissionsToTreeview(permissions.value)), 'position')
  ));
  const isApiPermissionsTab = computed(() => API_USER_PERMISSIONS_ROOT_GROUPS.includes(unref(activeTab)));
  const uiRoles = computed(() => filter(Object.values(rolesById.value), ['type', ROLE_TYPES.ui]));
  const apiRoles = computed(() => filter(Object.values(rolesById.value), ['type', ROLE_TYPES.api]));
  const roles = computed(() => (isApiPermissionsTab.value ? apiRoles.value : uiRoles.value));
  const hasChanges = computed(() => (
    Object.entries(changedRoles.value).some(([roleId, rolePermissions]) => (
      Object.keys(rolePermissions).some(permissionId => (
        !isEqual(
          rolesById.value[roleId].permissions[permissionId],
          originalRolesById.value[roleId].permissions[permissionId],
        )
      ))
    ))
  ));

  /**
   * Sets a permission change for a specific role in the changedRoles object.
   * If the role or permission doesn't exist in the tracking object, it creates the necessary structure.
   *
   * @param {string} roleId - The identifier of the role to be modified
   * @param {string} permissionId - The identifier of the permission to be set
   */
  const setChangedRolePermission = (roleId, permissionId) => {
    if (!changedRoles.value[roleId]) {
      set(changedRoles.value, roleId, { [permissionId]: true });
    }

    if (!changedRoles.value[roleId][permissionId]) {
      set(changedRoles.value[roleId], permissionId, true);
    }
  };

  /**
   * Sets multiple permission changes for a specific role by iterating through an array of permission IDs.
   *
   * @param {string|number} roleId - The identifier of the role to be modified
   * @param {Array<string|number>} permissionsIds - An array of permission identifiers to be set for the role
   */
  const setChangedRolePermissions = (roleId, permissionsIds) => (
    permissionsIds.forEach(permissionId => setChangedRolePermission(roleId, permissionId))
  );

  /**
   * Updates role permissions and their associated actions in the role-based access control system.
   * Handles both individual permissions and permissions with children (hierarchical permissions).
   *
   * @param {boolean} value - Flag indicating whether to add (true) or remove (false) the permission/action
   * @param {Object} role - The role object to be modified
   * @param {string|number} role._id - Unique identifier for the role
   * @param {Object} role.permissions - Current permissions map for the role
   * @param {Object} permission - The permission object to be added or removed
   * @param {string|number} permission._id - Unique identifier for the permission
   * @param {Array} [permission.allChildren] - Optional array of child permissions
   * @param {string} action - The specific action to be added or removed from the permission
   */
  const changeRole = (value, role, permission, action) => {
    if (!changedRoles.value[role._id]) {
      set(changedRoles.value, role._id, { [permission._id]: true });
    }

    if (!changedRoles.value[role._id][permission._id]) {
      set(changedRoles.value, permission._id, true);
    }

    /**
     * Handle permissions with children (hierarchical permissions)
     */
    if (permission.allChildren) {
      const allChildrenIds = mapIds(permission.allChildren);
      const newPermissions = value
        ? {
          ...role.permissions,
          ...permission.allChildren.reduce((acc, { _id: id, actions }) => {
            acc[id] = actions;

            return acc;
          }, {}),
        }
        : omit(role.permissions, allChildrenIds);

      setChangedRolePermissions(role._id, allChildrenIds);

      set(rolesById.value[role._id], 'permissions', newPermissions);

      return;
    }

    /**
     * Handle individual permission changes
     */
    setChangedRolePermission(role._id, permission._id);

    const currentActions = role.permissions[permission._id] ?? [];

    const newActions = value
      ? [...currentActions, action]
      : currentActions.filter(currentAction => currentAction !== action);

    if (!newActions.length) {
      set(rolesById.value[role._id], 'permissions', omit(role.permissions, permission._id));

      return;
    }

    set(rolesById.value[role._id].permissions, permission._id, newActions);
  };

  /**
   * Updates multiple roles' permissions in bulk and refreshes the roles list.
   * Processes only the roles that have been modified (tracked in changedRoles).
   *
   * @returns {Promise<void>} A promise that resolves when roles are updated and the list is refreshed
   */
  const updateRoles = async () => {
    const rolesForUpdate = Object.keys(changedRoles.value)
      .map(roleId => rolesById.value[roleId] && formToRolePermissions(rolesById.value[roleId]))
      .filter(Boolean);

    await bulkUpdateRolePermissions({ data: rolesForUpdate });

    return fetchList();
  };

  return {
    pending,
    roles,
    treeviewPermissions,
    hasChanges,
    resetRolesById,
    updateRoles,
    changeRole,
    fetchList,
  };
};

/**
 * Hook for managing role permission actions with modal confirmations.
 *
 * @param {Object} [options = {}] - Configuration options for the hook
 * @param {Function} [options.updateRoles] - Function to update modified roles' permissions
 * @param {Function} [options.resetRolesById] - Function to reset roles to their original state
 *
 * @returns {Object} An object containing:
 *   - submit {Function} - Shows confirmation modal and executes updateRoles when confirmed
 *   - cancel {Function} - Shows confirmation modal and executes resetRolesById when confirmed
 */
export const useRolePermissionActions = ({ updateRoles, resetRolesById } = {}) => {
  const modals = useModals();

  /**
   * Shows a confirmation modal dialog before executing the role permissions update.
   * When confirmed, triggers the updateRoles process through the modal's action handler.
   */
  const submit = () => modals.show({
    name: MODALS.confirmation,
    config: {
      action: updateRoles,
    },
  });

  /**
   * Shows a confirmation modal dialog before resetting the roles to their original state.
   * When confirmed, triggers the resetRolesById function through the modal's action handler.
   */
  const cancel = () => modals.show({
    name: MODALS.confirmation,
    config: {
      action: resetRolesById,
    },
  });

  return {
    submit,
    cancel,
  };
};
