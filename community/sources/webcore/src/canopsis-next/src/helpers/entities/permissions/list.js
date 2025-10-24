import { difference, isEmpty } from 'lodash';

import {
  PERMISSIONS_TYPES_TO_ACTIONS,
  CRUD_ACTIONS,
  USER_PERMISSIONS_GROUPS,
  INVERSE_CONDITIONAL_PERMISSIONS_MAP,
} from '@/constants';

import { filterObjectEntries, collectLeafKeys } from '@/helpers/collection';

/**
 * Check user access for a permission
 *
 * @param {Object} permission
 * @param {string} action
 * @returns {boolean}
 */
export const checkUserAccess = (permission, action) => {
  // Early return if no permission or actions
  if (!permission?.actions) {
    return false;
  }

  // Early return for 'can' action
  if (action === CRUD_ACTIONS.can) {
    return permission.actions.length >= 0;
  }

  return permission.actions.includes(action);
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
export const getPermissionActions = (permission) => {
  // Early return if no permission
  if (!permission?.type) {
    return [CRUD_ACTIONS.can];
  }

  return PERMISSIONS_TYPES_TO_ACTIONS[permission.type] || [CRUD_ACTIONS.can];
};

/**
 * Converts a flat permissions array into a hierarchical tree structure
 *
 * @param {Array} [permissions=[]] - Array of permission objects containing groups and actions
 * @returns {Object} Tree structure where:
 *   - Each node contains group/permission data
 *   - Groups have `children`, `allChildren`, and `actions` properties
 *   - Leaf nodes (permissions) contain the permission data and `actions`
 *   - Group nodes use group._id as keys
 *   - Permission nodes use permission._id as keys
 *
 * @example
 * const permissions = [{
 *   _id: 'perm1',
 *   groups: [{ _id: 'group1', name: 'Group 1' }],
 *   type: 'crud'
 * }];
 * const tree = permissionsToTreeview(permissions);
 * // Returns:
 * // {
 * //   group1: {
 * //     _id: 'group1',
 * //     name: 'Group 1',
 * //     children: {
 * //       perm1: { _id: 'perm1', actions: ['create', 'read', 'update', 'delete'] }
 * //     },
 * //     allChildren: [{ _id: 'perm1', actions: ['create', 'read', 'update', 'delete'] }],
 * //     actions: ['can']
 * //   }
 * // }
 */
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
 * Checks if a permission should be disabled and provides tooltip key based on conditional dependencies
 *
 * @param {Object} role - The role object containing permissions
 * @param {string} permissionId - The permission ID to check
 * @returns {Object|null} Object with disabled flag, inputValue, and tooltip key, or null if no conditions apply
 */
export const getConditionalPermissionState = (role, permissionId) => {
  if (!role?.permissions) {
    return null;
  }

  const conditionalDependencies = INVERSE_CONDITIONAL_PERMISSIONS_MAP[permissionId];

  if (!conditionalDependencies?.length) {
    return null;
  }

  const activeConditional = conditionalDependencies.find(
    ({ triggerPermission }) => role.permissions[triggerPermission]?.length > 0,
  );

  return activeConditional
    ? {
      disabled: true,
      tooltipKey: activeConditional.tooltipKey,
      inputValue: true,
    }
    : null;
};

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
 * - `disabled`: {boolean} Indicates if the checkbox should be disabled.
 * - `tooltipKey`: {string} i18n key for tooltip text for disabled checkboxes.
 */
export const getPermissionCheckboxProps = (role, permission, action) => {
  if (!role || !permission) {
    return { inputValue: false };
  }

  const conditionalState = getConditionalPermissionState(role, permission._id);

  if (permission.allChildren) {
    const childrenPermissionsDiffs = permission.allChildren
      .map(({ _id: id, actions }) => difference(actions ?? [], role.permissions?.[id] ?? []).length);
    const inputValue = childrenPermissionsDiffs.every(v => !v);
    const hasCheckedChildren = permission.allChildren.some(({ _id: id }) => !!role.permissions?.[id]?.length);

    return {
      inputValue,
      indeterminate: !inputValue && hasCheckedChildren,
      ...(conditionalState || {}),
    };
  }

  const inputValue = conditionalState?.inputValue
    ?? role.permissions?.[permission._id]?.includes(action)
    ?? false;

  return {
    inputValue,
    ...(conditionalState || {}),
  };
};

/**
 * Checks if a permission ID belongs to the API permission group
 *
 * @param {string} [permissionId=''] - The permission ID to check
 * @returns {boolean} True if the permission belongs to API group, false otherwise
 */
export const isApiPermission = (permissionId = '') => permissionId.startsWith(USER_PERMISSIONS_GROUPS.api);

/**
 * Filters treeview permissions by search term recursively
 *
 * @param {Object} treeviewPermissions - The treeview permissions object
 * @param {string} [search=''] - The search term to filter by
 * @returns {Object} Filtered treeview permissions object. If a parent has any matching child, the parent is included.
 *
 * @example
 * const filtered = filterTreeviewPermissions(permissions, 'alarm');
 * // Returns tree with only nodes matching 'alarm' and their parents
 */
export const filterTreeviewPermissions = (treeviewPermissions, search = '') => {
  if (!search || !search.trim()) {
    return treeviewPermissions;
  }

  const searchLower = search.toLowerCase().trim();

  /**
   * Check if a node matches the search term
   *
   * @param {Object} node - The node to check
   * @param {string} key - The key of the node
   * @returns {boolean} True if the node matches the search term
   */
  const matchesSearch = (node, key) => (
    node.name?.toLowerCase().includes(searchLower)
    || node.description?.toLowerCase().includes(searchLower)
    || key.toLowerCase().includes(searchLower)
  );

  /**
   * Recursively filters a node and its children based on search criteria
   *
   * @param {Object} node - The node to filter
   * @param {string} key - The key of the node
   * @returns {Object|null} Filtered node with children, or null if no matches found
   */
  const filterNode = (node, key) => {
    if (node.children) {
      const filteredChildren = filterObjectEntries(
        node.children,
        (childNode, childKey) => filterNode(childNode, childKey),
      );

      if (isEmpty(filteredChildren)) {
        return null;
      }

      const filteredAllChildren = node.allChildren.filter(
        child => collectLeafKeys(filteredChildren).has(child._id),
      );

      return { ...node, children: filteredChildren, allChildren: filteredAllChildren };
    }

    return matchesSearch(node, key) ? node : null;
  };

  return filterObjectEntries(treeviewPermissions, (node, key) => filterNode(node, key));
};
