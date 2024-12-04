/**
 * Converts an array of permissions into a single string by joining them with a delimiter
 *
 * @param {Array<string>} [permissions=[]] - Array of permission strings to be joined
 * @param {string} [delimiter=''] - The separator to use between permissions
 * @returns {string} - The joined permissions string
 */
export const groupedPermissionToPermission = (permissions = [], delimiter = '') => permissions.join(delimiter);
