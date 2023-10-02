import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createWidgetModule } from '@/store/plugins/entities';

export default createWidgetModule({ route: API_ROUTES.metrics.group }, {
  actions: {
    createExport(context, { data }) {
      return request.post(API_ROUTES.metrics.exportGroup, data);
    },
  },
});
