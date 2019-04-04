import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import request from '@/services/request';

import { createEntityModule } from '@/store/plugins/entities';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default createEntityModule({
  types,
  route: API_ROUTES.snmpRule.list,
  entityType: ENTITIES_TYPES.snmpRule,
}, {
  actions: {
    create(context, data) {
      return request.post(API_ROUTES.snmpRule.create, data);
    },
  },
});
