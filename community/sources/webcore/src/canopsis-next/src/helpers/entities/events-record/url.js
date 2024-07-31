import { API_HOST, API_ROUTES } from '@/config';

/**
 * Get tech events record download file url
 *
 * @return {string}
 */
export const getEventsRecordFileUrl = () => `${API_HOST}${API_ROUTES.eventsRecord}/download`;
