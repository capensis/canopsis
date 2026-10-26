import { API_ROUTES } from '@/config';

import { createJobManagementStore } from '@/stores/plugins/job-management';

export const useWebhookTokenJobStore = createJobManagementStore({
  id: 'webhookTokenJob',
  route: API_ROUTES.webhookTokenJobs,
  bulkStopRoute: API_ROUTES.bulkWebhookTokenJobs.stop,
});
