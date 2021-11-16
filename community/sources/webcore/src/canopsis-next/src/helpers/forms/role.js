import { CRUD_ACTIONS } from '@/constants';

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
 * @typedef {Object} Role
 * @property {string} _id
 * @property {string} name
 * @property {string} description
 * @property {Permission[]} permissions
 * @property {DefaultView} defaultview
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
 * Convert role to form object
 *
 * @param {Role} [role = {}]
 * @returns {RoleForm}
 */
export const roleToForm = (role = {}) => ({
  _id: role._id || '',
  name: role.name || '',
  description: role.description || '',
  defaultview: role.defaultview && role.defaultview._id,
  permissions: rolePermissionsToForm(role.permissions),
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
 * Convert role form to role object
 *
 * @param {RoleForm | {}} [form = {}]
 * @returns {RoleRequest}
 */
export const formToRole = (form = {}) => ({
  ...form,

  defaultview: form.defaultview,
  permissions: permissionsFormToRolePermissions(form.permissions),
});
