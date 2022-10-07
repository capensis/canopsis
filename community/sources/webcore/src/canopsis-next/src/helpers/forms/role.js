import { CRUD_ACTIONS, TIME_UNITS } from '@/constants';

import { durationToForm } from '@/helpers/date/duration';

/**
 * @typedef {'create' | 'update' | 'reed' | 'delete'} Action
 */

/**
 * @typedef {Object} Permission
 * @property {Action[]} actions
 * @property {string} name
 * @property {string} type
 * @property {string} _id
 */

/**
 * @typedef {Object.<string, Action[]>}  PermissionsForm
 */

/**
 * @typedef {Object} AuthConfig
 * @property {Duration} [inactivity_interval]
 * @property {Duration} [expiration_interval]
 */

/**
 * @typedef {AuthConfig} AuthConfigForm
 * @property {Boolean} intervals_enabled
 */

/**
 * @typedef {Object} Role
 * @property {string} _id
 * @property {string} name
 * @property {string} description
 * @property {Permission[]} permissions
 * @property {DefaultView} defaultview
 * @property {AuthConfig} auth_config
 */

/**
 * @typedef {Role} RoleForm
 * @property {PermissionsForm} permissions
 * @property {string} defaultview
 */

/**
 * @typedef {Role} RoleRequest
 * @property {string} defaultview
 */

/**
 * Convert role permissions to form permissions object
 *
 * @param {Permission[]} [permissions = []]
 * @return {PermissionsForm}
 */
const rolePermissionsToForm = (permissions = []) => permissions.reduce((acc, { _id: id, actions }) => {
  if (actions.length === 0) {
    acc[id] = [CRUD_ACTIONS.can];
  } else {
    acc[id] = actions;
  }

  return acc;
}, {});

/**
 * Convert role auth config to form auth config object
 *
 * @param {AuthConfig} [authConfig]
 * @return {AuthConfigForm}
 */
const roleAuthConfigToForm = (authConfig = {}) => ({
  intervals_enabled: Boolean(authConfig.inactivity_interval || authConfig.expiration_interval),
  inactivity_interval: durationToForm(authConfig.inactivity_interval ?? { value: 1, unit: TIME_UNITS.hour }),
  expiration_interval: durationToForm(authConfig.expiration_interval ?? { value: 1, unit: TIME_UNITS.day }),
});

/**
 * Convert role to form object
 *
 * @param {Role} [role = {}]
 * @returns {RoleForm}
 */
export const roleToForm = (role = {}) => ({
  _id: role._id ?? '',
  name: role.name ?? '',
  description: role.description ?? '',
  defaultview: role.defaultview?._id,
  permissions: rolePermissionsToForm(role.permissions),
  auth_config: roleAuthConfigToForm(role.auth_config),
});

/**
 * Convert form role permissions to role permissions object
 *
 * @param {PermissionsForm} [permissionsForm={}]
 * @return {PermissionsForm}
 */
const permissionsFormToRolePermissions = (permissionsForm = {}) => Object.entries(permissionsForm)
  .reduce((acc, [id, actions]) => {
    const [firstAction] = actions;

    if (firstAction === CRUD_ACTIONS.can) {
      acc[id] = [];
    } else if (!firstAction) {
      acc[id] = undefined;
    } else {
      acc[id] = actions;
    }

    return acc;
  }, {});

/**
 * Convert form role auth config to role auth config object
 *
 * @param {AuthConfigForm} [form={}]
 * @return {AuthConfig}
 */
const authConfigFormToRolePermissions = (form = {}) => ({
  inactivity_interval: form.intervals_enabled ? form.inactivity_interval : null,
  expiration_interval: form.intervals_enabled ? form.expiration_interval : null,
});

/**
 * Convert role form to role object
 *
 * @param {RoleForm | {}} [form = {}]
 * @returns {RoleRequest}
 */
export const formToRole = (form = {}) => ({
  ...form,

  defaultview: form.defaultview,
  permissions: permissionsFormToRolePermissions(form.permissions),
  auth_config: authConfigFormToRolePermissions(form.auth_config),
});
