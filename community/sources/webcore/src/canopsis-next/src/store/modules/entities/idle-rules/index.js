import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.idleRules,
}, {
  actions: {
    bulkEnable(context, { data }) {
      return request.put(API_ROUTES.bulkIdleRulesEnable, data);
    },

    bulkDisable(context, { data }) {
      return request.put(API_ROUTES.bulkIdleRulesDisable, data);
    },

    bulkRemove(context, { data }) {
      return request.delete(API_ROUTES.bulkIdleRules, { data });
    },
  },
});
