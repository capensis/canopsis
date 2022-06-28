import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import request from '@/services/request';

import { createEntityModule } from '@/store/plugins/entities';

export default createEntityModule({
  route: API_ROUTES.scenario.scenarios,
  entityType: ENTITIES_TYPES.scenario,
  dataPreparer: d => d.data,
  withMeta: true,
}, {
  actions: {
    checkPriority(context, { data }) {
      return request.post(API_ROUTES.scenario.checkPriority, data);
    },

    fetchMinimalPriority() {
      return request.get(API_ROUTES.scenario.minimalPriority);
    },
  },
});
