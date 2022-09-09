import Vue from 'vue';
import { merge, get } from 'lodash';

import request, { useRequestCancelling } from '@/services/request';
import i18n from '@/i18n';
import { alarmSchema } from '@/store/schemas';
import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',

  EXPORT_LIST: 'EXPORT_LIST',
  EXPORT_LIST_COMPLETED: 'EXPORT_LIST_COMPLETED',
  EXPORT_LIST_FAILED: 'EXPORT_LIST_FAILED',

  DOWNLOAD_LIST_COMPLETED: 'DOWNLOAD_LIST_COMPLETED',
};

export default {
  namespaced: true,
  state: {
    widgets: {},
  },
  getters: {
    getListByWidgetId: (state, getters, rootState, rootGetters) => widgetId => rootGetters['entities/getList'](
      ENTITIES_TYPES.alarm,
      get(state.widgets[widgetId], 'allIds', []),
    ),

    getMetaByWidgetId: state => widgetId => get(state.widgets[widgetId], 'meta', {}),
    getPendingByWidgetId: state => widgetId => get(state.widgets[widgetId], 'pending', false),
    getExportByWidgetId: state => widgetId => get(state.widgets[widgetId], 'exportData'),

    getItem: (state, getters, rootState, rootGetters) => id => rootGetters['entities/getItem'](
      ENTITIES_TYPES.alarm,
      id,
    ),
    getList: (state, getters, rootState, rootGetters) => ids => rootGetters['entities/getList'](
      ENTITIES_TYPES.alarm,
      ids,
    ),

    getFetchingParamsByWidgetId: state => widgetId => get(state.widgets[widgetId], 'fetchingParams'),
  },
  mutations: {
    [types.FETCH_LIST](state, { widgetId, params }) {
      Vue.setSeveral(state.widgets, widgetId, { pending: true, fetchingParams: params });
    },
    [types.FETCH_LIST_COMPLETED](state, { widgetId, allIds, meta }) {
      Vue.setSeveral(state.widgets, widgetId, { allIds, meta, pending: false });
    },
    [types.FETCH_LIST_FAILED](state, { widgetId }) {
      Vue.setSeveral(state.widgets, widgetId, { pending: false });
    },
    [types.EXPORT_LIST_COMPLETED](state, { widgetId, exportData }) {
      Vue.setSeveral(state.widgets, widgetId, { exportData });
    },
    [types.DOWNLOAD_LIST_COMPLETED](state, { widgetId }) {
      Vue.delete(state.widgets[widgetId], 'exportData');
    },
  },
  actions: {
    async fetchListWithoutStore(context, { params }) {
      return request.get(API_ROUTES.alarmList, { params });
    },

    async fetchList({ commit, dispatch }, { widgetId, params, withoutPending } = {}) {
      try {
        await useRequestCancelling(async (source) => {
          if (!withoutPending) {
            commit(types.FETCH_LIST, { widgetId, params });
          }

          const { normalizedData, data } = await dispatch('entities/fetch', {
            route: API_ROUTES.alarmList,
            schema: [alarmSchema],
            params,
            cancelToken: source.token,
            dataPreparer: d => d.data,
          }, { root: true });

          commit(types.FETCH_LIST_COMPLETED, {
            widgetId,
            allIds: normalizedData.result,
            meta: data.meta,
          });
        }, `alarms-list-${widgetId}`);
      } catch (err) {
        await dispatch('popups/error', { text: i18n.t('errors.default') }, { root: true });

        commit(types.FETCH_LIST_FAILED, { widgetId });
      }
    },

    fetchListWithPreviousParams({ dispatch, state }, { widgetId }) {
      return dispatch('fetchList', {
        widgetId,
        params: get(state, ['widgets', widgetId, 'fetchingParams'], {}),
        withoutPending: true,
      });
    },

    async fetchItem({ dispatch }, { id, params, dataPreparer = d => d.data }) {
      try {
        const paramsWithItemId = merge(params, { filter: { _id: id } });

        await dispatch('entities/fetch', {
          route: API_ROUTES.alarmList,
          schema: [alarmSchema],
          params: paramsWithItemId,
          dataPreparer,
        }, { root: true });
      } catch (err) {
        console.error(err);
      }
    },

    async createAlarmsListExport({ commit }, { widgetId, data = {} }) {
      const exportData = await request.post(API_ROUTES.alarmListExport, data);

      commit(types.EXPORT_LIST_COMPLETED, {
        widgetId,
        exportData,
      });

      return exportData;
    },

    async fetchAlarmsListExport({ commit }, { params, id, widgetId }) {
      const exportData = await request.get(`${API_ROUTES.alarmListExport}/${id}`, { params });

      commit(types.EXPORT_LIST_COMPLETED, {
        widgetId,
        exportData,
      });

      return exportData;
    },

    async fetchAlarmsListCsvFile({ commit }, { params, id, widgetId }) {
      const csvData = await request.get(`${API_ROUTES.alarmListExport}/${id}/download`, { params });

      commit(types.DOWNLOAD_LIST_COMPLETED, { widgetId });

      return csvData;
    },
  },
};
