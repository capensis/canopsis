import Vue from 'vue';
import { get } from 'lodash';

import { entitySchema } from '@/store/schemas';
import request from '@/services/request';
import i18n from '@/i18n';
import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES, REQUEST_METHODS, WIDGET_TYPES } from '@/constants';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',

  EDIT_FAILED: 'EDIT_FAILED',

  CREATION_FAILED: 'CREATION_FAILED',

  FETCH_GENERAL_LIST: 'FETCH_GENERAL_LIST',
  FETCH_GENERAL_LIST_COMPLETED: 'FETCH_GENERAL_LIST_COMPLETED',
  FETCH_GENERAL_LIST_FAILED: 'FETCH_GENERAL_LIST_FAILED',

  EXPORT_LIST: 'EXPORT_LIST',
  EXPORT_LIST_COMPLETED: 'EXPORT_LIST_COMPLETED',
  EXPORT_LIST_FAILED: 'EXPORT_LIST_FAILED',

  DOWNLOAD_LIST_COMPLETED: 'DOWNLOAD_LIST_COMPLETED',
};

export default {
  namespaced: true,
  state: {
    widgets: {},
    allIdsGeneralList: [],
    pendingGeneralList: false,
    fetchingParamsGeneralList: false,
  },
  getters: {
    getListByWidgetId: (state, getters, rootState, rootGetters) => widgetId =>
      rootGetters['entities/getList'](ENTITIES_TYPES.entity, get(state.widgets[widgetId], 'allIds', [])),

    getMetaByWidgetId: state => widgetId => get(state.widgets[widgetId], 'meta', {}),
    getPendingByWidgetId: state => widgetId => get(state.widgets[widgetId], 'pending'),
    getExportByWidgetId: state => widgetId => get(state.widgets[widgetId], 'exportData'),

    allIdsGeneralList: state => state.allIds,
    itemsGeneralList: (state, getters, rootState, rootGetters) => rootGetters['entities/getList']('entity', state.allIdsGeneralList),
    pendingGeneralList: state => state.pendingGeneralList,
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
    [types.FETCH_GENERAL_LIST](state, { params }) {
      state.pendingGeneralList = true;
      state.fetchingParamsGeneralList = params;
    },
    [types.FETCH_GENERAL_LIST_COMPLETED](state, { allIds }) {
      state.allIdsGeneralList = allIds;
      state.pendingGeneralList = false;
    },
    [types.FETCH_GENERAL_LIST_FAILED](state) {
      state.pendingGeneralList = false;
    },
  },
  actions: {
    async fetchListWithoutStore({ dispatch }, { params } = {}) {
      try {
        const { total, data } = await request.post(API_ROUTES.context, null, { params });

        return { total, entities: data };
      } catch (err) {
        await dispatch('popups/error', { text: i18n.t('errors.default') }, { root: true });

        return { total: 0, entities: [] };
      }
    },

    fetch({ dispatch }, { params } = {}) {
      return dispatch('entities/fetch', {
        route: API_ROUTES.context,
        schema: [entitySchema],
        method: REQUEST_METHODS.post,
        params,
        dataPreparer: d => d.data,
      }, { root: true });
    },

    async fetchList({ commit, dispatch }, { widgetId, params } = {}) {
      try {
        commit(types.FETCH_LIST, { widgetId, params });

        const { normalizedData, data } = await dispatch('fetch', { params });

        commit(types.FETCH_LIST_COMPLETED, {
          widgetId,
          allIds: normalizedData.result,
          meta: {
            total: data.total,
          },
        });
      } catch (err) {
        console.error(err);

        commit(types.FETCH_LIST_FAILED, { widgetId });
        await dispatch('popups/error', { text: i18n.t('errors.default') }, { root: true });
      }
    },

    async fetchGeneralList({ commit, dispatch }, { params } = {}) {
      try {
        commit(types.FETCH_GENERAL_LIST, { params });

        const { normalizedData } = await dispatch('fetch', { params });

        commit(types.FETCH_GENERAL_LIST_COMPLETED, {
          allIds: normalizedData.result,
        });
      } catch (err) {
        console.error(err);
        commit(types.FETCH_GENERAL_LIST_FAILED);
        await dispatch('popups/error', { text: i18n.t('errors.default') }, { root: true });
      }
    },

    refreshLists({ dispatch, getters, state }) {
      const widgetsIds = Object.keys(state.widgets);
      const fetchRequests = widgetsIds.map(widgetId => dispatch('fetchList', {
        widgetId,
        params: getters.getFetchingParamsByWidgetId(widgetId),
      }));

      return Promise.all(fetchRequests);
    },

    create(context, { data }) {
      // Need this special syntax for request params for the backend to handle it
      return request.put(API_ROUTES.createEntity, { entity: JSON.stringify(data) });
    },

    update(context, { data }) {
      return request.put(API_ROUTES.context, { entity: data, _type: WIDGET_TYPES.context });
    },

    async remove({ dispatch }, { id } = {}) {
      try {
        await request.delete(API_ROUTES.context, { params: { ids: id } });
        await dispatch('entities/removeFromStore', {
          id,
          type: ENTITIES_TYPES.entity,
        }, { root: true });
      } catch (err) {
        console.error(err);
      }
    },

    async createContextExport({ commit }, { widgetId, params }) {
      const exportData = await request.post(API_ROUTES.contextExport, {}, { params });

      commit(types.EXPORT_LIST_COMPLETED, {
        widgetId,
        exportData,
      });

      return exportData;
    },

    fetchContextExport({ commit }, { params, widgetId, id }) {
      const exportData = request.get(`${API_ROUTES.contextExport}/${id}`, { params });

      commit(types.EXPORT_LIST_COMPLETED, {
        widgetId,
        exportData,
      });

      return exportData;
    },

    async fetchContextCsvFile({ commit }, { params, id, widgetId }) {
      const csvData = await request.get(`${API_ROUTES.contextExport}/${id}/download`, { params });

      commit(types.DOWNLOAD_LIST_COMPLETED, { widgetId });

      return csvData;
    },
  },
};
