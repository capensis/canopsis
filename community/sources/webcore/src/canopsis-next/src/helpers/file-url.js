import { API_HOST, API_ROUTES } from '@/config';

/**
 * Get file url for test suite
 *
 * @param {string} id
 * @return {string}
 */
export const getTestSuiteFileUrl = id => `${API_HOST}${API_ROUTES.junit.file}/${id}`;

/**
 * Get tech metrics download file url
 *
 * @return {string}
 */
export const getTechMetricsDownloadFileUrl = () => `${API_HOST}${API_ROUTES.techMetrics}/download`;

/**
 * Get export metric download file url
 *
 * @param {string} id
 * @return {string}
 */
export const getExportMetricDownloadFileUrl = id => `${API_HOST}${API_ROUTES.metrics.exportMetric}/${id}/download`;

/**
 * Get alarm list export download file url
 *
 * @param {string} id
 * @return {string}
 */
export const getAlarmListExportDownloadFileUrl = id => `${API_HOST}${API_ROUTES.alarmListExport}/${id}/download`;

/**
 * Get context export download file url
 *
 * @param {string} id
 * @return {string}
 */
export const getContextExportDownloadFileUrl = id => `${API_HOST}${API_ROUTES.contextExport}/${id}/download`;
