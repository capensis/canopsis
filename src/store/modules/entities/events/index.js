import { normalize } from 'normalizr';
import { isEmpty } from 'lodash';

import { API_ROUTES } from '@/config';
import request from '@/services/request';
import { eventSchema } from '@/store/schemas';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_COMPLETED: 'FETCH_LIST_COMPLETED',
};

export default {
  namespaced: true,

  state: {
    data: {},
    byId: {},
    allIds: {},
    meta: {},
    fetchComplete: false,
  },

  getters: {
    data: state => state.data,
    byId: state => state.byId,
    allIds: state => state.allIds,
    meta: state => state.meta,
    fetchComplete: state => state.fetchComplete,
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
          meta: {
            first: data.first,
            last: data.last,
            total: data.total,
          },
        });
      } catch (error) {
        console.warn(error);
      }
    },
  },
};
