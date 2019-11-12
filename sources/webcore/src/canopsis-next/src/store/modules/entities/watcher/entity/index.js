import Vue from 'vue';
import { get } from 'lodash';

import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';
import { watcherEntitySchema } from '@/store/schemas';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default {
  namespaced: true,
  state: {
    watchers: {},
  },
  getters: {
    getListByWatcherId: (state, getters, rootState, rootGetters) => watcherId =>
      rootGetters['entities/getList'](ENTITIES_TYPES.watcherEntity, get(state.watchers[watcherId], 'allIds', [])),
    getPendingByWatcherId: state => watcherId => get(state.watchers[watcherId], 'pending'),
  },
  mutations: {
    [types.FETCH_LIST](state, { watcherId }) {
      Vue.set(state.watchers, watcherId, {
        ...state.watchers[watcherId],
        pending: true,
      });
    },
    [types.FETCH_LIST_COMPLETED](state, { watcherId, allIds }) {
      Vue.set(state.watchers, watcherId, {
        ...state.watchers[watcherId],
        pending: false,
        allIds,
      });
    },
    [types.FETCH_LIST_FAILED](state, { watcherId }) {
      Vue.set(state.watchers, watcherId, {
        ...state.watchers[watcherId],
        pending: false,
      });
    },
  },
  actions: {
    async fetchList({ commit, dispatch }, { watcherId, params }) {
      try {
        commit(types.FETCH_LIST, { watcherId });

        const { normalizedData } = await dispatch('entities/fetch', {
          route: `${API_ROUTES.weatherWatcher}/${watcherId}`,
          schema: [watcherEntitySchema],
          params,
        }, { root: true });

        commit(types.FETCH_LIST_COMPLETED, {
          watcherId,
          allIds: normalizedData.result,
        });
      } catch (err) {
        console.error(err);

        commit(types.FETCH_LIST_FAILED, { watcherId });
      }
    },
  },
};
