import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import { createEntityModule } from '@/store/plugins/entities';

export default createEntityModule({
  route: API_ROUTES.shareTokens,
  entityType: ENTITIES_TYPES.shareToken,
  dataPreparer: d => d.data,
  withFetchingParams: true,
  withMeta: true,
});
