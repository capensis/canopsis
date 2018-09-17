import request from '@/services/request';
import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';
import { viewSchema } from '@/store/schemas';

import groupModule from './group';

export const types = {
  FETCH_ITEM: 'FETCH_ITEM',
  FETCH_ITEM_COMPLETED: 'FETCH_ITEM_COMPLETED',
};

export default {
  namespaced: true,
  modules: {
    group: groupModule,
  },
  state: {
    viewId: null,
  },
  getters: {
    item: (state, getters, rootState, rootGetters) =>
      rootGetters['entities/getItem'](ENTITIES_TYPES.view, state.viewId),
  },
  mutations: {
    [types.FETCH_ITEM]: (state) => {
      state.pending = true;
    },
    [types.FETCH_ITEM_COMPLETED]: (state, viewId) => {
      state.viewId = viewId;
      state.pending = false;
    },
  },
  actions: {
    create(context, { data } = {}) {
      return request.post(API_ROUTES.view, data);
    },

    async fetchItem({ commit, dispatch }, { id }) {
      try {
        commit(types.FETCH_ITEM);

        const { normalizedData } = await dispatch('entities/fetch', {
          route: `${API_ROUTES.view}/${id}`,
          schema: viewSchema,
          dataPreparer: d => d,
        }, { root: true });

        commit(types.FETCH_ITEM_COMPLETED, normalizedData.result);
      } catch (err) {
        console.error(err);
      }
    },
  },
};
