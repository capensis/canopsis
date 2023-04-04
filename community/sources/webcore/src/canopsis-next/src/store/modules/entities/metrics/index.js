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

    createKpiAlarmAggregateExport(context, { data }) {
      return request.post(API_ROUTES.metrics.exportAggregate, data);
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

    fetchExternalMetricsListWithoutStore(context, { params }) {
      /* TODO: Should be added when backend part will be finished  */
      return new Promise((resolve) => {
        setTimeout(() => resolve({
          data: [{
            _id: 'instance2/cpu_surveillance',
            name: 'instance2/cpu_surveillance',
          }],
        }), 2000, params);
      });
    },
  },
};
