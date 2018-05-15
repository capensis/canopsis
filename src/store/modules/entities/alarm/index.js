import { normalize } from 'normalizr';

import { API_ROUTES } from '@/config';
import { alarmSchema } from '@/store/schemas';
import request from '@/services/request';

import entitiesTypes from '../types';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
};

export default {
  namespaced: true,
  state: {
    allIds: [],
    meta: {},
    pending: false,
  },
  getters: {
    allIds: state => state.allIds,
    meta: state => state.meta,
    items: (state, getters, rootState, rootGetters) => rootGetters['entities/getList']('alarm', state.allIds),
  },
  mutations: {
    [types.FETCH_LIST](state) {
      state.pending = true;
    },
    [types.FETCH_LIST_COMPLETED](state, { allIds, meta }) {
      state.allIds = allIds;
      state.meta = meta;
      state.pending = false;
    },
  },
  actions: {
    async fetchList({ commit }, params = {}) {
      commit(types.FETCH_LIST);

      try {
        const [data] = await request.get(API_ROUTES.alarmList, params);
        const normalizedData = normalize(data.alarms, [alarmSchema]);

        commit(`entities/${entitiesTypes.ENTITIES_UPDATE}`, normalizedData.entities, { root: true });
        commit(types.FETCH_LIST_COMPLETED, {
          allIds: normalizedData.result,
          meta: { total: data.total },
        });
      } catch (err) {
        console.error(err);
      }
    },
  },
};
