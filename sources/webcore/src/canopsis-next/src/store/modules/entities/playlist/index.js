import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import request from '@/services/request';

import { createEntityModule } from '@/store/plugins/entities';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default createEntityModule({
  types,
  route: API_ROUTES.playlist,
  entityType: ENTITIES_TYPES.playlist,
  withFetchingParams: true,
}, {
  actions: {
    fetchItemWithoutStore(context, { id }) {
      return request.get(`${API_ROUTES.playlist}/${id}`);
    },
  },
});
