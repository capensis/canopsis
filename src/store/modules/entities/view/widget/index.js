export const types = {
  SET_WIDGETS: 'SET_WIDGETS',
  SET_WIDGET: 'SET_WIDGET',
  UPDATE_QUERY: 'UPDATE_QUERY',
};

export default {
  namespaced: true,
  state: {
    widgets: null,
    queries: {},
  },
  mutations: {
    [types.UPDATE_QUERY]: (state, { id, query }) => {
      state.queries[id] = query;
    },
    [types.SET_WIDGETS]: (state, widgets) => {
      state.widgets = widgets;
    },
    [types.SET_WIDGET]: (state, widgetWrapper) => {
      state.widgets[widgetWrapper.id] = widgetWrapper;
    },
  },
  getters: {
    getQuery: state => id => state.queries[id],
    getItem: state => ({ widgetXType }) => {
      if (!state.widgets) {
        return null;
      }

      let widgetWrapper = null;

      Object.keys(state.widgets)
        .forEach((id) => {
          if (state.widgets[id].widget.xtype === widgetXType) {
            widgetWrapper = state.widgets[id];
          }
        });

      return widgetWrapper;
    },
    getItems: state => (asArray = false) => {
      if (asArray) {
        return Object.values(state.widgets);
      }

      return state.widgets;
    },
  },
  actions: {
    async saveItem({ commit, dispatch }, { widgetWrapper }) {
      commit(types.SET_WIDGET, widgetWrapper);
      await dispatch('view/saveItem', {}, { root: true });
    },
    updateQuery({ commit }, { id, query }) {
      commit(types.UPDATE_QUERY, { id, query });
    },
  },
};
