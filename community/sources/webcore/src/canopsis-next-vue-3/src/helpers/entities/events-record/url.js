import { API_HOST, API_ROUTES } from '@/config'

export const getEventsRecordFileUrl = (id = '') =>
  `${API_HOST}${API_ROUTES.eventsRecord.export}/${id}/download`
