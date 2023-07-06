import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import request from '@/services/request';

import { createEntityModule } from '@/store/plugins/entities';

export default createEntityModule({
  route: API_ROUTES.alarmTag.list,
  entityType: ENTITIES_TYPES.alarmTag,
  dataPreparer: d => d.data,
  withMeta: true,
}, {
  actions: {
    bulkRemove(context, { data }) {
      return request.delete(API_ROUTES.alarmTag.bulkList, { data });
    },
  },
});
