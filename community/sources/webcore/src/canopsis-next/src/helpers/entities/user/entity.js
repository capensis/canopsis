import { ROLE_TYPES } from '@/constants';

/**
 * Check if user has special role
 *
 * @param {User} [user = {}]
 * @param {Role} [role = {}]
 * @return {boolean}
 */
export const isUserHasRole = (user = {}, role = {}) => user.roles?.some?.(({ _id: id }) => id === role._id);

/**
 * Check if user has special role
 *
 * @param {User} [user = {}]
 * @return {boolean}
 */
export const isUserHasOnlyApiRole = (user = {}) => user.roles?.every?.(({ type }) => type === ROLE_TYPES.api);

/**
 * Get first role with non empty default view
 *
 * @param {User} user
 * @return {Role | undefined}
 */
export const getFirstRoleWithDefaultView = (user = {}) => user.roles?.find(role => !!role.defaultview);
