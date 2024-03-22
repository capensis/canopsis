import { API_HOST, API_ROUTES } from '@/config';

/**
 * Get availability download file url
 *
 * @param {string} id
 * @return {string}
 */
export const getAvailabilityDownloadFileUrl = id => `${API_HOST}${API_ROUTES.metrics.exportMetric}/${id}/download`;

/**
 * Get availability history download file url
 *
 * @param {string} id
 * @return {string}
 */
export const getAvailabilityHistoryDownloadFileUrl = id => `${API_HOST}${API_ROUTES.metrics.exportMetric}/${id}/download`;
