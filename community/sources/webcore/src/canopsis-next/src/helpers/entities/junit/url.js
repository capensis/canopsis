import { API_HOST, API_ROUTES } from '@/config';

/**
 * Get test suite file url
 *
 * @param {string} id
 * @return {string}
 */
export const getTestSuiteFileUrl = id => `${API_HOST}${API_ROUTES.junit.file}/${id}`;
