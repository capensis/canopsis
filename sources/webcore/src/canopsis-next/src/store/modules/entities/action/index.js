import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import { createEntityModule } from '@/store/plugins/entities';

export default createEntityModule({
  route: API_ROUTES.actions,
  entityType: ENTITIES_TYPES.action,
  withFetchingParams: true,
});
