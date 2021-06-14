import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import { createEntityModule } from '@/store/plugins/entities';

export default createEntityModule({
  route: API_ROUTES.stateSetting,
  entityType: ENTITIES_TYPES.stateSetting,
  dataPreparer: d => d.data,
});
