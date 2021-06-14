import { API_HOST, API_ROUTES } from '@/config';

/**
 * Remove trailing slashes from url (http://example.com//login -> http://example.com/login)
 *
 * @param {string} [url = '']
 * @returns {string}
 */
export const removeTrailingSlashes = (url = '') => url.replace(/([^:]\/)\/+/g, '$1');

/**
 * Get file url for test suite
 *
 * @param {string} id
 * @return {string}
 */
export const getTestSuiteFileUrl = id => `${API_HOST}${API_ROUTES.junit.file}/${id}`;
