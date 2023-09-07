/**
 * Convert object to FormData object for requests
 *
 * @param {Object} obj
 * @returns {FormData}
 */
export const convertObjectToFormData = (obj = {}) => Object.entries(obj).reduce((acc, [key, value]) => {
  acc.append(key, value);
  return acc;
}, new FormData());
