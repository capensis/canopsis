import { API_HOST, API_ROUTES } from '@/config';

/**
 * Get export metric download file url
 *
 * @param {string} id
 * @return {string}
 */
export const getExportMetricDownloadFileUrl = id => `${API_HOST}${API_ROUTES.metrics.exportMetric}/${id}/download`;

/**
 * Get tech metrics download file url
 *
 * @return {string}
 */
export const getTechMetricsDownloadFileUrl = () => `${API_HOST}${API_ROUTES.techMetrics}/download`;
