import request from '@/services/request';

import { API_ROUTES } from '@/config';

import { createCRUDStore } from '@/stores/plugins/crud';

export const usePlaylistStore = createCRUDStore({
  id: 'playlist',
  route: API_ROUTES.playlist,
  withFetchingParams: true,
  actions: {
    fetchItemWithoutStore({ id } = {}) {
      return request.get(`${API_ROUTES.playlist}/${id}`);
    },

    bulkEnable({ data } = {}) {
      return request.put(API_ROUTES.bulkPlaylistEnable, data);
    },

    bulkDisable({ data } = {}) {
      return request.put(API_ROUTES.bulkPlaylistDisable, data);
    },

    bulkRemove({ data } = {}) {
      return request.delete(API_ROUTES.bulkPlaylist, { data });
    },
  },
});
