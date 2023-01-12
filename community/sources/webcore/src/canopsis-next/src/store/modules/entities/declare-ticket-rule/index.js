import { API_ROUTES } from '@/config';

import { ENTITIES_TYPES } from '@/constants';

import { createEntityModule } from '@/store/plugins/entities';
import request from '@/services/request';

export default createEntityModule({
  route: API_ROUTES.declareTicketRules,
  entityType: ENTITIES_TYPES.declareTicketRule,
  dataPreparer: d => d.data,
  withMeta: true,
}, {
  actions: {
    bulkRemove(context, { data }) {
      return request.delete(API_ROUTES.bulkDeclareTicketRules, { data });
    },
  },
});
