import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.roles.list,
  withWithoutStore: true,
}, {
  actions: {
    fetchTemplatesListWithoutStore(context, { params }) {
      return request.get(API_ROUTES.roles.templates, { params });
    },
  },
});
