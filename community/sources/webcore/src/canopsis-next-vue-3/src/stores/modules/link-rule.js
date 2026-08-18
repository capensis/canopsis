import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDStore } from '@/stores/plugins/crud';

export const useLinkRuleStore = createCRUDStore({
  id: 'linkRule',
  route: API_ROUTES.linkRule,
  bulkRoute: API_ROUTES.bulkLinkRule,
  withFetchingParams: true,
  withWithoutStore: true,
  actions: {
    fetchLinkCategoriesWithoutStore({ params } = {}) {
      return request.get(API_ROUTES.linkCategories, { params });
    },

    bulkEnable({ data } = {}) {
      return request.put(API_ROUTES.bulkLinkRuleEnable, data);
    },

    bulkDisable({ data } = {}) {
      return request.put(API_ROUTES.bulkLinkRuleDisable, data);
    },
  },
});
