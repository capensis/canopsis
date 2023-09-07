import { API_ROUTES } from '@/config';

import { createCRUDModule } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.remediation.configurations,
  withFetchingParams: true,
  withWithoutStore: true,
});
