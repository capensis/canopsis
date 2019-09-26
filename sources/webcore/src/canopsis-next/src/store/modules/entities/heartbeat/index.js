import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import { createEntityModule } from '@/store/plugins/entities';
import request from '@/services/request';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default createEntityModule({
  types,
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
