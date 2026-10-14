import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.pbehavior.reasons,
  withFetchingParams: true,
  withWithoutStore: true,
}, {
  actions: {
    bulkHide(context, { data }) {
      return request.put(API_ROUTES.pbehavior.bulkReasonsHide, data);
    },

    bulkUnhide(context, { data }) {
      return request.put(API_ROUTES.pbehavior.bulkReasonsUnhide, data);
    },

    bulkRemove(context, { data }) {
      return request.delete(API_ROUTES.pbehavior.bulkReasons, { data });
    },
  },
});
