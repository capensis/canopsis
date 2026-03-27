import { VIEW_SCREEN_MODES } from '@/constants';

import activeWidgetsModule from './active-widgets';

export const types = {
  TOGGLE_EDITING: 'TOGGLE_EDITING',
  TOGGLE_EDITING_FINISHED: 'TOGGLE_EDITING_FINISHED',

  REGISTER_EDITING_HANDLER: 'REGISTER_EDITING_HANDLER',
  UNREGISTER_EDITING_HANDLER: 'UNREGISTER_EDITING_HANDLER',

  FETCH_ITEM: 'FETCH_ITEM',
  FETCH_ITEM_COMPLETED: 'FETCH_ITEM_COMPLETED',
  FETCH_ITEM_FAILED: 'FETCH_ITEM_FAILED',

  RESUME_PERIODIC_REFRESH: 'RESUME_PERIODIC_REFRESH',
  PAUSE_PERIODIC_REFRESH: 'PAUSE_PERIODIC_REFRESH',

  CLEAR: 'CLEAR',

  SET_SCREEN_MODE: 'SET_SCREEN_MODE',
};

export default {
  namespaced: true,
  modules: { activeWidgets: activeWidgetsModule },
  state: {
    id: null,
    item: null,
    pending: false,
    editing: false,
    editingProcess: false,
    screenMode: VIEW_SCREEN_MODES.default,
    editingOffHandlers: [],
    periodicRefreshPaused: false,
  },
  getters: {
    editing: state => state.editing,
    editingProcess: state => state.editingProcess,
    pending: state => state.pending,
    screenMode: state => state.screenMode,
    isKioskScreenMode: state => [VIEW_SCREEN_MODES.kiosk, VIEW_SCREEN_MODES.kioskFullscreen].includes(state.screenMode),
    item: state => state.item,
    periodicRefreshPaused: state => state.periodicRefreshPaused,
  },
  mutations: {
    [types.TOGGLE_EDITING]: (state) => {
      state.editingProcess = true;
    },

    [types.TOGGLE_EDITING_FINISHED]: (state) => {
      state.editing = !state.editing;
      state.editingProcess = false;
    },

    [types.REGISTER_EDITING_HANDLER]: (state, handler) => {
      state.editingOffHandlers.push(handler);
    },

    [types.UNREGISTER_EDITING_HANDLER]: (state, handler) => {
      state.editingOffHandlers = state.editingOffHandlers.filter(editingHandler => editingHandler !== handler);
    },

    [types.FETCH_ITEM]: (state, { id }) => {
      state.pending = true;
      state.id = id;
    },

    [types.FETCH_ITEM_COMPLETED]: (state, { item }) => {
      state.item = item;
      state.pending = false;
    },

    [types.FETCH_ITEM_FAILED]: (state) => {
      state.pending = false;
      state.id = null;
    },

    [types.RESUME_PERIODIC_REFRESH]: (state) => {
      state.periodicRefreshPaused = false;
    },

    [types.PAUSE_PERIODIC_REFRESH]: (state) => {
      state.periodicRefreshPaused = true;
    },

    [types.CLEAR]: (state) => {
      state.id = null;
      state.pending = false;
      state.editing = false;
      state.screenMode = VIEW_SCREEN_MODES.default;
      state.editingOffHandlers = [];
      state.periodicRefreshPaused = false;
    },

    [types.SET_SCREEN_MODE]: (state, screenMode) => {
      state.screenMode = screenMode;
    },
  },
  actions: {
    registerEditingOffHandler({ commit }, handler) {
      commit(types.REGISTER_EDITING_HANDLER, handler);
    },

    unregisterEditingOffHandler({ commit }, handler) {
      commit(types.UNREGISTER_EDITING_HANDLER, handler);
    },

    async toggleEditing({ commit, state }) {
      try {
        commit(types.TOGGLE_EDITING);

        if (state.editingOffHandlers.length && state.editing) {
          await Promise.all(state.editingOffHandlers.map(handler => handler()));
        }

        commit(types.TOGGLE_EDITING_FINISHED);
      } catch (err) {
        console.error(err);
      }
    },

    async fetch({ state, dispatch, commit }, { id = state.id } = {}) {
      try {
        if (!id) {
          throw new Error('Active view id is empty');
        }

        commit(types.FETCH_ITEM, { id });

        const item = await dispatch('view/fetchView', { id }, { root: true });

        commit(types.FETCH_ITEM_COMPLETED, { item: item ?? state.item });
      } catch (err) {
        commit(types.FETCH_ITEM_FAILED);

        console.error(err);

        throw err;
      }
    },

    resumePeriodicRefresh({ commit }) {
      commit(types.RESUME_PERIODIC_REFRESH);
    },

    pausePeriodicRefresh({ commit }) {
      commit(types.PAUSE_PERIODIC_REFRESH);
    },

    clear({ commit, dispatch }) {
      commit(types.CLEAR);
      dispatch('activeWidgets/clear');
    },

    setScreenMode({ commit, getters }, screenMode) {
      if (getters.screenMode === screenMode) {
        return;
      }

      commit(types.SET_SCREEN_MODE, screenMode);
    },
  },
};
