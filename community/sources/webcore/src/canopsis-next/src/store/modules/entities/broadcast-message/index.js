import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.broadcastMessage.list,
  withFetchingParams: true,
}, {
  actions: {
    async fetchListWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.broadcastMessage.list, { params });
    },

    async fetchActiveListWithoutStore() {
      return request.get(API_ROUTES.broadcastMessage.activeList);
    },
  },
});
