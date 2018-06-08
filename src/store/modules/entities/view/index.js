import axios from 'axios';
import { viewSchema } from '@/store/schemas';
import widgetModule, { types as widgetMutations } from './widget';

export const types = {
  SET_LOADED_VIEW: 'SET_LOADED_VIEW',
};

export default {
  namespaced: true,
  modules: {
    widget: widgetModule,
  },
  state: {
    loadedView: null,
  },
  mutations: {
    [types.SET_LOADED_VIEW]: (state, view) => {
      state.loadedView = view;
    },
    [types.SET_WIDGETS]: (state, widgets) => {
      state.loadedView.containerwidget.items = widgets;
    },
  },
  actions: {
    async loadView({ commit, dispatch }, viewId) {
      try {
        await axios.get('http://localhost:28082/account/me');
        await axios.get('http://localhost:28082/rest/default_rights/role/admin');
        await axios.get('http://localhost:28082/sessionstart?username=root');

        const result = await dispatch('entities/fetch', {
          route: `http://localhost:28082/rest/object/view/${viewId}`,
          schema: viewSchema,
          dataPreparer: d => d[0],
        }, { root: true });

        commit(types.SET_LOADED_VIEW, result.data[0]);
        commit(
          `view/widget/${widgetMutations.SET_WIDGETS}`,
          result.normalizedData.entities.widgetWrapper,
          { root: true },
        );
      } catch (e) {
        console.error(e);
      }
    },
    async save({ commit, rootGetters, getters }) {
      const view = getters.loadedView;

      view.containerwidget.items = rootGetters['view/widget/getWidgets'](true);

      try {
        await axios.put(`http://localhost:28082/rest/object/view/${view.id}`, view);

        commit(types.SET_LOADED_VIEW, view);
      } catch (e) {
        console.error(e);
      }
    },
  },
  getters: {
    loadedView: state => state.loadedView,
  },
};
