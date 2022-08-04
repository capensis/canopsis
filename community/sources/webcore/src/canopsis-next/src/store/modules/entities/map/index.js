import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import request from '@/services/request';

import { createEntityModule } from '@/store/plugins/entities';

export default createEntityModule({
  route: API_ROUTES.maps,
  entityType: ENTITIES_TYPES.map,
  dataPreparer: d => d.data,
  withFetchingParams: true,
  withMeta: true,
}, {
  actions: {
    bulkRemove(context, { data }) {
      return request.delete(API_ROUTES.bulkMaps, { data });
    },

    fetchListWithoutStore() {
      return request.get(API_ROUTES.maps);
    },
  },
});
