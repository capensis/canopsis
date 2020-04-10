import { pick } from 'lodash';

import { GROUPS_NAVIGATION_TYPES } from '@/constants';

/**
 * Convert user to form
 *
 * @param {Object} [user]
 * @returns {Object}
 */
export function userToForm(user) {
  if (!user) {
    return {
      _id: '',
      firstname: '',
      lastname: '',
      mail: '',
      password: '',
      role: null,
      ui_language: '',
      enable: true,
      defaultview: '',
      groupsNavigationType: GROUPS_NAVIGATION_TYPES.sideBar,
    };
  }

  const result = pick(user, [
    '_id',
    'firstname',
    'lastname',
    'mail',
    'password',
    'role',
    'ui_language',
    'enable',
    'defaultview',
    'groupsNavigationType',
  ]);

  if (!result.groupsNavigationType) {
    result.groupsNavigationType = GROUPS_NAVIGATION_TYPES.sideBar;
  }

  return result;
}
