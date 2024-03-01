import { API_HOST, API_ROUTES } from '@/config';

/**
 * Get availability download file url
 *
 * @param {string} id
 * @return {string}
 */
export const getAvailabilityDownloadFileUrl = id => `${API_HOST}${API_ROUTES.availability.list}/${id}/download`;
