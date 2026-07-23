import { API_ROUTES } from '@/config';

import { createCRUDStore } from '@/stores/plugins/crud';

export const useFilterStore = createCRUDStore({
  id: 'filter',
  route: API_ROUTES.filters,
  withFetchingParams: true,
  withWithoutStore: true,
});
