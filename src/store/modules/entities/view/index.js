import { normalize } from 'normalizr';

import request from '@/services/request';
import { viewSchema, widgetSchema } from '@/store/schemas';
import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';
import uuid from '@/helpers/uuid';
import { types as entitiesTypes } from '@/store/plugins/entities';

import widgetModule, { types as widgetMutations } from './widget';

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
    pending: state => state.pending,
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
    async fetchItem({ commit, dispatch }, { id }) {
      try {
        commit(types.FETCH_ITEM);
        const result = await dispatch('entities/fetch', {
          route: `${API_ROUTES.view}/${id}`,
          schema: viewSchema,
          dataPreparer: d => d.data[0],
        }, { root: true });

        commit(types.FETCH_ITEM_COMPLETED, result.normalizedData.result);
        commit(
          `view/widget/${widgetMutations.SET_WIDGETS}`,
          result.normalizedData.entities.widgetWrapper,
          { root: true },
        );
      } catch (e) {
        console.error(e);
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
    async createWidget({ commit, dispatch, getters }, { widget }) {
      const view = getters.item;

      view.items.push({
        id: uuid('widgetwrapper'),
        title: 'wrapper',
        xtype: 'widgetwrapper',
        mixins: [],
        widget,
      });

      commit(types.FETCH_ITEM);

      const result = await dispatch('entities/update', {
        route: `${API_ROUTES.view}/${view.id}`,
        schema: viewSchema,
        body: view,
        dataPreparer: d => d.data[0],
      }, { root: true });

      commit(types.FETCH_ITEM_COMPLETED, result.normalizedData.result);
    },
    async updateWidget({ commit, dispatch, getters }, { widget }) {
      const normalizedData = normalize(widget, widgetSchema);

      commit(entitiesTypes.ENTITIES_UPDATE, normalizedData.entities, { root: true });

      const view = getters.activeItem;

      commit(types.FETCH_ITEM);

      const result = await dispatch('entities/update', {
        route: `${API_ROUTES.view}/${view.id}`,
        schema: viewSchema,
        body: view,
        dataPreparer: d => d.data[0],
      }, { root: true });

      commit(types.FETCH_ITEM_COMPLETED, result.normalizedData.result);
    },
  },
};
