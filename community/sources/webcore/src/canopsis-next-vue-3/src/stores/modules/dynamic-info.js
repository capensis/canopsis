import { API_ROUTES } from '@/config';
import request from '@/services/request';
import { createCRUDStore } from '@/stores/plugins/crud';

export const useDynamicInfoStore = createCRUDStore({
  id: 'dynamicInfo',
  route: API_ROUTES.dynamicInfo,
  withFetchingParams: true,
  withMeta: true,
  withWithoutStore: true,
  actions: {
    fetchInfosKeysWithoutStore({ params } = {}) {
      return request.get(API_ROUTES.dynamicInfosDictionaryKeys, { params });
    },

    bulkEnable({ data } = {}) {
      return request.put(API_ROUTES.bulkDynamicInfoEnable, data);
    },

    bulkDisable({ data } = {}) {
      return request.put(API_ROUTES.bulkDynamicInfoDisable, data);
    },

    bulkRemove({ data } = {}) {
      return request.delete(API_ROUTES.bulkDynamicInfo, { data });
    },
  },
});
