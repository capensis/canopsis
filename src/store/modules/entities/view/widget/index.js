import get from 'lodash/get';
import omit from 'lodash/omit';

import { ENTITIES_TYPES } from '@/constants';

export default {
  namespaced: true,
  getters: {
    getItem: (state, getters, rootState, rootGetters) => id =>
      rootGetters['entities/getItem'](ENTITIES_TYPES.widget, id),
  },
  actions: {
    /**
     * This action does: create widget, put it into row and update row
     *
     * @param {function} dispatch
     * @param {Object} rootGetters
     * @param {Object} widget
     * @param {string} rowId
     */
    create({ dispatch, rootGetters }, { widget, rowId }) {
      const row = rootGetters['view/row/getItem'](rowId);

      row.widgets.push(omit(widget, ['_embedded']));

      return dispatch('view/row/update', { row }, { root: true });
    },

    /**
     * This action does: update widget, put it into view and update view
     *
     * @param {function} dispatch
     * @param {Object} rootGetters
     * @param {Object} widget
     * @param {string} rowId
     */
    async update({ rootGetters }, { widget, rowId }) {
      const oldRowId = get(widget, '_embedded.parentId');
      const view = rootGetters['view/item'];

      view.rows = view.rows.map((row) => {
        if (row === oldRowId && rowId !== oldRowId) {
          return { ...row, widgets: row.widgets.filter(({ _id }) => _id !== widget._id) };
        } else if (row === rowId) {
          return { ...row, widgets: [...row.widgets, omit(widget, ['_embedded'])] };
        }

        return row;
      });

      // return dispatch('view/row/update', { row }, { root: true });
    },
  },
};
