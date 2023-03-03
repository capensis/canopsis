import { API_ROUTES } from '@/config';

import { ENTITIES_TYPES } from '@/constants';

import { createEntityModule } from '@/store/plugins/entities';

import request from '@/services/request';

export default createEntityModule({
  route: API_ROUTES.dynamicInfo,
  entityType: ENTITIES_TYPES.dynamicInfo,
  dataPreparer: d => d.data,
  withMeta: true,
}, {
  actions: {
    fetchListWithoutStore(context, { params }) {
      return request.get(API_ROUTES.dynamicInfo, { params });
    },

    fetchInfosKeysWithoutStore(context, { params }) {
      return request.get(API_ROUTES.dynamicInfosDictionaryKeys, { params });
    },
  },
});
