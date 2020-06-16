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
    getList: (state, getters, rootState, rootGetters) => ids =>
      rootGetters['entities/getList'](ENTITIES_TYPES.alarm, ids),
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

        await dispatch('popups/error', { text: i18n.t('errors.default') }, { root: true });

        return { alarms: [], total: 0 };
      }
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
