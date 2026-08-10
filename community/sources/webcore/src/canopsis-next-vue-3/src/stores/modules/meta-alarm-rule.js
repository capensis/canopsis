import { API_ROUTES } from '@/config';
import request from '@/services/request';
import { createCRUDStore } from '@/stores/plugins/crud';

export const useMetaAlarmRuleStore = createCRUDStore({
  id: 'metaAlarmRule',
  route: API_ROUTES.metaAlarmRule,
  withFetchingParams: true,
  withMeta: true,
  withWithoutStore: true,
  actions: {
    bulkEnable({ data } = {}) {
      return request.put(API_ROUTES.bulkMetaAlarmRuleEnable, data);
    },

    bulkDisable({ data } = {}) {
      return request.put(API_ROUTES.bulkMetaAlarmRuleDisable, data);
    },

    bulkRemove({ data } = {}) {
      return request.delete(API_ROUTES.bulkMetaAlarmRule, { data });
    },

    checkPatternsAlarmsCount({ data } = {}) {
      return request.post(API_ROUTES.metaAlarmRuleAlarmCount, data);
    },
  },
});
