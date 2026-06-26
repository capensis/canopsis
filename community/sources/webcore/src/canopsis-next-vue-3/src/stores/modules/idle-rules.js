import request from '@/services/request';

import { API_ROUTES } from '@/config';

import { createCRUDStore } from '@/stores/plugins/crud';

export const useIdleRulesStore = createCRUDStore({
  id: 'idleRules',
  route: API_ROUTES.idleRules,
  withFetchingParams: true,
  withMeta: true,
  withWithoutStore: true,
  actions: {
    bulkEnable({ data } = {}) {
      return request.put(API_ROUTES.bulkIdleRulesEnable, data);
    },

    bulkDisable({ data } = {}) {
      return request.put(API_ROUTES.bulkIdleRulesDisable, data);
    },

    bulkRemove({ data } = {}) {
      return request.delete(API_ROUTES.bulkIdleRules, { data });
    },
  },
});
