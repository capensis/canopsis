import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.flappingRules,
}, {
  actions: {
    bulkEnable(context, { data }) {
      return request.put(API_ROUTES.bulkFlappingRulesEnable, data);
    },

    bulkDisable(context, { data }) {
      return request.put(API_ROUTES.bulkFlappingRulesDisable, data);
    },

    bulkRemove(context, { data }) {
      return request.delete(API_ROUTES.bulkFlappingRules, { data });
    },
  },
});
