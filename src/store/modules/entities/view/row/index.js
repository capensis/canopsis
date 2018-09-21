import { ENTITIES_TYPES } from '@/constants';
import { types as entitiesTypes } from '@/store/plugins/entities';

export const types = {
  UPDATE_ROWS_IDS: 'UPDATE_ROWS_IDS',
};

export default {
  namespaced: true,
  state: {
    activeRowsIds: [],
  },
  getters: {
    /**
     * Items of the active view
     *
     * @param {Object} state
     * @param {Object} getters
     * @param {Object} rootState
     * @param {Object} rootGetters
     */
    items: (state, getters, rootState, rootGetters) =>
      rootGetters['entities/getList'](ENTITIES_TYPES.viewRow, state.activeRowsIds),
  },
  mutations: {
    [types.UPDATE_ROWS_IDS]: (state, ids) => {
      state.activeRowsIds = ids;
    },
  },
  actions: {
    createInStore({ commit, rootGetters }, { row }) {
      const viewId = rootGetters['view/itemId'];

      commit(entitiesTypes.ENTITIES_MERGE, { viewRow: { [row._id]: row } }, { root: true });
      commit(entitiesTypes.ENTITIES_MERGE_CHILDREN, {
        path: `view.${viewId}.rows`,
        value: [row._id],
      }, { root: true });
    },

    create({ dispatch, rootGetters }, { row }) {
      const view = rootGetters['view/item'];

      view.rows.push(row);

      return dispatch('view/update', { view }, { root: true });
    },

    update({ dispatch, rootGetters }, { row }) {
      const view = rootGetters['view/item'];

      view.rows = view.rows.map((v) => {
        if (v._id === row._id) {
          return row;
        }

        return v;
      });

      return dispatch('view/update', { view }, { root: true });
    },
  },
};
