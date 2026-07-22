import { API_HOST, API_ROUTES } from '@/config'

export const getExportMetricDownloadFileUrl = (id) =>
  `${API_HOST}${API_ROUTES.metrics.exportMetric}/${id}/download`

export const getTechMetricsDownloadFileUrl = () => `${API_HOST}${API_ROUTES.techMetrics}/download`
