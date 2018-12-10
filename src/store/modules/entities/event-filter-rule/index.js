import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';
import { eventFilterRuleSchema } from '@/store/schemas';
import request from '@/services/request';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default {
  namespaced: true,
  state: {
    allIds: [],
    pending: false,
    fetchingParams: {},
  },
  getters: {
    items: (state, getters, rootState, rootGetters) => rootGetters['entities/getList'](ENTITIES_TYPES.eventFilterRule, state.allIds),
    pending: state => state.pending,
  },
  mutations: {
    [types.FETCH_LIST](state, { params }) {
      state.pending = true;
      state.fetchinParams = params;
    },
    [types.FETCH_LIST_COMPLETED](state, { allIds }) {
      state.pending = false;
      state.allIds = allIds;
    },
    [types.FETCH_LIST_FAILED](state) {
      state.pending = false;
    },
  },
  actions: {
    async fetchList({ commit, dispatch }, { params } = {}) {
      try {
        commit(types.FETCH_LIST, { params });

        const { normalizedData } = await dispatch('entities/fetch', {
          route: API_ROUTES.eventFilterRules,
          schema: [eventFilterRuleSchema],
          params,
        }, { root: true });

        commit(types.FETCH_LIST_COMPLETED, {
          allIds: normalizedData.result,
        });
      } catch (err) {
        console.error(err);
        commit(types.FETCH_LIST_FAILED);
      }
    },

    fetchListWithPreviousParams({ dispatch, state }) {
      return dispatch('fetchList', {
        params: state.fetchingParams,
      });
    },

    create(context, { data }) {
      return request.post(API_ROUTES.eventFilterRules, data);
    },

    edit(context, { id, data }) {
      return request.put(`${API_ROUTES.eventFilterRules}/${id}`, data);
    },

    async remove({ dispatch }, { id } = {}) {
      try {
        await request.delete(`${API_ROUTES.eventFilterRules}/${id}`);
        await dispatch('entities/removeFromStore', {
          id,
          type: ENTITIES_TYPES.eventFilterRule,
        }, { root: true });
      } catch (err) {
        console.error(err);
      }
    },
  },
};
