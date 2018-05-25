import { normalize } from 'normalizr';
import { isEmpty } from 'lodash';

import { API_ROUTES, API_HOST } from '@/config';
import { eventSchema } from '@/store/schemas';

import axios from 'axios';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_ERROR: 'FETCH_ERROR',
};

export default {
  namespaced: true,

  state: {
    data: {},
    byId: {},
    allIds: {},
    fetchComplete: false,
    fetchError: '',
    meta: {},
  },

  getters: {
    data: state => state.data,
    byId: state => state.byId,
    allIds: state => state.allIds,
    meta: state => state.meta,
    fetchComplete: state => state.fetchComplete,
    fetchError: state => state.fetchError,
  },

  mutations: {
    [types.FETCH_LIST](state) {
      state.byId = {};
      state.allIds = {};
      state.meta = {};
      state.fetchComplete = false;
    },
    [types.FETCH_COMPLETED](state, { byId, allIds, meta }) {
      state.byId = byId;
      state.allIds = allIds;
      state.meta = meta;
      state.fetchComplete = true;
    },
    [types.FETCH_ERROR](state, error) {
      state.fetchError = error.message;
    },
  },

  actions: {
    async fetchList({ commit }, params = {}) {
      commit(types.FETCH_LIST);

      try {
        const data = await axios.get(API_HOST + API_ROUTES.eventsList, { params });

        if (isEmpty(data.data.data)) {
          return;
        }
        const normalizedData = normalize(data.data.data, [eventSchema]);
        commit(types.FETCH_COMPLETED, {
          byId: normalizedData.entities.event,
          allIds: normalizedData.result,
          meta: {
            total: data.data.total,
          },
        });
      } catch (error) {
        commit(types.FETCH_ERROR, error);
      }
    },
  },
};
