import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDStore } from '@/stores/plugins/crud';

export const useFlappingRulesStore = createCRUDStore({
  id: 'flappingRules',
  route: API_ROUTES.flappingRules,
  withFetchingParams: true,
  withMeta: true,
  withWithoutStore: true,
  actions: {
    bulkEnable({ data } = {}) {
      return request.put(API_ROUTES.bulkFlappingRulesEnable, data);
    },

    bulkDisable({ data } = {}) {
      return request.put(API_ROUTES.bulkFlappingRulesDisable, data);
    },

    bulkRemove({ data } = {}) {
      return request.delete(API_ROUTES.bulkFlappingRules, { data });
    },
  },
});
