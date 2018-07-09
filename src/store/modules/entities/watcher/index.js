import { watcherSchema } from '@/store/schemas';
import request from '@/services/request';
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
    allIds: [],
    meta: {},
    pending: false,
    fetchingParams: {},
    allIdsGeneralList: [],
    pendingGeneralList: false,
  },
  getters: {
    allIds: state => state.allIds,
    items: (state, getters, rootState, rootGetters) => rootGetters['entities/getList']('watcher', state.allIds),
    pending: state => state.pending,
    meta: state => state.meta,
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
  },
  actions: {
    fetch({ dispatch }, { params } = {}) {
      return dispatch('entities/fetch', {
        route: API_ROUTES.watcher,
        schema: [watcherSchema],
        params,
        dataPreparer: d => d.data,
        isPost: true,
      }, { root: true });
    },

    async fetchList({ commit, dispatch }, { params } = {}) {
      try {
        commit(types.FETCH_LIST, { params });

        const { normalizedData, data } = await dispatch('fetch', { params });

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

    async create({ dispatch }, params = {}) {
      try {
        console.log(params);
        await request.post(API_ROUTES.watcher, params);
      } catch (err) {
        console.warn(err);
      }
    },

    async remove({ dispatch }, { id } = {}) {
      try {
        await request.delete(API_ROUTES.watcher, { params: { watcher_id: id } });

        await dispatch('entities/removeFromStore', {
          id,
          type: ENTITIES_TYPES.watcher,
        }, { root: true });
      } catch (err) {
        console.warn(err);
      }
    },
  },
};
