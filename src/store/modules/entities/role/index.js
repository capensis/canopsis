import i18n from '@/i18n';

import { roleSchema } from '@/store/schemas';
import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';
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
  },
  getters: {
    items: (state, getters, rootState, rootGetters) => rootGetters['entities/getList'](ENTITIES_TYPES.role, state.allIds),
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
        commit(types.FETCH_LIST);
        const { normalizedData, data } = await dispatch('entities/fetch', {
          route: API_ROUTES.role.list,
          schema: [roleSchema],
          params,
          dataPreparer: d => d.data,
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

        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });
      }
    },

    remove(context, { id }) {
      return request.delete(`${API_ROUTES.role.remove}/${id}`);
    },
  },
};
