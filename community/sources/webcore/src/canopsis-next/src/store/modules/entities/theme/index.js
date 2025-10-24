import { API_ROUTES } from '@/config';

import { createCRUDModule } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.themes.list,
  bulkRoute: API_ROUTES.themes.bulkList,
  withFetchingParams: true,
  withWithoutStore: true,
});
