import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.eventFilter.rules,
  withFetchingParams: true,
  withMeta: true,
}, {
  actions: {
    markNewEventFilterFailuresAsRead(context, { id }) {
      return request.put(`${API_ROUTES.eventFilter.list}/${id}/failures`);
    },

    fetchEventFilterFailuresListWithoutStore(context, { id, params } = {}) {
      return request.get(`${API_ROUTES.eventFilter.list}/${id}/failures`, { params });
    },
  },
});
