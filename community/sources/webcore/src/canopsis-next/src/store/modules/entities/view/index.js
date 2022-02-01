import { normalize } from 'normalizr';

import request, { useRequestCancelling } from '@/services/request';
import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';
import { viewSchema } from '@/store/schemas';
import { types as entitiesTypes } from '@/store/plugins/entities';

import groupModule from './group';

export const types = {
  FETCH_ITEM: 'FETCH_ITEM',
  FETCH_ITEM_COMPLETED: 'FETCH_ITEM_COMPLETED',
  FETCH_ITEM_FAILED: 'FETCH_ITEM_FAILED',
};

export default {
  namespaced: true,
  modules: {
    group: groupModule,
  },
  state: {
    activeViewId: null,
    pending: true,
  },
  getters: {
    itemId: state => state.activeViewId,
    pending: state => state.pending,
    item: (state, getters, rootState, rootGetters) => rootGetters['entities/getItem'](
      ENTITIES_TYPES.view,
      state.activeViewId,
    ),
    getItemById: (state, getters, rootState, rootGetters) => id => rootGetters['entities/getItem'](
      ENTITIES_TYPES.view,
      id,
    ),
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
    [types.FETCH_ITEM_FAILED]: (state) => {
      state.pending = false;
    },
  },
  actions: {
    async fetchItem({ commit, dispatch }, { id }) {
      try {
        await useRequestCancelling(async (source) => {
          commit(types.FETCH_ITEM, id);

          const { normalizedData } = await dispatch('entities/fetch', {
            route: `${API_ROUTES.view}/${id}`,
            schema: viewSchema,
            cancelToken: source.token,
          }, { root: true });

          commit(types.FETCH_ITEM_COMPLETED, normalizedData.result);
        }, 'activeView');
      } catch (err) {
        console.error(err);

        commit(types.FETCH_ITEM_FAILED);
      }
    },

    create(context, { data } = {}) {
      return request.post(API_ROUTES.view, data);
    },

    async update({ commit }, { id, data }) {
      const result = await request.put(`${API_ROUTES.view}/${id}`, data);

      const { entities } = normalize(result, viewSchema);

      commit(entitiesTypes.ENTITIES_UPDATE, entities, { root: true });

      return result;
    },

    updateWithoutStore(context, { id, data }) {
      return request.put(`${API_ROUTES.view}/${id}`, data);
    },

    updatePositions(context, { data }) {
      return request.put(API_ROUTES.viewPosition, data);
    },

    remove(context, { id }) {
      return request.delete(`${API_ROUTES.view}/${id}`);
    },

    bulkCreateWithoutStore(context, { data }) {
      return request.post(API_ROUTES.bulkView, data);
    },
  },
};
