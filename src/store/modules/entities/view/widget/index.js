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

      row.widgets.push(widget);

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
    async update({ dispatch, rootGetters }, { widget, rowId }) {
      const row = rootGetters['view/row/getItem'](rowId);

      row.widgets = row.widgets.map((v) => {
        if (v._id === widget._id) {
          return widget;
        }

        return v;
      });

      return dispatch('view/row/update', { row }, { root: true });
    },
  },
};
