import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.linkRule,
  bulkRoute: API_ROUTES.bulkLinkRule,
}, {
  actions: {
    fetchLinkCategoriesWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.linkCategories, { params });
    },

    bulkEnable(context, { data }) {
      return request.put(API_ROUTES.bulkLinkRuleEnable, data);
    },

    bulkDisable(context, { data }) {
      return request.put(API_ROUTES.bulkLinkRuleDisable, data);
    },
  },
});
