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
    getWidget: state => ({ wrapperId = null, widgetXType = null }) => {
      if (!state.widgets) {
        return null;
      }

      if (!wrapperId && !widgetXType) {
        return null;
      }

      if (wrapperId && !widgetXType) {
        return state.widgets[wrapperId];
      }

      if (!wrapperId && widgetXType) {
        let widgetWrapper = null;

        Object.keys(state.widgets)
          .forEach((id) => {
            if (state.widgets[id].widget.xtype === widgetXType) {
              widgetWrapper = state.widgets[id];
            }
          });

        return widgetWrapper;
      }

      if (state.widgets[wrapperId].widget.xtype === widgetXType) {
        return state.widgets[wrapperId];
      }

      return null;
    },
    getWidgets: state => (asArray = false) => {
      if (asArray) {
        return Object.values(state.widgets);
      }

      return state.widgets;
    },
  },
  actions: {
    async save({ commit, dispatch }, widgetWrapper) {
      commit(types.SET_WIDGET, widgetWrapper);
      await dispatch('view/save', {}, { root: true });
    },
  },
};
