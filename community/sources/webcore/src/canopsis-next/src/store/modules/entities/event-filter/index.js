import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import { createEntityModule } from '@/store/plugins/entities';

export default createEntityModule({
  route: API_ROUTES.eventFilter.rules,
  entityType: ENTITIES_TYPES.eventFilter,
  dataPreparer: d => d.data,
  withFetchingParams: true,
  withMeta: true,
});
