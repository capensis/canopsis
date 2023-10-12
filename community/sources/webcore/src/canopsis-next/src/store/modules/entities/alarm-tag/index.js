import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.alarmTag.list,
  withWithoutStore: true,
}, {
  actions: {
    bulkRemove(context, { data }) {
      return request.delete(API_ROUTES.alarmTag.bulkList, { data });
    },
  },
});
