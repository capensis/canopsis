import { API_ROUTES } from '@/config';

import { ENTITIES_TYPES } from '@/constants';

import { createEntityModule } from '@/store/plugins/entities';

export default createEntityModule({
  route: API_ROUTES.scenario.scenarios,
  entityType: ENTITIES_TYPES.scenario,
  dataPreparer: d => d.data,
  withMeta: true,
});
