import request from '@/services/request';

import { createCRUDStore } from './crud';

export const createJobManagementStore = ({ id, route, bulkStopRoute }) => createCRUDStore({
  id,
  route,
  withFetchingParams: true,
  actions: {
    stop({ data } = {}) {
      return request.post(bulkStopRoute, data);
    },
  },
});
