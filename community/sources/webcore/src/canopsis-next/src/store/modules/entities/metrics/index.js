import { API_ROUTES } from '@/config';

import request from '@/services/request';

const types = {
  FETCH_EXTERNAL_METRICS_LIST: 'FETCH_EXTERNAL_METRICS_LIST',
  FETCH_EXTERNAL_METRICS_LIST_COMPLETED: 'FETCH_EXTERNAL_METRICS_LIST_COMPLETED',
  FETCH_EXTERNAL_METRICS_LIST_FAILED: 'FETCH_EXTERNAL_METRICS_LIST_FAILED',
};

export default {
  namespaced: true,
  state: {
    externalMetrics: [],
    pending: false,
  },
  getters: {
    externalMetrics: state => state.externalMetrics,
    pending: state => state.pending,
  },
  mutations: {
    [types.FETCH_EXTERNAL_METRICS_LIST](state) {
      state.pending = true;
    },
    [types.FETCH_EXTERNAL_METRICS_LIST_COMPLETED](state, externalMetrics) {
      state.pending = false;
      state.externalMetrics = externalMetrics;
    },
    [types.FETCH_EXTERNAL_METRICS_LIST_FAILED](state) {
      state.pending = false;
    },
  },
  actions: {
    fetchSliMetricsWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.metrics.sli, { params });
    },

    fetchRatingMetricsWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.metrics.rating, { params });
    },

    fetchAlarmsMetricsWithoutStore(context, { params } = {}) {
      return request.post(API_ROUTES.metrics.alarm, params);
    },

    createKpiAlarmExport(context, { data }) {
      return request.post(API_ROUTES.metrics.exportAlarm, null, { params: data });
    },

    createKpiAlarmAggregateExport(context, { data }) {
      return request.post(API_ROUTES.metrics.exportAggregate, data);
    },

    createRemediationExport(context, { data }) {
      return request.post(API_ROUTES.metrics.exportRemediation, data);
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

    async fetchExternalMetricsList({ commit }, { params }) {
      commit(types.FETCH_EXTERNAL_METRICS_LIST);

      try {
        const { data } = await request.get(API_ROUTES.metrics.perfDataMetrics, { params });

        commit(types.FETCH_EXTERNAL_METRICS_LIST_COMPLETED, data);
      } catch (err) {
        commit(types.FETCH_EXTERNAL_METRICS_LIST_FAILED);
      }
    },
  },
};
