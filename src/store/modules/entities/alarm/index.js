import { API_ROUTES } from '@/config';
import { alarmSchema } from '@/store/schemas';

import fetchModule from '../fetch';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default {
  namespaced: true,
  modules: [fetchModule],
  state: {
    allIds: [],
    meta: {},
    pending: false,
    requestParams: {},
  },
  getters: {
    allIds: state => state.allIds,
    meta: state => state.meta,
    items: (state, getters, rootState, rootGetters) => rootGetters['entities/getList']('alarm', state.allIds),
  },
  mutations: {
    [types.FETCH_LIST](state, { params }) {
      state.pending = true;
      state.requestParams = params;
    },
    [types.FETCH_LIST_COMPLETED](state, { allIds, meta }) {
      state.allIds = allIds;
      state.meta = meta;
      state.pending = false;
    },
    [types.FETCH_LIST_FAILED](state) {
      state.pending = false;
    },
  },
  actions: {
    async fetchList({ commit, dispatch }, { params } = {}) {
      try {
        commit(types.FETCH_LIST, { params });

        const { normalizedData, data } = await dispatch('fetch', {
          route: API_ROUTES.alarmList,
          schema: [alarmSchema],
          params,
          dataPreparer: d => d.alarms,
        });

        commit(types.FETCH_LIST_COMPLETED, {
          allIds: normalizedData.result,
          meta: { total: data.total },
        });
      } catch (err) {
        console.error(err);

        commit(types.FETCH_LIST_FAILED);
      }
    },

    fetchListWithPreviousParams({ dispatch, state }) {
      return dispatch('fetchList', { params: state.requestParams });
    },

    async fetchItem({ dispatch }, { id }) {
      try {
        await dispatch('fetch', {
          route: API_ROUTES.alarmList,
          schema: [alarmSchema],
          params: { filter: { _id: id } },
          dataPreparer: d => d.alarms,
        });
      } catch (err) {
        console.error(err);
      }
    },
  },
};
