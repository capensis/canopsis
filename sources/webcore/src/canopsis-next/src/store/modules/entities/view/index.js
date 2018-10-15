import { normalize } from 'normalizr';

import request from '@/services/request';
import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';
import { viewSchema } from '@/store/schemas';
import { types as entitiesTypes } from '@/store/plugins/entities';

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
    activeViewId: null,
  },
  getters: {
    item: (state, getters, rootState, rootGetters) =>
      rootGetters['entities/getItem'](ENTITIES_TYPES.view, state.activeViewId),
  },
  mutations: {
    [types.FETCH_ITEM]: (state, viewId) => {
      state.pending = true;
      state.activeViewId = viewId;
    },
    [types.FETCH_ITEM_COMPLETED]: (state, viewId) => {
      state.activeViewId = viewId;
      state.pending = false;
    },
  },
  actions: {
    create(context, { data } = {}) {
      return request.post(API_ROUTES.view, data);
    },

    async fetchItem({ commit, dispatch }, { id }) {
      try {
        commit(types.FETCH_ITEM, id);

        const { normalizedData } = await dispatch('entities/fetch', {
          route: `${API_ROUTES.view}/${id}`,
          schema: viewSchema,
        }, { root: true });

        commit(types.FETCH_ITEM_COMPLETED, normalizedData.result);
      } catch (err) {
        console.error(err);
      }
    },

    async update({ commit }, { view }) {
      try {
        await request.put(`${API_ROUTES.view}/${view._id}`, view);

        const { entities } = normalize(view, viewSchema);

        commit(entitiesTypes.ENTITIES_UPDATE, entities, { root: true });
      } catch (err) {
        console.warn(err);
      }
    },
  },
};
