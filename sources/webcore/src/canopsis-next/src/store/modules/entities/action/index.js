import request from '@/services/request';
import i18n from '@/i18n';
import schemas from '@/store/schemas';
import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
  FETCH_BY_ID: 'FETCH_BY_ID',
  FETCH_BY_ID_COMPLETED: 'FETCH_BY_ID_COMPLETED',
  FETCH_BY_ID_FAILED: 'FETCH_BY_ID_FAILED',
};

export default {
  namespaced: true,
  state: {
    pending: false,
    allIds: [],
  },
  getters: {
    allIds: state => state.allIds,
    items: (state, getters, rootState, rootGetters) =>
      rootGetters['entities/getList'](ENTITIES_TYPES.action, state.allIds),
    pending: state => state.pending,
  },
  mutations: {
    [types.FETCH_LIST](state) {
      state.pending = true;
    },
    [types.FETCH_LIST_COMPLETED](state, { allIds }) {
      state.allIds = allIds;
      state.pending = false;
    },
  },
  actions: {
    async fetchList({ dispatch, commit }, { params } = {}) {
      try {
        commit(types.FETCH_LIST);

        const { normalizedData } = await dispatch('entities/fetch', {
          route: API_ROUTES.actions,
          schema: [schemas.action],
          params,
        }, { root: true });

        commit(types.FETCH_LIST_COMPLETED, {
          allIds: normalizedData.result,
        });
      } catch (err) {
        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });
        commit(types.FETCH_LIST_FAILED);
      }
    },

    async remove({ dispatch }, { id }) {
      try {
        await request.delete(`${API_ROUTES.actions}/${id}`);
        await dispatch('entities/removeFromStore', { id, type: ENTITIES_TYPES.action }, { root: true });
      } catch (err) {
        console.warn(err);
      }
    },
  },
};
