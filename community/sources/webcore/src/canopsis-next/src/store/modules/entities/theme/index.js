import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.themes.list,
  withFetchingParams: true,
  withWithoutStore: true,
}, {
  bulkRemove(context, { data }) {
    return request.delete(API_ROUTES.themes.bulkList, { data });
  },
});
