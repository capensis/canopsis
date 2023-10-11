import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import { createCRUDModule } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.snmpRule,
  entityType: ENTITIES_TYPES.snmpRule,
  dataPreparer: d => d.data,
});
