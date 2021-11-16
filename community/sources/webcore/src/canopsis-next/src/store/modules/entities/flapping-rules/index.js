import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import { createEntityModule } from '@/store/plugins/entities';

export default createEntityModule({
  route: API_ROUTES.flappingRules,
  entityType: ENTITIES_TYPES.flappingRules,
  dataPreparer: d => d.data,
  withMeta: true,
});
