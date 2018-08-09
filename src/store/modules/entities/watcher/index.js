import Vue from 'vue';
import get from 'lodash/get';
import request from '@/services/request';
import { API_ROUTES } from '@/config';
import { watcherSchema } from '@/store/schemas';
import i18n from '@/i18n';
import { ENTITIES_TYPES, WIDGET_TYPES } from '@/constants';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default {
  namespaced: true,
  state: {
    widgets: {},
  },
  getters: {
    getListByWidgetId: (state, getters, rootState, rootGetters) => widgetId =>
      rootGetters['entities/getList'](ENTITIES_TYPES.watcher, get(state.widgets[widgetId], 'allIds', [])),
    getMetaByWidgetId: state => widgetId => get(state.widgets[widgetId], 'meta', {}),
    getPendingByWidgetId: state => widgetId => get(state.widgets[widgetId], 'pending'),
  },
  mutations: {
    [types.FETCH_LIST](state, { widgetId, params }) {
      Vue.set(state.widgets, widgetId, {
        ...state.widgets[widgetId],
        pending: true,
        fetchingParams: params,
      });
    },
    [types.FETCH_LIST_COMPLETED](state, { widgetId, allIds, meta }) {
      Vue.set(state.widgets, widgetId, {
        ...state.widgets[widgetId],
        pending: false,
        allIds,
        meta,
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
    async create(context, params = {}) {
      try {
        await request.post(API_ROUTES.watcher, params);
      } catch (err) {
        console.warn(err);
      }
    },
    async edit(context, { data }) {
      try {
        await request.put(API_ROUTES.context, { entity: data, _type: WIDGET_TYPES.context });
      } catch (err) {
        console.warn(err);
      }
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

    async fetchList({ dispatch, commit }, { widgetId, params, filter = {} } = {}) {
      try {
        const { normalizedData, data } = await dispatch('entities/fetch', {
          route: `${API_ROUTES.weatherWatcher}/${JSON.stringify(filter)}`,
          schema: [watcherSchema],
          params,
          dataPreparer: d => d,
        }, { root: true });
        commit(types.FETCH_LIST_COMPLETED, {
          widgetId,
          allIds: normalizedData.result,
          meta: {
            first: data.first,
            last: data.last,
            total: data.total,
          },
        });
      } catch (err) {
        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });
        commit(types.FETCH_LIST_FAILED, { widgetId });
      }
    },
  },
};
