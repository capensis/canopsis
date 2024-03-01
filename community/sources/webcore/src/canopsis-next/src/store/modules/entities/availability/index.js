import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createWidgetModule } from '@/store/plugins/entities';

export default createWidgetModule({ route: API_ROUTES.availability.list }, {
  actions: {
    createAvailabilityExport(context, { data }) {
      return request.post(API_ROUTES.availability.exportList, data);
    },

    fetchAvailabilityExport(context, { params, id }) {
      return request.get(`${API_ROUTES.availability.exportList}/${id}`, { params });
    },
  },
});
