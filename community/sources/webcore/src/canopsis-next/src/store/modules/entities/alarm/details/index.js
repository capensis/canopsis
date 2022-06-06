import Vue from 'vue';

import { API_ROUTES } from '@/config';

import { alarmDetailsSchema } from '@/store/schemas';
import { ENTITIES_TYPES } from '@/constants';

export const types = {
  FETCH_ITEM: 'FETCH_ITEM',
  FETCH_ITEM_COMPLETED: 'FETCH_ITEM_COMPLETED',
};

export default {
  namespaced: true,
  state: {
    pendingsById: {},
    widgets: {},
  },
  getters: {
    getItem: (state, getters, rootState, rootGetters) => (widgetId, id) => rootGetters['entities/getItem'](
      ENTITIES_TYPES.alarmDetails,
      id,
    ),
    getPending: state => id => state.pendingsById[id] ?? false,
  },
  mutations: {
    [types.FETCH_ITEM]: (state, { id }) => {
      Vue.set(state.pendingsById, id, true);
    },

    [types.FETCH_ITEM_COMPLETED]: (state, { id }) => {
      Vue.set(state.pendingsById, id, false);
    },
  },
  actions: {
    async fetchItem({ dispatch, commit }, { id, query } = {}) {
      try {
        commit(types.FETCH_ITEM, { id });

        await dispatch('entities/create', {
          route: API_ROUTES.alarmDetails,
          schema: [alarmDetailsSchema],
          body: [query],
        }, { root: true });
      } catch (err) {
        console.error(err);
      } finally {
        commit(types.FETCH_ITEM_COMPLETED, { id });
      }
    },
  },
};
