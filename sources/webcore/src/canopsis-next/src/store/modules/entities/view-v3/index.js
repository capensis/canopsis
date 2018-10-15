import request from '@/services/request';
import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';
import { viewV3Schema } from '@/store/schemas';
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
  },
  getters: {
    item: (state, getters, rootState, rootGetters) =>
      rootGetters['entities/getItem'](ENTITIES_TYPES.viewV3, state.viewId),
  },
  mutations: {
    [types.FETCH_ITEM]: () => {
    },
    [types.FETCH_ITEM_COMPLETED]: () => {
    },
  },
  actions: {
    create(context, { data } = {}) {
      return request.post(API_ROUTES.viewV3, data);
    },

    async fetchItem({ commit, dispatch }, { id }) {
      try {
        commit(types.FETCH_ITEM);

        const normalizedData = await dispatch('entities/fetch', {
          route: `${API_ROUTES.viewV3.view}/${id}`,
          schema: viewV3Schema,
          dataPreparer: d => d,
        }, { root: true });

        commit(types.FETCH_ITEM_COMPLETED, normalizedData);
      } catch (err) {
        console.error(err);
      }
    },
  },
};
