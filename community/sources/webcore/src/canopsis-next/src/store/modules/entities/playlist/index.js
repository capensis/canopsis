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

    bulkEnable(context, { data }) {
      return request.put(API_ROUTES.bulkPlaylistEnable, data);
    },

    bulkDisable(context, { data }) {
      return request.put(API_ROUTES.bulkPlaylistDisable, data);
    },

    bulkRemove(context, { data }) {
      return request.delete(API_ROUTES.bulkPlaylist, { data });
    },
  },
});
