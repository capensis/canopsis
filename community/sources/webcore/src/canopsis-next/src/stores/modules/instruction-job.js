import { API_ROUTES } from '@/config';

import { createJobManagementStore } from '@/stores/plugins/job-management';

export const useInstructionJobStore = createJobManagementStore({
  id: 'instructionJob',
  route: API_ROUTES.remediationJobs,
  bulkStopRoute: API_ROUTES.bulkRemediationJobs.stop,
});
