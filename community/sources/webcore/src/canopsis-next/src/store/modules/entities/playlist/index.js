import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.playlist,
  withFetchingParams: true,
}, {
  actions: {
    fetchItemWithoutStore(context, { id }) {
      return request.get(`${API_ROUTES.playlist}/${id}`);
    },
  },
});
