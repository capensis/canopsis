import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import request from '@/services/request';

import { createEntityModule } from '@/store/plugins/entities';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default createEntityModule({
  types,
  route: API_ROUTES.planning.exceptions,
  entityType: ENTITIES_TYPES.pbehaviorExceptions,
  dataPreparer: d => d.data,
  withFetchingParams: true,
  withMeta: true,
}, {
  actions: {
    fetchListWithoutStore(context, { params }) {
      return request.get(API_ROUTES.planning.exceptions, { params });
    },
  },
});
