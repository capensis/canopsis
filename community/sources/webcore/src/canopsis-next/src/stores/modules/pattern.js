import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDStore } from '@/stores/plugins/crud';

export const usePatternStore = createCRUDStore({
  id: 'pattern',
  route: API_ROUTES.pattern.list,
  bulkRoute: API_ROUTES.pattern.bulkList,
  withFetchingParams: true,
  withMeta: true,
  withWithoutStore: true,
  actions: {
    checkPatternsEntitiesCount({ data } = {}) {
      return request.post(API_ROUTES.pattern.entitiesCount, data);
    },

    checkPatternsAlarmsCount({ data } = {}) {
      return request.post(API_ROUTES.pattern.alarmsCount, data);
    },
  },
});
