import { API_HOST, API_ROUTES } from '@/config';

/**
 * Get context export download file url
 *
 * @param {string} id
 * @return {string}
 */
export const getContextExportDownloadFileUrl = id => `${API_HOST}${API_ROUTES.contextExport}/${id}/download`;
