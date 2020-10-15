import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import { createEntityModule } from '@/store/plugins/entities';
import request from '@/services/request';

export default createEntityModule({
  route: API_ROUTES.heartbeat,
  entityType: ENTITIES_TYPES.heartbeat,
  withFetchingParams: true,
}, {
  actions: {
    remove(context, { id } = {}) {
      return request.delete(`${API_ROUTES.heartbeat}${id}`);
    },
  },
});
