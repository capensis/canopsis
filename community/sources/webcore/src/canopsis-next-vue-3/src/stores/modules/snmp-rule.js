import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDStore } from '@/stores/plugins/crud';

export const useSnmpRuleStore = createCRUDStore({
  id: 'snmpRule',
  route: API_ROUTES.snmpRule,
  dataPreparer: (response) => response?.data,
  withFetchingParams: true,
  withMeta: true,
  withWithoutStore: true,
  actions: {
    bulkEnable({ data } = {}) {
      return request.put(API_ROUTES.bulkSnmpRuleEnable, data);
    },

    bulkDisable({ data } = {}) {
      return request.put(API_ROUTES.bulkSnmpRuleDisable, data);
    },

    bulkRemove({ data } = {}) {
      return request.delete(API_ROUTES.bulkSnmpRule, { data });
    },
  },
});
