import { API_HOST, API_ROUTES } from '@/config';

/**
 * Get tech external data table records download file url
 *
 * @return {string}
 */
export const getExternalDataTableRecordsFileUrl = (id = '') => (
  `${API_HOST}${API_ROUTES.externalDataExport}/${id}/download`
);
