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
  },
});
