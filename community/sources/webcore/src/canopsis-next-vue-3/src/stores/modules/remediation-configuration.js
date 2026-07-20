import { API_ROUTES } from '@/config'

import { createCRUDStore } from '@/stores/plugins/crud'

export const useRemediationConfigurationStore = createCRUDStore({
  id: 'remediationConfiguration',
  route: API_ROUTES.remediation.configurations,
  withFetchingParams: true,
  withWithoutStore: true,
})
