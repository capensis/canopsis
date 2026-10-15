import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDStore } from '@/stores/plugins/crud';

export const useMapStore = createCRUDStore({
  id: 'map',
  route: API_ROUTES.maps,
  bulkRoute: API_ROUTES.bulkMaps,
  withFetchingParams: true,
  withWithoutStore: true,
  actions: {
    fetchItemWithoutStore({ id, params } = {}) {
      return request.get(`${API_ROUTES.maps}/${id}`, { params });
    },

    fetchItemStateWithoutStore({ id, params } = {}) {
      return request.get(`${API_ROUTES.mapState}/${id}`, { params });
    },
  },
});
