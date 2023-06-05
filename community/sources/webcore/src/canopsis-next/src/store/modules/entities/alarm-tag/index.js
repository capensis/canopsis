import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import { createEntityModule } from '@/store/plugins/entities';

export default createEntityModule({
  route: API_ROUTES.alarmTags,
  entityType: ENTITIES_TYPES.alarmTag,
  dataPreparer: d => d.data,
});
