import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import { createEntityModule } from '@/store/plugins/entities';
import request from '@/services/request';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default createEntityModule({
  types,
  route: API_ROUTES.view.groups,
  entityType: ENTITIES_TYPES.group,
  dataPreparer: d => d.data,
  withMeta: true,
}, {
  actions: {
    fetchListWithoutStore(context, { params }) {
      return request.get(API_ROUTES.view.groups, { params });
    },
  },
});
