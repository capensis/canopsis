import { normalize } from 'normalizr';

import uuid from '@/helpers/uuid';
import { ENTITIES_TYPES } from '@/constants';
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
    items: (state, getters, rootState, rootGetters) =>
      rootGetters['entities/getList'](ENTITIES_TYPES.widget, state.allIds),
    getItem: state => ({ widgetXType }) => {
      if (!state.widgets) {
        return null;
      }

      let widgetWrapper = null;

      Object.keys(state.widgets)
        .forEach((id) => {
          if (state.widgets[id].widget.xtype === widgetXType) {
            widgetWrapper = state.widgets[id];
          }
        });

      return widgetWrapper;
    },
    getItems: state => (asArray = false) => {
      if (asArray) {
        return Object.values(state.widgets);
      }

      return state.widgets;
    },
  },
  actions: {
    async saveItem({ commit, dispatch }, { widgetWrapper }) {
      commit(types.SET_WIDGET, widgetWrapper);
      await dispatch('view/saveItem', {}, { root: true });
    },
    async create({ dispatch, rootGetters }, { widget }) {
      const view = rootGetters['view/item'];
      const widgetWrapper = {
        id: uuid('widgetwrapper'),
        title: 'wrapper',
        xtype: 'widgetwrapper',
        mixins: [],
        widget,
      };

      view.containerwidget.items.push(widgetWrapper);

      await dispatch('view/update', { view }, { root: true });
    },
    async update({ commit, dispatch, rootGetters }, { widget }) {
      const normalizedData = normalize(widget, widgetSchema);

      commit(entitiesTypes.ENTITIES_UPDATE, normalizedData.entities, { root: true });

      const view = rootGetters['view/item'];

      await dispatch('view/update', { view }, { root: true });
    },
  },
};
