import Vue from 'vue';
import { get } from 'lodash';

import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';
import { weatherServiceSchema } from '@/store/schemas';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default {
  namespaced: true,
  state: {
    services: {},
  },
  getters: {
    getListByServiceId: (state, getters, rootState, rootGetters) => serviceId => rootGetters['entities/getList'](
      ENTITIES_TYPES.weatherService,
      get(state.services[serviceId], 'allIds', []),
    ),
    getPendingByServiceId: state => serviceId => get(state.services[serviceId], 'pending'),
    getMetaByServiceId: state => serviceId => get(state.services[serviceId], 'meta', {}),
  },
  mutations: {
    [types.FETCH_LIST](state, { id }) {
      Vue.setSeveral(state.services, id, { pending: true });
    },
    [types.FETCH_LIST_COMPLETED](state, { id, meta, allIds }) {
      Vue.setSeveral(state.services, id, { allIds, meta, pending: false });
    },
    [types.FETCH_LIST_FAILED](state, { id, error = {} }) {
      Vue.setSeveral(state.services, id, { error, pending: false });
    },
  },
  actions: {
    async fetchList({ commit, dispatch }, { id, params = {} }) {
      try {
        commit(types.FETCH_LIST, { id });

        const { normalizedData, data } = await dispatch('entities/fetch', {
          route: `${API_ROUTES.weatherService}/${id}`,
          schema: [weatherServiceSchema],
          dataPreparer: d => d.data,
          params,
        }, { root: true });

        commit(types.FETCH_LIST_COMPLETED, {
          id,
          allIds: normalizedData.result,
          meta: data.meta,
        });
      } catch (err) {
        console.error(err);

        commit(types.FETCH_LIST_FAILED, { id });
      }
    },
  },
};
