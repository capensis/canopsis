import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.dynamicInfo,
  withWithoutStore: true,
}, {
  actions: {
    fetchInfosKeysWithoutStore(context, { params }) {
      return request.get(API_ROUTES.dynamicInfosDictionaryKeys, { params });
    },

    bulkEnable(context, { data }) {
      return request.put(API_ROUTES.bulkDynamicInfoEnable, data);
    },

    bulkDisable(context, { data }) {
      return request.put(API_ROUTES.bulkDynamicInfoDisable, data);
    },

    bulkRemove(context, { data }) {
      return request.delete(API_ROUTES.bulkDynamicInfo, { data });
    },
  },
});
