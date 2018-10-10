import i18n from '@/i18n';
import request from '@/services/request';
import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';
import { groupSchema } from '@/store/schemas';

export const types = {
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
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
    [types.FETCH_LIST_COMPLETED](state, { allIds }) {
      state.allIds = allIds;
    },
  },
  actions: {
    async fetchList({ commit, dispatch }) {
      try {
        const { normalizedData } = await dispatch('entities/fetch', {
          route: API_ROUTES.view,
          schema: [groupSchema],
          dataPreparer: d => Object.keys(d.groups).map(key => ({ _id: key, ...d.groups[key] })),
        }, { root: true });
        commit(types.FETCH_LIST_COMPLETED, {
          allIds: normalizedData.result,
        });
      } catch (err) {
        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });
      }
    },

    create(context, { data } = {}) {
      return request.post(API_ROUTES.viewGroup, data);
    },

    update(context, { id, data } = {}) {
      return request.put(`${API_ROUTES.viewGroup}/${id}`, data);
    },

    remove(context, { id }) {
      return request.delete(`${API_ROUTES.viewGroup}/${id}`);
    },
  },
};
