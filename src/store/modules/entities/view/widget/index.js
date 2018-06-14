export const types = {
  SET_WIDGETS: 'SET_WIDGETS',
  SET_WIDGET: 'SET_WIDGET',
};

export default {
  namespaced: true,
  state: {
    widgets: null,
  },
  mutations: {
    [types.SET_WIDGETS]: (state, widgets) => {
      state.widgets = widgets;
    },
    [types.SET_WIDGET]: (state, widgetWrapper) => {
      state.widgets[widgetWrapper.id] = widgetWrapper;
    },
  },
  getters: {
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
  },
};
