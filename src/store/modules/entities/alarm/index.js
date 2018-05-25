import { API_ROUTES } from '@/config';
import { alarmSchema } from '@/store/schemas';

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
  },
  getters: {
    allIds: state => state.allIds,
    items: (state, getters, rootState, rootGetters) => rootGetters['entities/getList']('alarm', state.allIds),
    meta: state => state.meta,
    pending: state => state.pending,
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
    [types.FETCH_ITEM](state) {
      state.item = {};
      state.itemPending = true;
    },
    [types.FETCH_ITEM_COMPLETED](state, { item }) {
      state.item = item;
      state.itemPending = false;
    },
  },
  actions: {
    async fetchList({ commit, dispatch }, { params } = {}) {
      try {
        commit(types.FETCH_LIST, { params });

        const { normalizedData, data } = await dispatch('entities/fetch', {
          route: API_ROUTES.alarmList,
          schema: [alarmSchema],
          params,
          dataPreparer: d => d.alarms,
        }, { root: true });

        commit(types.FETCH_LIST_COMPLETED, {
          allIds: normalizedData.result,
          meta: {
            first: data.first,
            last: data.last,
            total: data.total,
          },
        });
      } catch (err) {
        console.error(err);

        commit(types.FETCH_LIST_FAILED);
      }
    },

    fetchListWithPreviousParams({ dispatch, state }) {
      return dispatch('fetchList', { params: state.fetchingParams });
    },

    async fetchItem({ dispatch }, { id }) {
      try {
        await dispatch('entities/fetch', {
          route: API_ROUTES.alarmList,
          schema: [alarmSchema],
          params: { filter: { _id: id } },
          dataPreparer: d => d.alarms,
        }, { root: true });
      } catch (err) {
        console.error(err);
      }
    },

  },
};
