import request from '@/services/request';
import { viewSchema } from '@/store/schemas';
import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';
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
    activeViewId: null,
    loadedView: null, // TODO:
    pending: false,
  },
  getters: {
    activeItem: (state, getters, rootState, rootGetters) =>
      rootGetters['entities/getItem'](ENTITIES_TYPES.view, state.activeViewId),
    pending: state => state.pending,
  },
  mutations: {
    [types.FETCH_ITEM]: (state) => {
      state.pending = true;
    },
    [types.FETCH_ITEM_COMPLETED]: (state, viewId) => {
      state.activeViewId = viewId;
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
      const view = getters.getItem;

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
