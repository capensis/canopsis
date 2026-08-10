import { API_ROUTES } from '@/config';
import request from '@/services/request';
import { createCRUDStore } from '@/stores/plugins/crud';

export const useBroadcastMessageStore = createCRUDStore({
  id: 'broadcastMessage',
  route: API_ROUTES.broadcastMessage.list,
  withFetchingParams: true,
  withMeta: true,
  withWithoutStore: true,
  actions: {
    fetchActiveListWithoutStore() {
      return request.get(API_ROUTES.broadcastMessage.activeList);
    },

    markAsRead({ id } = {}) {
      return request.put(`${API_ROUTES.broadcastMessage.list}/${id}/read`);
    },
  },
});
