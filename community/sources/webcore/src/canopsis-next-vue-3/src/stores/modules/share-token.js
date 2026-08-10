import { API_ROUTES } from '@/config';
import { createCRUDStore } from '@/stores/plugins/crud';

export const useShareTokenStore = createCRUDStore({
  id: 'shareToken',
  route: API_ROUTES.shareTokens,
  withFetchingParams: true,
});
