import { isString } from 'lodash';

import { ROLE_TYPES } from '@/constants';

/**
 * Check if user has special role
 *
 * @param {User} [user = {}]
 * @param {Role} [role = {}]
 * @return {boolean}
 */
export const isUserHasRole = (user = {}, role = {}) =>
  user.roles?.some?.(({ _id: id }) => id === role._id);

/**
 * Check if user has special role
 *
 * @param {User} [user = {}]
 * @return {boolean}
 */
export const isUserHasOnlyApiRole = (user = {}) =>
  user.roles?.every?.(({ type }) => type === ROLE_TYPES.api);

/**
 * Get first role with non empty default view
 *
 * @param {User} user
 * @return {Role | undefined}
 */
export const getFirstRoleWithDefaultView = (user = {}) =>
  user.roles?.find(role => !!role.defaultview || !!role.default_view);

/**
 * Get default view id from default view entity or id.
 *
 * @param {DefaultView|string} defaultView
 * @return {string|undefined}
 */
export const getDefaultViewId = defaultView =>
  isString(defaultView) ? defaultView : defaultView?._id;

/**
 * Gets user or role default view from legacy and API naming variants.
 *
 * @param {User|Role} entity
 * @returns {DefaultView|string|undefined}
 */
export const getDefaultView = (entity = {}) => entity.defaultview ?? entity.default_view;
