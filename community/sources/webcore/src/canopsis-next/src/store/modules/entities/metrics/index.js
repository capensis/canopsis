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

    fetchAlarmsMetricsWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.metrics.alarm, { params });
    },

    createKpiAlarmExport(context, { data }) {
      return request.post(API_ROUTES.metrics.exportAlarm, null, { params: data });
    },

    createKpiRatingExport(context, { data }) {
      return request.post(API_ROUTES.metrics.exportRating, null, { params: data });
    },

    createKpiSliExport(context, { data }) {
      return request.post(API_ROUTES.metrics.exportSli, null, { params: data });
    },

    fetchMetricExport(context, { id }) {
      return request.get(`${API_ROUTES.metrics.exportMetric}/${id}`);
    },
  },
};
