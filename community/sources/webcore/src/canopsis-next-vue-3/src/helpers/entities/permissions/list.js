import { difference } from 'lodash';

import {
  CRUD_ACTIONS,
  INVERSE_CONDITIONAL_PERMISSIONS_MAP,
  PERMISSIONS_TYPES_TO_ACTIONS,
  USER_PERMISSIONS_GROUPS,
} from '@/constants';

/**
 * Get actions for permission by type.
 *
 * @param {Object} permission
 * @returns {string[]}
 */
export const getPermissionActions = (permission) => {
  if (!permission?.type) {
    return [CRUD_ACTIONS.can];
  }

  return PERMISSIONS_TYPES_TO_ACTIONS[permission.type] || [CRUD_ACTIONS.can];
};

/**
 * Converts a flat permissions array into a hierarchical tree structure.
 *
 * @param {Array} [permissions=[]]
 * @returns {Object}
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

      activeGroup[group._id].allChildren.push({
        _id: permission._id,
        actions: getPermissionActions(permission),
      });
      activeGroup = activeGroup[group._id].children;
    });

    activeGroup[permission._id] = {
      ...permission,

      actions: getPermissionActions(permission),
    };

    return acc;
  }, {});

/**
 * Checks if a permission should be disabled and provides tooltip key based on conditional dependencies.
 *
 * @param {Object} role
 * @param {string} permissionId
 * @returns {Object|null}
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
 * @param {Object} role
 * @param {Object} permission
 * @param {string} action
 * @returns {Object}
 */
export const getPermissionCheckboxProps = (role, permission, action) => {
  if (!role || !permission) {
    return { inputValue: false };
  }

  const conditionalState = getConditionalPermissionState(role, permission._id);

  if (permission.allChildren) {
    const childrenPermissionsDiffs = permission.allChildren.map(
      ({ _id: id, actions }) => difference(actions ?? [], role.permissions?.[id] ?? []).length,
    );
    const inputValue = childrenPermissionsDiffs.every(value => !value);
    const hasCheckedChildren = permission.allChildren.some(
      ({ _id: id }) => !!role.permissions?.[id]?.length,
    );

    return {
      inputValue,
      indeterminate: !inputValue && hasCheckedChildren,
      ...conditionalState,
    };
  }

  const inputValue =
    conditionalState?.inputValue ?? role.permissions?.[permission._id]?.includes(action) ?? false;

  return {
    inputValue,
    ...conditionalState,
  };
};

/**
 * Checks if a permission ID belongs to the API permission group.
 *
 * @param {string} [permissionId='']
 * @returns {boolean}
 */
export const isApiPermission = (permissionId = '') => permissionId.startsWith(USER_PERMISSIONS_GROUPS.api);

/**
 * Filters treeview permissions by search term recursively.
 *
 * @param {Array} items
 * @param {string} [search='']
 * @param {number|null} [searchDepth]
 * @returns {Array}
 */
export const filterTreeviewPermissions = (items, search = '', searchDepth) => {
  if (!search || !search.trim()) {
    return items;
  }

  const searchLower = search.toLowerCase().trim();

  const matchesSearch = node => node?.title?.toLowerCase().includes(searchLower);

  const collectLeafIds = (nodes) => {
    const ids = new Set();

    nodes.forEach((node) => {
      if (node.children && node.children.length > 0) {
        collectLeafIds(node.children).forEach(id => ids.add(id));

        return;
      }

      ids.add(node._id);
    });

    return ids;
  };

  const filterNodesArray = (nodes, filterFn) => nodes.reduce((acc, node) => {
      const filteredNode = filterFn(node);

      if (filteredNode) {
        acc.push(filteredNode);
      }

      return acc;
    }, []);

  const filterNode = (node, currentDepth = 0) => {
    if (currentDepth === searchDepth || !node.children) {
      return matchesSearch(node) ? node : null;
    }

    const filteredChildren = filterNodesArray(node.children, childNode => filterNode(childNode, currentDepth + 1),
    );

    if (!filteredChildren.length) {
      return null;
    }

    const leafIds = collectLeafIds(filteredChildren);
    const filteredAllChildren = node.allChildren.filter(child => leafIds.has(child._id));

    return { ...node, children: filteredChildren, allChildren: filteredAllChildren };
  };

  return filterNodesArray(items, item => filterNode(item, 0));
};
