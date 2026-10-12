import { cloneDeep, filter, isEqual, keyBy, omit, set, sortBy } from 'lodash';
import { computed, ref } from 'vue';

import {
  CONDITIONAL_PERMISSIONS_MAP,
  CRUD_ACTIONS,
  MAX_LIMIT,
  MODALS,
  ROLE_TYPES,
  USER_PERMISSIONS_GROUPS,
} from '@/constants';

import { mapIds } from '@/helpers/array';
import { permissionsToTreeview } from '@/helpers/entities/permissions/list';
import { formToRolePermissions, roleToPermissionForm } from '@/helpers/entities/role/form';

import { useModals } from '@/hooks/modals';
import { usePendingHandler } from '@/hooks/query/pending';
import { usePermissions } from '@/hooks/store/modules/permissions';
import { useRole } from '@/hooks/store/modules/role';

/**
 * Hook for managing role permissions fetching and manipulation.
 *
 * @returns {Object}
 */
export const useRolePermissionFetching = () => {
  const originalRolesById = ref({});
  const rolesById = ref({});
  const changedRoles = ref({});
  const permissions = ref([]);

  const { fetchPermissionsListWithoutStore } = usePermissions();
  const { fetchRolesListWithoutStore, bulkUpdateRolePermissions } = useRole();

  const resetRolesById = () => {
    rolesById.value = cloneDeep(originalRolesById.value);
    changedRoles.value = {};
  };

  const { pending, handler: fetchList } = usePendingHandler(async () => {
    const [rolesResponse, permissionsResponse] = await Promise.all([
      fetchRolesListWithoutStore({ params: { limit: MAX_LIMIT, with_flags: true } }),
      fetchPermissionsListWithoutStore({ params: { limit: MAX_LIMIT } }),
    ]);

    originalRolesById.value = keyBy(
      Object.values(rolesResponse.data).map(roleToPermissionForm),
      '_id',
    );
    permissions.value = permissionsResponse.data;
    resetRolesById();
  });

  const treeviewPermissions = computed(() =>
    sortBy(Object.values(permissionsToTreeview(permissions.value)), 'position'),
  );

  const uiRoles = computed(() => filter(Object.values(rolesById.value), ['type', ROLE_TYPES.ui]));
  const apiRoles = computed(() => filter(Object.values(rolesById.value), ['type', ROLE_TYPES.api]));

  const hasChanges = computed(() =>
    Object.entries(changedRoles.value).some(([roleId, rolePermissions]) =>
      Object.keys(rolePermissions).some(
        permissionId =>
          !isEqual(
            rolesById.value[roleId].permissions[permissionId],
            originalRolesById.value[roleId].permissions[permissionId],
          ),
      ),
    ),
  );

  const setChangedRolePermission = (roleId, permissionId) => {
    if (!changedRoles.value[roleId]) {
      set(changedRoles.value, roleId, { [permissionId]: true });
    }

    if (!changedRoles.value[roleId][permissionId]) {
      set(changedRoles.value[roleId], permissionId, true);
    }
  };

  const setChangedRolePermissions = (roleId, permissionsIds) => {
    permissionsIds.map(permissionId => setChangedRolePermission(roleId, permissionId));
  };

  const enableDependentPermissions = (role, permissionId, value) => {
    if (!value) {
      return;
    }

    const dependencies = CONDITIONAL_PERMISSIONS_MAP[permissionId];

    if (!dependencies?.length) {
      return;
    }

    const crudActions = [
      CRUD_ACTIONS.create,
      CRUD_ACTIONS.read,
      CRUD_ACTIONS.update,
      CRUD_ACTIONS.delete,
    ];

    dependencies.forEach(({ dependentPermission }) => {
      setChangedRolePermission(role._id, dependentPermission);
      set(rolesById.value[role._id].permissions, dependentPermission, crudActions);
    });
  };

  const changeRole = (value, role, permission, action) => {
    setChangedRolePermission(role._id, permission._id);

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

    const currentActions = role.permissions[permission._id] ?? [];
    const newActions = value
      ? [...currentActions, action]
      : currentActions.filter(currentAction => currentAction !== action);

    if (!newActions.length) {
      set(rolesById.value[role._id], 'permissions', omit(role.permissions, permission._id));
    } else {
      set(rolesById.value[role._id].permissions, permission._id, newActions);
    }

    enableDependentPermissions(role, permission._id, value && newActions.length > 0);
  };

  const updateRoles = async () => {
    const rolesForUpdate = Object.keys(changedRoles.value)
      .map(roleId => rolesById.value[roleId] && formToRolePermissions(rolesById.value[roleId]))
      .filter(Boolean);

    await bulkUpdateRolePermissions({ data: rolesForUpdate });

    return fetchList();
  };

  return {
    pending,
    uiRoles,
    apiRoles,
    treeviewPermissions,
    hasChanges,
    resetRolesById,
    updateRoles,
    changeRole,
    fetchList,
    apiTabId: USER_PERMISSIONS_GROUPS.api,
  };
};

/**
 * Hook for managing role permission actions with modal confirmations.
 *
 * @param {Object} [options={}]
 * @param {Function} [options.updateRoles]
 * @param {Function} [options.resetRolesById]
 * @returns {{ submit: Function, cancel: Function }}
 */
export const useRolePermissionActions = ({ updateRoles, resetRolesById } = {}) => {
  const modals = useModals();

  const submit = () =>
    modals.show({
      name: MODALS.confirmation,
      config: {
        action: updateRoles,
      },
    });

  const cancel = () =>
    modals.show({
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
