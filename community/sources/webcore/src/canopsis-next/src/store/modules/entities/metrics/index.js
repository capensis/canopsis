import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    fetchSliMetricsWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.metrics.sli, { params });
    },

    fetchRatingMetricsWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.metrics.rating, { params });
    },
  },
};
