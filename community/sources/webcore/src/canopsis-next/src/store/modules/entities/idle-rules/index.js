import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import { createEntityModule } from '@/store/plugins/entities';
import request from '@/services/request';

export default createEntityModule({
  route: API_ROUTES.idleRules,
  entityType: ENTITIES_TYPES.idleRules,
  dataPreparer: d => d.data,
  withMeta: true,
}, {
  actions: {
    fetchEntitiesCountWithoutStore(context, { data }) {
      return request.post(API_ROUTES.idleRulesCount, data);
    },
  },
});