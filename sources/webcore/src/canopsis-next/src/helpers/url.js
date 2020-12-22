/**
 * Remove trailing slashes from url (http://example.com//login -> http://example.com/login)
 *
 * @param {string} [url = '']
 * @returns {string}
 */
export const removeTrailingSlashes = (url = '') => url.replace(/([^:]\/)\/+/g, '$1');
