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
    getMetaByWatcherId: state => watcherId => get(state.watchers[watcherId], 'meta', {}),
  },
  mutations: {
    [types.FETCH_LIST](state, { watcherId }) {
      Vue.setSeveral(state.watchers, watcherId, { pending: true });
    },
    [types.FETCH_LIST_COMPLETED](state, { watcherId, meta, allIds }) {
      Vue.setSeveral(state.watchers, watcherId, { allIds, meta, pending: false });
    },
    [types.FETCH_LIST_FAILED](state, { watcherId, error = {} }) {
      Vue.setSeveral(state.watchers, watcherId, { error, pending: false });
    },
  },
  actions: {
    async fetchList({ commit, dispatch }, { watcherId, params = {} }) {
      try {
        commit(types.FETCH_LIST, { watcherId });

        const { normalizedData, data } = await dispatch('entities/fetch', {
          route: `${API_ROUTES.weatherWatcher}/${watcherId}`,
          schema: [watcherEntitySchema],
          dataPreparer: d => d.data,
          params,
        }, { root: true });

        commit(types.FETCH_LIST_COMPLETED, {
          watcherId,
          allIds: normalizedData.result,
          meta: data.meta,
        });
      } catch (err) {
        console.error(err);

        commit(types.FETCH_LIST_FAILED, { watcherId });
      }
    },
  },
};
