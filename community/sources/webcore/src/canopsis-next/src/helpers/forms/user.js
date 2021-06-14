import { GROUPS_NAVIGATION_TYPES } from '@/constants';

/**
 * @typedef {Object} User
 * @property {string} crecord_name
 * @property {string} firstname
 * @property {string} password
 * @property {string|null} role
 * @property {string} mail
 * @property {boolean} enable
 * @property {string} groupsNavigationType
 * @property {string} lastname
 * @property {string} ui_language
 */

/**
 * Convert user to form object
 *
 * @param {User} [user = {}]
 * @returns {User}
 */
export const userToForm = (user = {}) => ({
  _id: user._id || '',
  crecord_name: user.crecord_name || '',
  firstname: user.firstname || '',
  lastname: user.lastname || '',
  mail: user.mail || '',
  password: user.password || '',
  role: user.role || null,
  ui_language: user.ui_language || '',
  enable: !!user.enable,
  defaultview: user.defaultview || '',
  groupsNavigationType: user.groupsNavigationType || GROUPS_NAVIGATION_TYPES.sideBar,
});
