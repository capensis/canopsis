import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDStore } from '@/stores/plugins/crud';

export const useResolveRulesStore = createCRUDStore({
  id: 'resolveRules',
  route: API_ROUTES.resolveRules,
  withFetchingParams: true,
  withMeta: true,
  withWithoutStore: true,
  actions: {
    bulkEnable({ data } = {}) {
      return request.put(API_ROUTES.bulkResolveRulesEnable, data);
    },

    bulkDisable({ data } = {}) {
      return request.put(API_ROUTES.bulkResolveRulesDisable, data);
    },

    bulkRemove({ data } = {}) {
      return request.delete(API_ROUTES.bulkResolveRules, { data });
    },
  },
});
