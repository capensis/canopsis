import { API_HOST, API_ROUTES } from '@/config';

/**
 * Get alarm list export download file url
 *
 * @param {string} id
 * @returns {string}
 */
export const getAlarmListExportDownloadFileUrl = id =>
  `${API_HOST}${API_ROUTES.alarmListExport}/${id}/download`;

/**
 * Get alarm timeline export download file url
 *
 * @param {string} id
 * @returns {string}
 */
export const getAlarmTimelineExportDownloadFileUrl = id =>
  `${API_HOST}${API_ROUTES.alarmTimelineExport}/${id}/download`;
