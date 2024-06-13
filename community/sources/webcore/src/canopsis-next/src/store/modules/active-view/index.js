import activeWidgets from './active-widgets';

export const types = {
  TOGGLE_EDITING: 'TOGGLE_EDITING',
  TOGGLE_EDITING_FINISHED: 'TOGGLE_EDITING_FINISHED',

  REGISTER_EDITING_HANDLER: 'REGISTER_EDITING_HANDLER',
  UNREGISTER_EDITING_HANDLER: 'UNREGISTER_EDITING_HANDLER',

  FETCH_ITEM: 'FETCH_ITEM',
  FETCH_ITEM_COMPLETED: 'FETCH_ITEM_COMPLETED',
  FETCH_ITEM_FAILED: 'FETCH_ITEM_FAILED',

  CLEAR: 'CLEAR',
};

export default {
  namespaced: true,
  modules: { activeWidgets },
  state: {
    id: null,
    pending: false,
    editing: false,
    editingProcess: false,
    editingOffHandlers: [],
  },
  getters: {
    editing: state => state.editing,
    editingProcess: state => state.editingProcess,
    pending: state => state.pending,
    item: (state, getters, rootState, rootGetters) => rootGetters['view/getViewById'](state.id),
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

    [types.FETCH_ITEM_COMPLETED]: (state) => {
      state.pending = false;
    },

    [types.FETCH_ITEM_FAILED]: (state) => {
      state.pending = false;
      state.id = null;
    },

    [types.CLEAR]: (state) => {
      state.id = null;
      state.pending = false;
      state.editing = false;
      state.editingOffHandlers = [];
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

        await dispatch('view/fetchView', { id }, { root: true });

        commit(types.FETCH_ITEM_COMPLETED);
      } catch (err) {
        commit(types.FETCH_ITEM_FAILED);

        console.error(err);

        throw err;
      }
    },

    clear({ commit, dispatch }) {
      commit(types.CLEAR);
      dispatch('activeWidgets/clear');
    },
  },
};
