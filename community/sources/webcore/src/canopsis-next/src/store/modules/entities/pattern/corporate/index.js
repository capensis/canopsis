import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import { patternSchema } from '@/store/schemas';

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
    meta: {},
  },

  getters: {
    getItemById: (state, getters, rootState, rootGetters) => (
      id => rootGetters['entities/getItem'](ENTITIES_TYPES.pattern, id)
    ),

    items: (state, getters, rootState, rootGetters) => (
      rootGetters['entities/getList'](ENTITIES_TYPES.pattern, state.allIds)
    ),

    pending: state => state.pending,
    meta: state => state.meta,
  },

  mutations: {
    [types.FETCH_LIST](state, { params }) {
      state.pending = true;
      state.fetchingParams = params;
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
        const preparedParams = {
          ...params,

          corporate: true,
        };

        commit(types.FETCH_LIST, { params: preparedParams });

        const { normalizedData, data } = await dispatch('entities/fetch', {
          route: API_ROUTES.patterns,
          params: preparedParams,
          dataPreparer: d => d.data,
          schema: [patternSchema],
        }, { root: true });

        commit(types.FETCH_LIST_COMPLETED, {
          ...data,
          allIds: normalizedData.result,
        });
      } catch (err) {
        console.error(err);
        commit(types.FETCH_LIST_FAILED);

        throw err;
      }
    },

    fetchListWithPreviousParams({ dispatch, state }) {
      return dispatch('fetchList', {
        params: state.fetchingParams,
      });
    },
  },
};
