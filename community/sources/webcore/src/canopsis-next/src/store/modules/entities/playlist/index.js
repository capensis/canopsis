import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import request from '@/services/request';

import { createEntityModule } from '@/store/plugins/entities';

export default createEntityModule({
  route: API_ROUTES.playlist,
  entityType: ENTITIES_TYPES.playlist,
  dataPreparer: d => d.data,
  withFetchingParams: true,
  withMeta: true,
}, {
  actions: {
    fetchItemWithoutStore(context, { id }) {
      return request.get(`${API_ROUTES.playlist}/${id}`);
    },
  },
});
