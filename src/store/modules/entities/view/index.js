// SERVICES
import request from '@/services/request';
// OTHERS
import { viewSchema } from '@/store/schemas';
import { API_ROUTES } from '@/config';
import widgetModule, { types as widgetMutations } from '@/store/modules/entities/view/widget';

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
    loadedView: null,
    pending: false,
  },
  mutations: {
    [types.FETCH_ITEM]: (state) => {
      state.pending = true;
    },
    [types.FETCH_ITEM_COMPLETED]: (state, view) => {
      state.loadedView = view;
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

        commit(types.FETCH_ITEM_COMPLETED, result.data[0]);
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
  getters: {
    getItem: state => state.loadedView,
    pending: state => state.pending,
  },
};
