import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import { createEntityModule } from '@/store/plugins/entities';
import request from '@/services/request';

export default createEntityModule({
  route: API_ROUTES.pbehavior.types,
  entityType: ENTITIES_TYPES.pbehaviorTypes,
  dataPreparer: d => d.data,
  withFetchingParams: true,
  withMeta: true,
}, {
  actions: {
    fetchListWithoutStore(context, { params }) {
      return request.get(API_ROUTES.pbehavior.types, { params });
    },
  },
});
