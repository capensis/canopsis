import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

import labelModule from './label';

export default createCRUDModule({
  route: API_ROUTES.alarmTag.list,
  withWithoutStore: true,
}, {
  modules: {
    label: labelModule,
  },
  actions: {
    bulkRemove(context, { data }) {
      return request.delete(API_ROUTES.alarmTag.bulkList, { data });
    },
  },
});
