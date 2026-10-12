import { API_ROUTES } from '@/config';

import { createCRUDStore } from '@/stores/plugins/crud';

export const useAlarmTagStore = createCRUDStore({
  id: 'alarmTag',
  route: API_ROUTES.alarmTag.list,
  bulkRoute: API_ROUTES.alarmTag.bulkList,
  withFetchingParams: true,
  withMeta: true,
  withWithoutStore: true,
});
