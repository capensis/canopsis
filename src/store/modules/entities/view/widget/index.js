import { normalize } from 'normalizr';

import { ENTITIES_TYPES } from '@/constants';
import { widgetSchema } from '@/store/schemas';
import { types as entitiesTypes } from '@/store/plugins/entities';

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
      const view = rootGetters['view/item'];

      view.rows = view.rows.map((row) => {
        if (row._id === rowId) {
          return { ...row, widgets: [...row.widgets, widget] };
        }

        return row;
      });

      const row = rootGetters['view/row/getItem'](rowId);

      if (row) {
        row.widgets.push(widget);
      }

      return dispatch('view/row/update', { row }, { root: true });
    },

    /**
     * This action does: update widget, put it into view and update view
     *
     * @param {function} commit
     * @param {function} dispatch
     * @param {Object} getters
     * @param {Object} rootGetters
     * @param {Object} widget
     */
    async update({
      commit, dispatch, getters, rootGetters,
    }, { widget }) {
      const cachedWidget = getters.getItem(widget._id);

      try {
        const normalizedData = normalize(widget, widgetSchema);

        commit(entitiesTypes.ENTITIES_UPDATE, normalizedData.entities, { root: true });

        const view = rootGetters['view/item'];

        await dispatch('view/update', { view }, { root: true });
      } catch (err) {
        const normalizedData = normalize(cachedWidget, widgetSchema);

        commit(entitiesTypes.ENTITIES_UPDATE, normalizedData.entities, { root: true });
      }
    },
  },
};
