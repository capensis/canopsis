import request from '@/services/request';
import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';
import i18n from '@/i18n';
import { groupSchema } from '@/store/schemas';

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
    items: (state, getters, rootState, rootGetters) => rootGetters['entities/getList'](ENTITIES_TYPES.group, state.allIds),
  },
  mutations: {
    [types.FETCH_LIST](state) {
      state.pending = true;
    },
    [types.FETCH_LIST_COMPLETED](state, { allIds }) {
      state.allIds = allIds;
      state.pending = false;
    },
    [types.FETCH_LIST_FAILED](state) {
      state.pending = false;
    },
  },
  actions: {
    create(context, params = {}) {
      return request.post(API_ROUTES.viewV3.groups, params);
    },
    async fetchList({ commit, dispatch }, { id } = {}) {
      try {
        let route = API_ROUTES.viewV3.groups;
        if (id) {
          route += `/${id}`;
        }
        const { normalizedData } = await dispatch('entities/fetch', {
          route,
          schema: [groupSchema],
          dataPreparer: d => d,
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
