import i18n from '@/i18n';
import request from '@/services/request';
import { actionSchema } from '@/store/schemas';
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
    pending: false,
    allIds: [],
  },
  getters: {
    pending: state => state.pending,
    items: (state, getters, rootState, rootGetters) =>
      rootGetters['entities/getList'](ENTITIES_TYPES.action, state.allIds),
  },
  mutations: {
    [types.FETCH_LIST](state) {
      state.pending = true;
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
    async fetchListWithoutStore({ dispatch }, { params }) {
      try {
        return await request.get(API_ROUTES.action, { params });
      } catch (err) {
        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });

        return { data: [], total: 0 };
      }
    },

    async fetchList({ commit, dispatch }, { params } = {}) {
      try {
        commit(types.FETCH_LIST);

        const { normalizedData } = await dispatch('entities/fetch', {
          route: API_ROUTES.action,
          schema: [actionSchema],
          params,
          dataPreparer: d => d.data,
        }, { root: true });

        commit(types.FETCH_LIST_COMPLETED, {
          allIds: normalizedData.result,
        });
      } catch (err) {
        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });

        commit(types.FETCH_LIST_FAILED);
      }
    },
  },
};
