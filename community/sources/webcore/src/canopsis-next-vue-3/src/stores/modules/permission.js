import { API_ROUTES } from '@/config';
import { createCRUDStore } from '@/stores/plugins/crud';

export const usePermissionStore = createCRUDStore({
  id: 'permission',
  route: API_ROUTES.permissions,
  withWithoutStore: true,
});
