import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import { createEntityModule } from '@/store/plugins/entities';

export default createEntityModule({
  route: API_ROUTES.tags,
  entityType: ENTITIES_TYPES.tag,
  dataPreparer: d => d.data,
  withFetchingParams: true,
  withMeta: true,
});
