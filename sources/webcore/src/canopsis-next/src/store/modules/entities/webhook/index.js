import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import { createEntityModule } from '@/store/plugins/entities';

export default createEntityModule({
  route: API_ROUTES.webhook,
  entityType: ENTITIES_TYPES.webhook,
  withFetchingParams: true,
  dataPreparer: d => d.data,
});
