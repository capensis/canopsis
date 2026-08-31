/**
 * Converts an array of permissions into a single string by joining them with a delimiter
 *
 * @param {Array<string>} [permissions=[]] - Array of permission strings to be joined
 * @param {string} [delimiter=''] - The separator to use between permissions
 * @returns {string} - The joined permissions string
 */
export const groupedPermissionToPermission = (permissions = [], delimiter = '') => permissions.join(delimiter);

const CAN_ACTION = 'can';

export const checkUserAccess = (permission, action) => {
  if (!permission?.actions) {
    return false;
  }

  if (Array.isArray(permission.actions)) {
    if (action === CAN_ACTION) {
      return true;
    }

    return permission.actions.includes(action);
  }

  if (permission[action] !== undefined) {
    return Boolean(permission[action]);
  }

  return Boolean(permission.actions?.[action] ?? permission.permissions?.[action]);
};

export const checkUserAnyAccess = (permissions, action) => Object.values(permissions).some(permission => checkUserAccess(permission, action));

/**
 * Groups permissions by their view property into an object structure.
 *
 * @param {Array} [permissions=[]] - Array of permission objects. Each permission object should have 'view'
 * and 'name' properties.
 * @returns {Object} An object where:
 *   - Keys are the view names from the permissions
 *   - Values are objects containing permissions as key-value pairs, where:
 *     - Keys are permission names
 *     - Values are the original permission objects
 *
 * @example
 * const permissions = [
 *   { view: 'users', name: 'read', ... },
 *   { view: 'users', name: 'write', ... },
 *   { view: 'posts', name: 'delete', ... }
 * ];
 * const grouped = viewPermissionsGroupedPermissions(permissions);
 * // Result:
 * // {
 * //   users: { read: {...}, write: {...} },
 * //   posts: { delete: {...} }
 * // }
 */
export const viewPermissionsGroupedPermissions = (permissions = []) => permissions.reduce((acc, permission) => {
    if (!permission.view) {
      return acc;
    }

    if (!acc[permission.view]) {
      acc[permission.view] = {};
    }

    acc[permission.view][permission.name] = permission;

    return acc;
  }, {});
