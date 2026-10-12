import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    fetchItemGanttIntervalsWithoutStore(context, { id, params } = {}) {
      return request.get(API_ROUTES.junit.entityGantt, { params: { ...params, _id: id } });
    },
  },
};
