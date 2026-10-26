import { API_ROUTES } from '@/config';

import { createJobManagementStore } from '@/stores/plugins/job-management';

export const useWebhookJobStore = createJobManagementStore({
  id: 'webhookJob',
  route: API_ROUTES.webhookJobs,
  bulkStopRoute: API_ROUTES.bulkWebhookJobs.stop,
});
