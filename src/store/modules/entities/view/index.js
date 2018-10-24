import { normalize } from 'normalizr';
import i18n from '@/i18n';

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
    items: (state, getters, rootState) => rootState.entities.view,
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

    create(context, { data } = {}) {
      return request.post(API_ROUTES.view, data);
    },

    async update({ commit, dispatch }, { id, data }) {
      try {
        await request.put(`${API_ROUTES.view}/${id}`, data);

        const { entities } = normalize(data, viewSchema);

        commit(entitiesTypes.ENTITIES_UPDATE, entities, { root: true });
      } catch (err) {
        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });
        console.warn(err);
      }
    },

    remove(context, { id }) {
      return request.delete(`${API_ROUTES.view}/${id}`);
    },
  },
};
