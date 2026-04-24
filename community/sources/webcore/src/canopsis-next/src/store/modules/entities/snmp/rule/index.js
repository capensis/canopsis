import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.snmpRule,
  entityType: ENTITIES_TYPES.snmpRule,
  dataPreparer: d => d.data,
}, {
  actions: {
    bulkEnable(context, { data }) {
      return request.put(API_ROUTES.bulkSnmpRuleEnable, data);
    },

    bulkDisable(context, { data }) {
      return request.put(API_ROUTES.bulkSnmpRuleDisable, data);
    },

    bulkRemove(context, { data }) {
      return request.delete(API_ROUTES.bulkSnmpRule, { data });
    },
  },
});
