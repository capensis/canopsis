import Vue from 'vue';
import { get } from 'lodash';

import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import request from '@/services/request';

import i18n from '@/i18n';

import { serviceSchema } from '@/store/schemas';
import { createEntityModule } from '@/store/plugins/entities';

import entityModule from './entity';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default createEntityModule({
  route: API_ROUTES.service,
  entityType: ENTITIES_TYPES.service,
  dataPreparer: d => d.data,
  withFetchingParams: true,
  withMeta: true,
  types,
}, {
  modules: {
    entity: entityModule,
  },
  state: {
    widgets: {},
  },
  getters: {
    getListByWidgetId: (state, getters, rootState, rootGetters) => widgetId => rootGetters['entities/getList'](
      ENTITIES_TYPES.service,
      get(state.widgets[widgetId], 'allIds', []),
    ),
    getPendingByWidgetId: state => widgetId => get(state.widgets[widgetId], 'pending'),
    getErrorByWidgetId: state => widgetId => get(state.widgets[widgetId], 'error'),
    getItem: (state, getters, rootState, rootGetters) => id => rootGetters['entities/getItem'](
      ENTITIES_TYPES.service,
      id,
    ),
  },
  mutations: {
    [types.FETCH_LIST](state, { widgetId }) {
      Vue.setSeveral(state.widgets, widgetId, { pending: true, error: null });
    },
    [types.FETCH_LIST_COMPLETED](state, { widgetId, allIds }) {
      Vue.setSeveral(state.widgets, widgetId, { allIds, pending: false });
    },
    [types.FETCH_LIST_FAILED](state, { widgetId, error = {} }) {
      Vue.setSeveral(state.widgets, widgetId, { error, pending: false });
    },
  },
  actions: {
    async fetchList({ dispatch, commit }, { widgetId, params } = {}) {
      try {
        commit(types.FETCH_LIST, { widgetId });

        const { normalizedData, data } = await dispatch('entities/fetch', {
          route: API_ROUTES.weatherService,
          schema: [serviceSchema],
          dataPreparer: d => d.data,
          params,
        }, { root: true });

        commit(types.FETCH_LIST_COMPLETED, {
          widgetId,
          allIds: normalizedData.result,
          ...data,
        });
      } catch (err) {
        console.error(err);

        commit(types.FETCH_LIST_FAILED, { widgetId, error: err });

        await dispatch('popups/error', { text: i18n.t('errors.default') }, { root: true });
      }
    },

    fetchDependenciesWithoutStore(context, { id, params = {} }) {
      return request.get(API_ROUTES.serviceDependencies, {
        params: { ...params, _id: id },
      });
    },

    fetchImpactsWithoutStore(context, { id, params = {} }) {
      return request.get(API_ROUTES.serviceImpacts, {
        params: { ...params, _id: id },
      });
    },

    fetchAlarmsWithoutStore(context, { id, params }) {
      return request.get(`${API_ROUTES.serviceAlarms}/${id}`, { params });
    },

    fetchItemWithoutStore(context, { id, params }) {
      return request.get(`${API_ROUTES.service}/${id}`, { params });
    },

    fetchInfosKeysWithoutStore(context, { params }) {
      return request.get(API_ROUTES.entityInfosDictionaryKeys, { params });
    },
  },
});
