import Vue from 'vue';
import { get } from 'lodash';

import i18n from '@/i18n';
import request from '@/services/request';
import { API_ROUTES } from '@/config';
import { watcherSchema } from '@/store/schemas';
import { ENTITIES_TYPES, WIDGET_TYPES } from '@/constants';

import watcherEntityModule from './entity';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default {
  namespaced: true,
  modules: {
    entity: watcherEntityModule,
  },
  state: {
    widgets: {},
  },
  getters: {
    getListByWidgetId: (state, getters, rootState, rootGetters) => widgetId =>
      rootGetters['entities/getList'](ENTITIES_TYPES.watcher, get(state.widgets[widgetId], 'allIds', [])),
    getPendingByWidgetId: state => widgetId => get(state.widgets[widgetId], 'pending'),
    getErrorByWidgetId: state => widgetId => get(state.widgets[widgetId], 'error'),
    getItem: (state, getters, rootState, rootGetters) => id =>
      rootGetters['entities/getItem'](ENTITIES_TYPES.watcher, id),
  },
  mutations: {
    [types.FETCH_LIST](state, { widgetId }) {
      Vue.set(state.widgets, widgetId, {
        ...state.widgets[widgetId],
        pending: true,
      });
    },
    [types.FETCH_LIST_COMPLETED](state, { widgetId, allIds }) {
      Vue.set(state.widgets, widgetId, {
        ...state.widgets[widgetId],
        pending: false,
        allIds,
      });
    },
    [types.FETCH_LIST_FAILED](state, { widgetId, error }) {
      Vue.set(state.widgets, widgetId, {
        ...state.widgets[widgetId],
        pending: false,
        error,
      });
    },
  },
  actions: {
    createWatcher(context, { data }) {
      return request.put(API_ROUTES.createEntity, { entity: JSON.stringify(data) });
    },

    createWatcherNg(context, { data }) {
      return request.post(API_ROUTES.watcherng, data);
    },

    editWatcher(context, { data }) {
      return request.put(API_ROUTES.context, { entity: data, _type: WIDGET_TYPES.context });
    },

    editWatcherNg(context, { data }) {
      return request.put(`${API_ROUTES.watcherng}/${data._id}`, data);
    },

    async fetchList({ dispatch, commit }, { widgetId, params, filter } = {}) {
      try {
        const requestFilter = filter || '{}';

        commit(types.FETCH_LIST, { widgetId });

        const { normalizedData } = await dispatch('entities/fetch', {
          route: `${API_ROUTES.weatherWatcher}/${requestFilter}`,
          schema: [watcherSchema],
          params,
        }, { root: true });

        commit(types.FETCH_LIST_COMPLETED, {
          widgetId,
          allIds: normalizedData.result,
        });
      } catch (err) {
        commit(types.FETCH_LIST_FAILED, {
          widgetId,
          error: err,
        });

        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });
      }
    },
  },
};
