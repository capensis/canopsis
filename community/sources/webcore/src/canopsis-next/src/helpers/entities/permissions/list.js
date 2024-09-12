import { get, isUndefined } from 'lodash';
import flatten from 'flat';

import {
  NOT_COMPLETED_USER_PERMISSIONS,
  USERS_PERMISSIONS,
  PERMISSIONS_TYPES_TO_ACTIONS,
  USER_PERMISSIONS_PREFIXES,
  CRUD_ACTIONS,
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

  return permissions.reduce((acc, permission) => {
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
      ...Object.values(USER_PERMISSIONS_PREFIXES.business).reduce((acc, key) => {
        acc[key] = [];

        return acc;
      }, {}),
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
};
