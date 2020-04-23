export const types = {
  SET_PLAYLIST: 'SET_PLAYLIST',
  PLAY: 'PLAY',
  STOP: 'STOP',
  GO_TO_TAB: 'GO_TO_TAB',
};

export default {
  namespaced: true,
  state: {
    playlist: null,
    playing: false,
    activeTabIndex: 0,
  },
  getters: {
    playing: state => state.playing,
    activeTabIndex: state => state.activeTabIndex,
    playlist: state => state.playlist,
  },
  mutations: {
    [types.SET_PLAYLIST]: (state, { playlist }) => {
      state.playing = false;
      state.activeTabIndex = 0;
      state.playlist = playlist;
    },
    [types.PLAY]: (state) => {
      state.playing = true;
    },
    [types.STOP]: (state) => {
      state.playing = false;
    },
    [types.GO_TO_TAB]: (state, { index }) => {
      state.activeTabIndex = index;
    },
  },
  actions: {
    setPlaylist({ commit }, { playlist = null } = {}) {
      commit(types.SET_PLAYLIST, { playlist });
    },

    play({ commit }) {
      commit(types.PLAY);
    },

    stop({ commit }) {
      commit(types.STOP);
    },

    prevTab({ commit, state }) {
      const { tabs = [] } = state.playlist || {};

      if (tabs.length) {
        const index = state.activeTabIndex <= 0 ? 0 : state.activeTabIndex - 1;

        commit(types.GO_TO_TAB, { index });
      }
    },

    nextTab({ commit, state }) {
      const { tabs = [] } = state.playlist || {};

      if (tabs.length) {
        const lastIndex = tabs.length - 1;
        const index = state.activeTabIndex >= lastIndex ? lastIndex : state.activeTabIndex + 1;

        commit(types.GO_TO_TAB, { index });
      }
    },
  },
};
