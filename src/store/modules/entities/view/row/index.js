import { normalize } from 'normalizr';

import { ENTITIES_TYPES } from '@/constants';
import { viewRowSchema } from '@/store/schemas';
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
    create({ dispatch, rootGetters }, { row }) {
      const view = rootGetters['view/item'];

      view.rows.push(row);

      return dispatch('view/update', { view }, { root: true });
    },

    update({ dispatch, commit, rootGetters }, { row }) {
      const normalizedData = normalize(row, viewRowSchema);

      commit(entitiesTypes.ENTITIES_UPDATE, normalizedData.entities, { root: true });

      const view = rootGetters['view/item'];

      return dispatch('view/update', { view }, { root: true });
    },
  },
};
