import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import { createEntityModule } from '@/store/plugins/entities';

export default createEntityModule({
  route: API_ROUTES.scenarios,
  entityType: ENTITIES_TYPES.scenario,
  withFetchingParams: true,
  dataPreparer: d => d.data,
});
