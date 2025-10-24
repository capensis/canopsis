/**
 * Check if user has special role
 *
 * @param {User} [user = {}]
 * @param {Role} [role = {}]
 * @return {boolean}
 */
export const isUserHasRole = (user = {}, role = {}) => user.roles?.some(({ _id: id }) => id === role._id);

/**
 * Get first role with non empty default view
 *
 * @param {User} user
 * @return {Role | undefined}
 */
export const getFirstRoleWithDefaultView = (user = {}) => user.roles?.find(role => !!role.defaultview);
