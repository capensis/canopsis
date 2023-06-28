import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import request from '@/services/request';

import { createEntityModule } from '@/store/plugins/entities';

export default createEntityModule({
  route: API_ROUTES.roles.list,
  entityType: ENTITIES_TYPES.role,
  dataPreparer: d => d.data,
  withMeta: true,
}, {
  actions: {
    fetchListWithoutStore(context, { params }) {
      return request.get(API_ROUTES.roles.list, { params });
    },

    fetchTemplatesListWithoutStore(context, { params }) {
      return request.get(API_ROUTES.roles.templates, { params });
    },
  },
});
