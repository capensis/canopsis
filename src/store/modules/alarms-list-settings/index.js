const types = {
  OPEN_SETTINGS_PANEL: 'OPEN_SETTINGS_PANEL',
  CLOSE_SETTINGS_PANEL: 'CLOSE_SETTINGS_PANEL',
  FETCH_SETTINGS: 'FETCH_SETTINGS',
  FETCH_SETTINGS_COMPLETED: 'FETCH_SETTINGS_COMPLETED',
};

export default {
  namespaced: true,

  state: {
    isPanelOpen: false,
    settings: {},
  },

  getters: {
    isPanelOpen: state => state.isPanelOpen,
  },

  mutations: {
    [types.OPEN_SETTINGS_PANEL](state) {
      state.isPanelOpen = true;
    },
    [types.CLOSE_SETTINGS_PANEL](state) {
      state.isPanelOpen = false;
    },
    [types.FETCH_SETTINGS](state) {
      state.settings = {};
    },
    [types.FETCH_SETTINGS_COMPLETED]() {
      // TODO : Store settings data
    },
  },

  actions: {
    closePanel(context) {
      context.commit(types.CLOSE_SETTINGS_PANEL);
    },
    openPanel(context) {
      context.commit(types.OPEN_SETTINGS_PANEL);
    },
    async fetchSettings(context) {
      context.commit(types.FETCH_SETTINGS);

      try {
        // TODO : REQUEST TO GET SETTINGS -> COMMIT FETCH_COMPLETED
      } catch (e) {
        console.warn(e);
      }
    },
  },
};
