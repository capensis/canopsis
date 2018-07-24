import request from '@/services/request';
import { viewSchema } from '@/store/schemas';
import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import widgetModule, { types as widgetTypes } from './widget';

export const types = {
  FETCH_ITEM: 'FETCH_ITEM',
  FETCH_ITEM_COMPLETED: 'FETCH_ITEM_COMPLETED',
  SET_LOADED_VIEW: 'SET_LOADED_VIEW',
};

export default {
  namespaced: true,
  modules: {
    widget: widgetModule,
  },
  state: {
    viewId: null,
    pending: false,
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
    [types.SET_LOADED_VIEW]: (state, view) => {
      state.view = view;
    },
  },
  actions: {
    /**
     * This action updates view id and updates widgets ids in the store
     *
     * @param {function} commit
     * @param {Object} normalizedData
     */
    fetchedItem({ commit }, { normalizedData }) {
      commit(types.FETCH_ITEM_COMPLETED, normalizedData.result);
      commit(
        `view/widget/${widgetTypes.UPDATE_WIDGETS_IDS}`,
        Object.keys(normalizedData.entities.widget),
        { root: true },
      );
    },

    /**
     * This action fetches view by id
     *
     * @param {function} commit
     * @param {function} dispatch
     * @param {string|number} id
     */
    async fetchItem({ commit, dispatch }, { id }) {
      try {
        commit(types.FETCH_ITEM);
        const result = await dispatch('entities/fetch', {
          route: `${API_ROUTES.view}/${id}`,
          schema: viewSchema,
          dataPreparer: d => d.data[0],
        }, { root: true });

        await dispatch('fetchedItem', result);
      } catch (e) {
        console.error(e);
      }
    },

    /**
     * This action updates view
     *
     * @param {function} dispatch
     * @param {Object} view
     */
    async update({ dispatch }, { view }) {
      try {
        const result = await dispatch('entities/update', {
          route: `${API_ROUTES.view}/${view.id}`,
          schema: viewSchema,
          body: view,
          dataPreparer: d => d.data[0],
        }, { root: true });

        await dispatch('fetchedItem', result);
      } catch (err) {
        console.warn(err);
      }
    },

    async saveItem({ commit, rootGetters, getters }) {
      const view = getters.item;

      view.containerwidget.items = rootGetters['view/widget/getItems'](true);

      try {
        await request.put(`${API_ROUTES.view}/${view.id}`, view);

        commit(types.SET_LOADED_VIEW, view);
      } catch (e) {
        console.error(e);
      }
    },
  },
};
