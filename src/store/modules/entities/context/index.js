import { API_ROUTES } from '@/config';
import { contextSchema } from '@/store/schemas';
import request from '@/services/request';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
  EDIT_FAILED: 'EDIT_FAILED',
};

export default {
  namespaced: true,
  state: {
    allIds: [],
    meta: {},
    pending: false,
    fetchingParams: {},
    error: '',
  },
  getters: {
    allIds: state => state.allIds,
    items: (state, getters, rootState, rootGetters) => rootGetters['entities/getList']('context', state.allIds),
    meta: state => state.meta,
    pending: state => state.pending,
    error: state => state.error,
  },
  mutations: {
    [types.FETCH_LIST](state, { params }) {
      state.pending = true;
      state.fetchingParams = params;
    },
    [types.FETCH_LIST_COMPLETED](state, { allIds, meta }) {
      state.allIds = allIds;
      state.meta = meta;
      state.pending = false;
    },
    [types.FETCH_LIST_FAILED](state) {
      state.pending = false;
    },
    [types.EDIT_FAILED](state, err) {
      state.error = err;
    },
  },
  actions: {
    async fetchList({ commit, dispatch }, { params } = {}) {
      try {
        commit(types.FETCH_LIST, { params });

        const { normalizedData, data } = await dispatch('entities/fetch', {
          route: API_ROUTES.context,
          schema: [contextSchema],
          params,
          dataPreparer: d => d.data,
          isPost: true,
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
      }
    },

    async edit({ commit }, { data }) {
      try {
        await request.put(API_ROUTES.context, { entity: data, _type: 'crudcontext' });
      } catch (err) {
        commit(types.EDIT_FAILED, err);
      }
    },
  },
};
