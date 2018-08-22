import Vue from 'vue';
import merge from 'lodash/merge';
import get from 'lodash/get';

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
    async fetchListWithoutStore({ dispatch }, { params }) {
      try {
        const { data: [result] } = await request.get(API_ROUTES.alarmList, { params });

        return result;
      } catch (err) {
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
        const paramsWithItemId = merge(params, { filter: { d: id } });

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
