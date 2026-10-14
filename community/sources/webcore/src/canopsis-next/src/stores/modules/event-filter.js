import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDStore } from '@/stores/plugins/crud';

export const useEventFilterStore = createCRUDStore({
  id: 'eventFilter',
  route: API_ROUTES.eventFilter.rules,
  withFetchingParams: true,
  withMeta: true,
  withWithoutStore: true,
  actions: {
    markNewEventFilterFailuresAsRead({ id } = {}) {
      return request.put(`${API_ROUTES.eventFilter.list}/${id}/failures`);
    },

    fetchEventFilterFailuresListWithoutStore({ id, params } = {}) {
      return request.get(`${API_ROUTES.eventFilter.list}/${id}/failures`, { params });
    },

    bulkEnable({ data } = {}) {
      return request.put(API_ROUTES.eventFilter.bulkEnable, data);
    },

    bulkDisable({ data } = {}) {
      return request.put(API_ROUTES.eventFilter.bulkDisable, data);
    },

    bulkRemove({ data } = {}) {
      return request.delete(API_ROUTES.eventFilter.bulk, { data });
    },
  },
});
