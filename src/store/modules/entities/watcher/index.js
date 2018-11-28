import Vue from 'vue';
import get from 'lodash/get';

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
    [types.FETCH_LIST_FAILED](state, { widgetId }) {
      Vue.set(state.widgets, widgetId, {
        ...state.widgets[widgetId],
        pending: false,
      });
    },
  },
  actions: {
    async create(context, { data }) {
      await request.put(API_ROUTES.createEntity, { entity: JSON.stringify(data) });
    },
    async edit(context, { data }) {
      await request.put(API_ROUTES.context, { entity: data, _type: WIDGET_TYPES.context });
    },

    async remove({ dispatch }, { id } = {}) {
      try {
        await request.delete(API_ROUTES.watcher, { params: { watcher_id: id } });

        await dispatch('entities/removeFromStore', {
          id,
          type: ENTITIES_TYPES.watcher,
        }, { root: true });
      } catch (err) {
        console.warn(err);
      }
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
        commit(types.FETCH_LIST_FAILED, { widgetId });

        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });
      }
    },
  },
};
