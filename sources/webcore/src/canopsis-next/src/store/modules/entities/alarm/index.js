import Vue from 'vue';
import { merge, get } from 'lodash';

import request from '@/services/request';
import i18n from '@/i18n';
import { alarmSchema } from '@/store/schemas';
import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

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
      rootGetters['entities/getList'](ENTITIES_TYPES.alarm, get(state.widgets[widgetId], 'allIds', [])),

    getMetaByWidgetId: state => widgetId => get(state.widgets[widgetId], 'meta', {}),
    getPendingByWidgetId: state => widgetId => get(state.widgets[widgetId], 'pending'),
    getItem: (state, getters, rootState, rootGetters) => id =>
      rootGetters['entities/getItem'](ENTITIES_TYPES.alarm, id),
  },
  mutations: {
    [types.FETCH_LIST](state, { widgetId, params }) {
      if (!state.widgets[widgetId]) {
        Vue.set(state.widgets, widgetId, {
          pending: true,
          fetchingParams: params,
        });
      } else {
        Vue.set(state.widgets[widgetId], 'pending', true);
        Vue.set(state.widgets[widgetId], 'fetchingParams', params);
      }
    },
    [types.FETCH_LIST_COMPLETED](state, { widgetId, allIds, meta }) {
      if (!state.widgets[widgetId]) {
        Vue.set(state.widgets, widgetId, {
          pending: false,
          allIds,
          meta,
        });
      } else {
        Vue.set(state.widgets[widgetId], 'pending', false);
        Vue.set(state.widgets[widgetId], 'allIds', allIds);
        Vue.set(state.widgets[widgetId], 'meta', meta);
      }
    },
    [types.FETCH_LIST_FAILED](state, { widgetId }) {
      Vue.set(state.widgets, widgetId, {
        ...state.widgets[widgetId],
        pending: false,
      });
    },
  },
  actions: {
    async fetchListWithoutStore({ dispatch }, { params, withoutCatch = false }) {
      try {
        const { data: [result] } = await request.get(API_ROUTES.alarmList, { params });

        return result;
      } catch (err) {
        if (withoutCatch) {
          throw err;
        }

        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });

        return { alarms: [], total: 0 };
      }
    },
    async fetchList({ commit, dispatch }, { widgetId, params, withoutPending } = {}) {
      try {
        if (!withoutPending) {
          commit(types.FETCH_LIST, { widgetId, params });
        }

        const { normalizedData, data } = await dispatch('entities/fetch', {
          route: API_ROUTES.alarmList,
          schema: [alarmSchema],
          params,
          dataPreparer: d => d.data[0].alarms,
        }, { root: true });

        const total = data.data[0].total ? data.data[0].total : normalizedData.result.length;

        commit(types.FETCH_LIST_COMPLETED, {
          widgetId,
          allIds: normalizedData.result,
          meta: {
            total,
            first: data.data[0].first,
            last: data.data[0].last,
          },
        });
      } catch (err) {
        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });

        commit(types.FETCH_LIST_FAILED, { widgetId });
      }
    },

    fetchListWithPreviousParams({ dispatch, state }, { widgetId }) {
      return dispatch('fetchList', {
        widgetId,
        params: state.widgets[widgetId].fetchingParams,
        withoutPending: true,
      });
    },

    async fetchItem({ dispatch }, { id, params }) {
      try {
        const paramsWithItemId = merge(params, { filter: { _id: id } });

        await dispatch('entities/fetch', {
          route: API_ROUTES.alarmList,
          schema: [alarmSchema],
          params: paramsWithItemId,
          dataPreparer: d => d.data[0].alarms,
        }, { root: true });
      } catch (err) {
        console.error(err);
      }
    },

  },
};
