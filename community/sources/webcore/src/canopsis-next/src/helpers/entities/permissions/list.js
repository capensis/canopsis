import { PERMISSIONS_TYPES_TO_ACTIONS, CRUD_ACTIONS } from '@/constants';

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
