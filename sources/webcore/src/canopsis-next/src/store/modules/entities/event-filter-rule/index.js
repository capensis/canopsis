import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import { createEntityModule } from '@/store/plugins/entities';

export default createEntityModule({
  route: API_ROUTES.eventFilterRules,
  entityType: ENTITIES_TYPES.eventFilterRule,
  dataPreparer: d => d.data,
  withFetchingParams: true,
  withMeta: true,
});
