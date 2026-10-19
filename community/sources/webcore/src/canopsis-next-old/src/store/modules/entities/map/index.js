import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.maps,
  bulkRoute: API_ROUTES.bulkMaps,
  withFetchingParams: true,
  withWithoutStore: true,
}, {
  actions: {
    fetchItemWithoutStore(context, { id, params }) {
      return request.get(`${API_ROUTES.maps}/${id}`, { params });
    },

    fetchItemStateWithoutStore(context, { id, params }) {
      return request.get(`${API_ROUTES.mapState}/${id}`, { params });
    },
  },
});
