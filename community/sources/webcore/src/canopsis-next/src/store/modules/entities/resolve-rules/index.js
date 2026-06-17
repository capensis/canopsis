import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.resolveRules,
}, {
  actions: {
    bulkEnable(context, { data }) {
      return request.put(API_ROUTES.bulkResolveRulesEnable, data);
    },

    bulkDisable(context, { data }) {
      return request.put(API_ROUTES.bulkResolveRulesDisable, data);
    },

    bulkRemove(context, { data }) {
      return request.delete(API_ROUTES.bulkResolveRules, { data });
    },
  },
});
