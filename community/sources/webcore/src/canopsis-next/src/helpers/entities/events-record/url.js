import { API_HOST, API_ROUTES } from '@/config';

/**
 * Get tech events record download file url
 *
 * @return {string}
 */
export const getEventsRecordFileUrl = (id = '') => `${API_HOST}${API_ROUTES.eventsRecordExport}/${id}/download`;
