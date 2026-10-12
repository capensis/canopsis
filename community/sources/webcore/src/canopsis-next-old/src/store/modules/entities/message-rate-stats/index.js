import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    fetchMessageRateStatsWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.messageRateStats, { params });
    },
  },
};
