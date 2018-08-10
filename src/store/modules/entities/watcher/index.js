import request from '@/services/request';
import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES, WIDGET_TYPES } from '@/constants';
import { entitySchema } from '@/store/schemas';

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
    async create(context, { params = {} } = {}) {
      try {
        await request.post(
          API_ROUTES.watcher,
          { _id: params._id, mfilter: params.mfilter, display_name: params.display_name },
        );
      } catch (err) {
        console.warn(err);
      }
    },

    async edit({ dispatch }, { data }) {
      try {
        await request.put(API_ROUTES.context, { entity: data, _type: WIDGET_TYPES.context });
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
