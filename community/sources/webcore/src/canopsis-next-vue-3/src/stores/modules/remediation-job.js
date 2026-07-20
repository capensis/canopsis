import { API_ROUTES } from '@/config'

import { createCRUDStore } from '@/stores/plugins/crud'

export const useRemediationJobStore = createCRUDStore({
  id: 'remediationJob',
  route: API_ROUTES.remediation.jobs,
  withFetchingParams: true,
  withWithoutStore: true,
})
