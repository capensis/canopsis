import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.maps,
  withFetchingParams: true,
  withWithoutStore: true,
}, {
  actions: {
    fetchItemWithoutStore(context, { id, params }) {
      return request.get(`${API_ROUTES.maps}/${id}`, { params });
    },

    bulkRemove(context, { data }) {
      return request.delete(API_ROUTES.bulkMaps, { data });
    },

    fetchItemStateWithoutStore(context, { id, params }) {
      return request.get(`${API_ROUTES.mapState}/${id}`, { params });
    },
  },
});
