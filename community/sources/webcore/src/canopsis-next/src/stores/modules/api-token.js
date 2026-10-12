import { API_ROUTES } from '@/config';

import { createCRUDStore } from '@/stores/plugins/crud';

export const useApiTokenStore = createCRUDStore({
  id: 'apiToken',
  route: API_ROUTES.apiTokens,
  withFetchingParams: true,
});
