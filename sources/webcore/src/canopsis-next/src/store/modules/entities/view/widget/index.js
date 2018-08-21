import { normalize } from 'normalizr';

import uuid from '@/helpers/uuid';
import { ENTITIES_TYPES, WIDGET_TYPES } from '@/constants';
import { widgetSchema } from '@/store/schemas';
import { types as entitiesTypes } from '@/store/plugins/entities';

export const types = {
  SET_WIDGET: 'SET_WIDGET',
  UPDATE_WIDGETS_IDS: 'UPDATE_WIDGETS_IDS',
};

export default {
  namespaced: true,
  state: {
    allIds: [],
  },
  mutations: {
    [types.UPDATE_WIDGETS_IDS]: (state, allIds) => {
      state.allIds = allIds;
    },
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
      rootGetters['entities/getList'](ENTITIES_TYPES.widget, state.allIds),
  },
  actions: {
    /**
     * This action does: create widget, put it into view and update view
     *
     * @param {function} dispatch
     * @param {Object} rootGetters
     * @param {Object} widget
     */
    async create({ dispatch, rootGetters }, { widget }) {
      const view = rootGetters['view/item'];
      const widgetWrapper = {
        id: uuid(WIDGET_TYPES.widgetWrapper),
        title: 'wrapper',
        xtype: WIDGET_TYPES.widgetWrapper,
        mixins: [],
        widget,
      };

      view.containerwidget.items.push(widgetWrapper);

      await dispatch('view/update', { view }, { root: true });
    },
    /**
     * This action does: update widget, put it into view and update view
     *
     * @param {function} commit
     * @param {function} dispatch
     * @param {Object} rootGetters
     * @param {Object} widget
     */
    async update({ commit, dispatch, rootGetters }, { widget }) {
      const normalizedData = normalize(widget, widgetSchema);

      commit(entitiesTypes.ENTITIES_UPDATE, normalizedData.entities, { root: true });

      const view = rootGetters['view/item'];

      await dispatch('view/update', { view }, { root: true });
    },
  },
};
