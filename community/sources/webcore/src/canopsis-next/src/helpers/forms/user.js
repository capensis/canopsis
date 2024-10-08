import { THEMES_NAMES } from '@/config';

import { GROUPS_NAVIGATION_TYPES } from '@/constants';

/**
 * @typedef {Object} DefaultView
 * @property {string} _id
 * @property {string} name
 */

/**
 * @typedef {Object} User
 * @property {string} name
 * @property {string} firstname
 * @property {string} lastname
 * @property {string} password
 * @property {Role} role
 * @property {string} email
 * @property {boolean} enable
 * @property {Object} ui_tours
 * @property {DefaultView} defaultview
 * @property {string} groups_navigation_type
 * @property {string} ui_language
 * @property {string} ui_theme
 */

/**
 * @typedef {User} UserForm
 * @property {string} defaultview
 * @property {string} role
 */

/**
 * Convert user to form object
 *
 * @param {User} [user = {}]
 * @returns {UserForm}
 */
export const userToForm = (user = {}) => ({
  _id: user._id ?? '',
  name: user.name ?? '',
  firstname: user.firstname ?? '',
  lastname: user.lastname ?? '',
  email: user.email ?? '',
  password: user.password ?? '',
  role: user.role,
  enable: user.enable ?? true,
  defaultview: user.defaultview ? user.defaultview._id : '',
  ui_language: user.ui_language ?? '',
  ui_tours: user.ui_tours ?? {},
  ui_groups_navigation_type: user.ui_groups_navigation_type ?? GROUPS_NAVIGATION_TYPES.sideBar,
  ui_theme: user.ui_theme || THEMES_NAMES.canopsis,
});

/**
 * Convert user form to user object
 *
 * @param {UserForm} [form = {}]
 * @returns {User}
 */
export const formToUser = (form = {}) => ({
  ...form,

  role: form.role ? form.role._id : '',
});
