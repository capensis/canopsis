import { normalize } from 'normalizr';
import { isEmpty } from 'lodash';

import { API_ROUTES } from '@/config';
import request from '@/services/request';
import { eventSchema } from '@/store/schemas';

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
  },

  getters: {
    data: state => state.data,
    byId: state => state.byId,
    allIds: state => state.allIds,
    fetchComplete: state => state.fetchComplete,
    fetchError: state => state.fetchError,
  },

  mutations: {
    [types.FETCH_LIST](state) {
      state.byId = {};
      state.allIds = {};
      state.fetchComplete = false;
    },
    [types.FETCH_COMPLETED](state, { byId, allIds }) {
      state.byId = byId;
      state.allIds = allIds;
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
        const data = await request.get(API_ROUTES.eventsList, { params });

        if (isEmpty(data)) {
          return;
        }

        const normalizedData = normalize(data, [eventSchema]);

        commit(types.FETCH_COMPLETED, {
          byId: normalizedData.entities.event,
          allIds: normalizedData.result,
        });
      } catch (error) {
        commit(types.FETCH_ERROR, error);
      }
    },
  },
};
