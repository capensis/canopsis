import { API_ROUTES } from '@/config';
import { createCRUDStore } from '@/stores/plugins/crud';

export const useThemeStore = createCRUDStore({
  id: 'theme',
  route: API_ROUTES.themes.list,
  bulkRoute: API_ROUTES.themes.bulkList,
  withFetchingParams: true,
  withWithoutStore: true,
});
