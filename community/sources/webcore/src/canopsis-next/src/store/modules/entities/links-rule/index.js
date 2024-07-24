import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.linkRule,
}, {
  actions: {
    bulkRemove(context, { data }) {
      return request.delete(API_ROUTES.bulkLinkRule, { data });
    },

    fetchLinkCategoriesWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.linkCategories, { params });
    },
  },
});
