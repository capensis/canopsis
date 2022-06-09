import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES, REQUEST_METHODS } from '@/constants';

import request from '@/services/request';

import { snmpRuleSchema } from '@/store/schemas';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default {
  namespaced: true,
  state: {
    allIds: [],
    meta: {},
    pending: false,
  },
  getters: {
    items: (state, getters, rootState, rootGetters) => rootGetters['entities/getList'](
      ENTITIES_TYPES.snmpRule,
      state.allIds,
    ),

    meta: state => state.meta,
    pending: state => state.pending,
  },
  mutations: {
    [types.FETCH_LIST](state) {
      state.pending = true;
    },
    [types.FETCH_LIST_COMPLETED](state, { allIds, meta }) {
      state.pending = false;
      state.allIds = allIds;
      state.meta = meta;
    },
    [types.FETCH_LIST_FAILED](state) {
      state.pending = false;
    },
  },
  actions: {
    async fetchList({ commit, dispatch }, { params } = {}) {
      try {
        commit(types.FETCH_LIST, { params });

        const { data, normalizedData } = await dispatch('entities/fetch', {
          body: params,
          method: REQUEST_METHODS.post,
          route: API_ROUTES.snmpRule.list,
          dataPreparer: d => d.data,
          schema: [snmpRuleSchema],
        }, { root: true });

        commit(types.FETCH_LIST_COMPLETED, {
          allIds: normalizedData.result,
          meta: {
            total: data.total,
          },
        });
      } catch (err) {
        console.error(err);
        commit(types.FETCH_LIST_FAILED);

        throw err;
      }
    },

    create(context, { data } = {}) {
      return request.post(API_ROUTES.snmpRule.create, data);
    },

    remove(context, { data = {} } = {}) {
      return request.delete(API_ROUTES.snmpRule.list, { data });
    },
  },
};
