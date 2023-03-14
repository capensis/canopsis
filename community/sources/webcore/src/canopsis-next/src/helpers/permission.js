import { get, isUndefined, omit, sortBy } from 'lodash';
import flatten from 'flat';

import {
  NOT_COMPLETED_USER_PERMISSIONS,
  USERS_PERMISSIONS,
  PERMISSIONS_TYPES_TO_ACTIONS,
  USER_PERMISSIONS_PREFIXES, CRUD_ACTIONS,
} from '@/constants';

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
 * Get actions for permission by type
 *
 * @param {Object} permission
 * @returns {*}
 */
export const getPermissionActions = permission => (PERMISSIONS_TYPES_TO_ACTIONS[permission.type]
  ? PERMISSIONS_TYPES_TO_ACTIONS[permission.type]
  : [CRUD_ACTIONS.can]);

/**
 * Get right checkbox value for the form
 *
 * @param {string} permissionId
 * @param {PermissionForm} permissions
 * @param {PermissionForm} changedPermissions
 * @param {string} [action=CRUD_ACTIONS.can]
 * @returns {boolean}
 */
export const getCheckboxValue = (
  permissionId,
  permissions,
  changedPermissions,
  action = CRUD_ACTIONS.can,
) => {
  const actions = permissions[permissionId] || [];

  const changedActions = get(changedPermissions, [permissionId]);
  const currentActions = isUndefined(changedActions) ? actions : changedActions;

  return currentActions.includes(action);
};

/**
 * Get prepared grouped permissions for the permissions page
 *
 * @param {Array<Object>} permissions
 * @returns {*}
 */
export const getGroupedPermissions = (permissions) => {
  const allBusinessPermissionsIds = Object.values(flatten(USERS_PERMISSIONS.business));
  const generalApiPermissions = Object.values(USERS_PERMISSIONS.api.general);
  const rulesApiPermissions = Object.values(USERS_PERMISSIONS.api.rules);
  const remediationApiPermissions = Object.values(USERS_PERMISSIONS.api.remediation);
  const pbehaviorApiPermissions = Object.values(USERS_PERMISSIONS.api.pbehavior);

  const {
    exploitation: exploitationTechnicalPermissions,
    notification: notificationTechnicalPermissions,
    profile: profileTechnicalPermissions,
    ...adminTechnicalPermissions
  } = USERS_PERMISSIONS.technical;
  const adminTechnicalPermissionsValues = Object.values(adminTechnicalPermissions);
  const exploitationTechnicalPermissionsValues = Object.values(exploitationTechnicalPermissions);
  const notificationTechnicalPermissionsValues = Object.values(notificationTechnicalPermissions);
  const profileTechnicalPermissionsValues = Object.values(profileTechnicalPermissions);

  const groupedPermissions = permissions.reduce((acc, permission) => {
    const permissionId = String(permission._id);

    if (permission.view && permission.view_group) {
      acc.view.push(permission);
    } else if (permission.playlist) {
      acc.playlist.push(permission);
    } else if (adminTechnicalPermissionsValues.includes(permissionId)) {
      acc.technical.admin.push(permission);
    } else if (exploitationTechnicalPermissionsValues.includes(permissionId)) {
      acc.technical.exploitation.push(permission);
    } else if (notificationTechnicalPermissionsValues.includes(permissionId)) {
      acc.technical.notification.push(permission);
    } else if (profileTechnicalPermissionsValues.includes(permissionId)) {
      acc.technical.profile.push(permission);
    } else if (
      allBusinessPermissionsIds.includes(permissionId)
      || NOT_COMPLETED_USER_PERMISSIONS.some(id => permissionId.startsWith(id))
    ) {
      const [parentKey] = permission._id.split('_');

      if (acc.business[parentKey]) {
        acc.business[parentKey].push(permission);
      }
    } else if (generalApiPermissions.includes(permissionId)) {
      acc.api.general.push(permission);
    } else if (rulesApiPermissions.includes(permissionId)) {
      acc.api.rules.push(permission);
    } else if (remediationApiPermissions.includes(permissionId)) {
      acc.api.remediation.push(permission);
    } else if (pbehaviorApiPermissions.includes(permissionId)) {
      acc.api.pbehavior.push(permission);
    }

    return acc;
  }, {
    business: {
      [USER_PERMISSIONS_PREFIXES.business.common]: [],
      [USER_PERMISSIONS_PREFIXES.business.alarmsList]: [],
      [USER_PERMISSIONS_PREFIXES.business.context]: [],
      [USER_PERMISSIONS_PREFIXES.business.serviceWeather]: [],
      [USER_PERMISSIONS_PREFIXES.business.counter]: [],
      [USER_PERMISSIONS_PREFIXES.business.testingWeather]: [],
      [USER_PERMISSIONS_PREFIXES.business.map]: [],
    },
    view: [],
    playlist: [],
    technical: {
      admin: [],
      exploitation: [],
      notification: [],
      profile: [],
    },
    api: {
      general: [],
      rules: [],
      remediation: [],
      pbehavior: [],
    },
  });

  /**
   * We are using order which one we've defined on the reduce accumulator initial value.
   * For not `number`/`number string` object keys ordering is staying like we define
   */
  groupedPermissions.business = Object.entries(groupedPermissions.business)
    .map(([key, groupPermissions]) => ({
      key: `permission.business.${key}`,
      permissions: sortBy(groupPermissions, ['description']),
    }));

  /**
   * Ordering behavior have the same behavior as for `business`
   */
  groupedPermissions.technical = Object.entries(groupedPermissions.technical)
    .map(([key, groupPermissions]) => ({
      key: `permission.technical.${key}`,
      permissions: sortBy(groupPermissions, ['description']),
    }));

  /**
   * Ordering behavior have the same behavior as for `api`
   */
  groupedPermissions.api = Object.entries(groupedPermissions.api)
    .map(([key, groupPermissions]) => ({
      key: `permission.api.${key}`,
      permissions: sortBy(groupPermissions, ['description']),
    }));

  const viewsPermissionsByGroupTitle = groupedPermissions.view.reduce((acc, permission) => {
    const { view_group: viewGroup, ...rest } = permission;

    if (!acc[viewGroup._id]) {
      acc[viewGroup._id] = {
        viewGroup,
        permissions: [],
      };
    }

    acc[viewGroup._id].permissions.push(rest);

    return acc;
  }, {});

  groupedPermissions.view = sortBy(Object.values(viewsPermissionsByGroupTitle), ['viewGroup.position'])
    .map(({ viewGroup, permissions: viewGroupPermissions }) => ({
      name: viewGroup.title,
      permissions: sortBy(viewGroupPermissions, ['view.position']),
    }));

  groupedPermissions.view.push({
    key: 'common.playlist',
    permissions: sortBy(groupedPermissions.playlist, ['playlist.name']),
  });

  return omit(groupedPermissions, ['playlist']);
};
