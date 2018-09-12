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
  },
  getters: {
    allIds: state => state.allIds,
    items: (state, getters, rootState, rootGetters) => rootGetters['entities/getList'](ENTITIES_TYPES.group, state.allIds),
  },
  mutations: {
    [types.FETCH_LIST]() {
    },
    [types.FETCH_LIST_COMPLETED](state, { allIds }) {
      state.allIds = allIds;
    },
    [types.FETCH_LIST_FAILED]() {
    },
  },
  actions: {
    create(context, { data } = {}) {
      return request.post(API_ROUTES.viewV3Group, data);
    },
    async fetchList({ commit, dispatch }) {
      try {
        const { normalizedData } = await dispatch('entities/fetch', {
          route: API_ROUTES.viewV3,
          schema: [groupSchema],
          dataPreparer: d => d.groups,
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
