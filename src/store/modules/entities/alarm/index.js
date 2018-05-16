import { normalize } from 'normalizr';

import { API_ROUTES } from '@/config';
import { alarmSchema } from '@/store/schemas';
import request from '@/services/request';

import entitiesTypes from '../types';
import filtersModule from './filters';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
};

export default {
  namespaced: true,
  modules: {
    filters: filtersModule,
  },
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
  },
  actions: {
    async fetch({ commit }, { params }) {
      try {
        const [data] = await request.get(API_ROUTES.alarmList, { params });
        const normalizedData = normalize(data.alarms, [alarmSchema]);

        commit(`entities/${entitiesTypes.ENTITIES_UPDATE}`, normalizedData.entities, { root: true });

        return { data, normalizedData };
      } catch (err) {
        console.error(err);
      }

      return { data: {}, normalizedData: { result: [], entities: {} } };
    },

    async fetchList({ commit, dispatch }, { params } = {}) {
      commit(types.FETCH_LIST, { params });

      const { normalizedData, data } = await dispatch('fetch', { params });

      commit(types.FETCH_LIST_COMPLETED, {
        allIds: normalizedData.result,
        meta: { total: data.total },
      });
    },

    fetchListWithPreviousParams({ dispatch, state }) {
      return dispatch('fetchList', { params: state.requestParams });
    },

    fetchItem({ dispatch }, { id }) {
      return dispatch('fetch', { params: { filter: { _id: id } } });
    },
  },
};
