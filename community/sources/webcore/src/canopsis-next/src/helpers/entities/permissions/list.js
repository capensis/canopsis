import { difference } from 'lodash';

import { PERMISSIONS_TYPES_TO_ACTIONS, CRUD_ACTIONS, USER_PERMISSIONS_GROUPS } from '@/constants';

/**
 * Check user access for a permission
 *
 * @param {Object} permission
 * @param {string} action
 * @returns {boolean}
 */
export const checkUserAccess = (permission, action) => {
  if (permission && permission.actions) {
    const { actions } = permission;

    return action === CRUD_ACTIONS.can
      ? actions.length >= 0
      : actions.includes(action);
  }

  return false;
};

/**
 * Check user access for a permission
 *
 * @param {Object} currentUserPermissionsById
 * @param {Object} permissions
 * @param {string} action
 * @returns {boolean}
 */
export const checkUserAnyAccess = (permissions, action) => (
  Object.values(permissions).some(permission => checkUserAccess(permission, action))
);

/**
 * Get actions for permission by type
 *
 * @param {Object} permission
 * @returns {*}
 */
export const getPermissionActions = permission => (PERMISSIONS_TYPES_TO_ACTIONS[permission.type]
  ? PERMISSIONS_TYPES_TO_ACTIONS[permission.type]
  : [CRUD_ACTIONS.can]);

export const permissionsToTreeview = (permissions = []) => permissions.reduce((acc, permission) => {
  let activeGroup = acc;

  permission.groups.forEach((group) => {
    if (!activeGroup[group._id]) {
      activeGroup[group._id] = {
        ...group,

        children: {},
        allChildren: [],
        actions: [CRUD_ACTIONS.can],
      };
    }
    activeGroup[group._id].allChildren.push({ _id: permission._id, actions: getPermissionActions(permission) });
    activeGroup = activeGroup[group._id].children;
  });

  activeGroup[permission._id] = {
    ...permission,

    actions: getPermissionActions(permission),
  };

  return acc;
}, {});

/**
 * Get properties for a permission checkbox based on role and permission details.
 *
 * @param {Object} role - The role object containing permissions.
 * @param {Object} permission - The permission object which may contain child permissions.
 * @param {string} action - The specific action to check within the permission.
 * @returns {Object} An object containing properties for the checkbox:
 * - `inputValue`: {boolean} Indicates if the checkbox should be checked.
 * - `indeterminate`: {boolean} Indicates if the checkbox should be in an indeterminate state
 * (only present if `permission.allChildren` exists).
 */
export const getPermissionCheckboxProps = (role, permission, action) => {
  if (permission.allChildren) {
    const childrenPermissionsDiffs = permission.allChildren
      .map(({ _id: id, actions }) => difference(actions ?? [], role.permissions[id] ?? []).length);
    const inputValue = childrenPermissionsDiffs.every(v => !v);
    const hasCheckedChildren = permission.allChildren.some(({ _id: id }) => !!role.permissions[id]?.length);

    return {
      inputValue,
      indeterminate: !inputValue && hasCheckedChildren,
    };
  }

  return {
    inputValue: role.permissions[permission._id]?.includes(action) ?? false,
  };
};

// TODO: add jsdoc
export const isApiPermission = (permissionId = '') => permissionId.startsWith(USER_PERMISSIONS_GROUPS.api);
