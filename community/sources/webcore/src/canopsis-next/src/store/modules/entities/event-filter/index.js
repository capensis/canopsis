import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import request from '@/services/request';

import { createEntityModule } from '@/store/plugins/entities';

export default createEntityModule({
  route: API_ROUTES.eventFilter.rules,
  entityType: ENTITIES_TYPES.eventFilter,
  dataPreparer: d => d.data,
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
