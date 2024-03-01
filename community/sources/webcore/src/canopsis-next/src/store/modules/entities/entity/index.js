import Vue from 'vue';
import { get } from 'lodash';

import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import request from '@/services/request';

import i18n from '@/i18n';

import { entitySchema } from '@/store/schemas';
import { createEntityModule } from '@/store/plugins/entities';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',

  DOWNLOAD_LIST_COMPLETED: 'DOWNLOAD_LIST_COMPLETED',
};

export default createEntityModule({
  route: API_ROUTES.entity,
  entityType: ENTITIES_TYPES.entity,
  dataPreparer: d => d.data,
  withFetchingParams: true,
  withMeta: true,
  types,
}, {
  state: {
    widgets: {},
  },
  getters: {
    getListByWidgetId: (state, getters, rootState, rootGetters) => widgetId => rootGetters['entities/getList'](
      ENTITIES_TYPES.entity,
      get(state.widgets[widgetId], 'allIds', []),
    ),

    getMetaByWidgetId: state => widgetId => get(state.widgets[widgetId], 'meta', {}),
    getPendingByWidgetId: state => widgetId => get(state.widgets[widgetId], 'pending'),

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
    async fetchList({ commit, dispatch }, { widgetId, params } = {}) {
      try {
        commit(types.FETCH_LIST, { widgetId, params });

        const { normalizedData, data } = await dispatch('entities/fetch', {
          route: API_ROUTES.entity,
          params,
          dataPreparer: d => d.data,
          schema: [entitySchema],
        }, { root: true });

        commit(types.FETCH_LIST_COMPLETED, {
          widgetId,
          allIds: normalizedData.result,
          ...data,
        });
      } catch (err) {
        console.error(err);

        commit(types.FETCH_LIST_FAILED, { widgetId });
        await dispatch('popups/error', { text: i18n.t('errors.default') }, { root: true });
      }
    },

    fetchListWithPreviousParams({ dispatch, getters, state }) {
      const widgetsIds = Object.keys(state.widgets);
      const fetchRequests = widgetsIds.map(widgetId => dispatch('fetchList', {
        widgetId,
        params: getters.getFetchingParamsByWidgetId(widgetId),
      }));

      return Promise.all(fetchRequests);
    },

    update(context, { data, id }) {
      return request.put(API_ROUTES.entityBasics, data, {
        params: { _id: id },
      });
    },

    create(context, { data, id }) {
      return request.post(API_ROUTES.entityBasics, data, {
        params: { _id: id },
      });
    },

    remove(context, { id }) {
      return request.delete(API_ROUTES.entityBasics, {
        params: { _id: id },
      });
    },

    fetchBasicEntityWithoutStore(context, { id } = {}) {
      return request.get(API_ROUTES.entityBasics, {
        params: { _id: id },
      });
    },

    fetchContextGraphWithoutStore(context, { id } = {}) {
      return request.get(API_ROUTES.entityContextGraph, {
        params: { _id: id },
      });
    },

    checkStateSetting(context, { data }) {
      return request.post(API_ROUTES.entityCheckStateSetting, data);
    },

    fetchStateSettingWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.entityStateSetting, { params });
    },

    async fetchAvailabilityWithoutStore(context, { id, params } = {}) {
      // return request.get(API_ROUTES.entityAvailability, { params });
      // eslint-disable-next-line no-console
      console.info(id, params);
      /**
       * TODO: Should be replaced on real fetch function
       */
      await new Promise(r => setTimeout(r, 2000));

      const minDate = new Date();
      minDate.setDate(minDate.getDate() - 3);

      return {
        availability: {
          uptime: Math.round(Math.random() * 100000),
          downtime: Math.round(Math.random() * 100000),
          inactive_time: Math.round(Math.random() * 1000),
        },
        min_date: Math.round(minDate.getTime() / 1000),
      };
    },

    fetchListWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.entity, { params });
    },

    async createContextExport(context, { data = {} }) {
      return request.post(API_ROUTES.contextExport, data);
    },

    fetchContextExport(context, { params, id }) {
      return request.get(`${API_ROUTES.contextExport}/${id}`, { params });
    },

    archiveDisabledEntitiesData(context, { data }) {
      return request.post(`${API_ROUTES.entity}/archive-disabled`, data);
    },

    archiveUnlinkedEntitiesData(context, { data }) {
      return request.post(`${API_ROUTES.entity}/archive-unlinked`, data);
    },

    cleanArchivedEntitiesData() {
      return request.post(`${API_ROUTES.entity}/clean-archived`);
    },

    bulkEnable(context, { data }) {
      return request.put(API_ROUTES.bulkEntitiesEnable, data);
    },

    bulkDisable(context, { data }) {
      return request.put(API_ROUTES.bulkEntitiesDisable, data);
    },
  },
});
