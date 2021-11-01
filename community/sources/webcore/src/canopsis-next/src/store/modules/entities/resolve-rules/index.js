import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import { createEntityModule } from '@/store/plugins/entities';

export default createEntityModule({
  route: API_ROUTES.resolveRules,
  entityType: ENTITIES_TYPES.resolveRules,
  dataPreparer: d => d.data,
  withMeta: true,
});
