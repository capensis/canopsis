import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.metaAlarmRule,
  withWithoutStore: true,
}, {
  actions: {
    bulkEnable(context, { data }) {
      return request.put(API_ROUTES.bulkMetaAlarmRuleEnable, data);
    },

    bulkDisable(context, { data }) {
      return request.put(API_ROUTES.bulkMetaAlarmRuleDisable, data);
    },

    bulkRemove(context, { data }) {
      return request.delete(API_ROUTES.bulkMetaAlarmRule, { data });
    },
  },
});
