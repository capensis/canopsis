import { API_HOST, API_ROUTES } from '@/config';

/**
 * Get alarm list export download file url
 *
 * @param {string} id
 * @return {string}
 */
export const getAlarmListExportDownloadFileUrl = id => `${API_HOST}${API_ROUTES.alarmListExport}/${id}/download`;
