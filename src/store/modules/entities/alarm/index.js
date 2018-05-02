import { normalize } from 'normalizr';

import { API_ROUTES } from '@/config';
import { alarmSchema } from '@/store/schemas';
import request from '@/services/request';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
};

export default {
  namespaced: true,
  state: {
    byId: {},
    allIds: [],
    meta: {},
    fetchComplete: false,
  },
  getters: {
    byId: state => state.byId,
    allIds: state => state.allIds,
    items: state => state.allIds.map(id => state.byId[id]),
    meta: state => state.meta,
    fetchComplete: state => state.fetchComplete,
  },
  mutations: {
    [types.FETCH_LIST](state) {
      state.byId = {};
      state.allIds = [];
      state.meta = {};
      state.fetchComplete = false;
    },
    [types.FETCH_LIST_COMPLETED](state, { byId, allIds, meta }) {
      state.byId = byId;
      state.allIds = allIds;
      state.meta = meta;
      state.fetchComplete = true;
    },
  },
  actions: {
    async fetchList({ commit }, params = {}) {
      commit(types.FETCH_LIST);

      try {
        const [data] = await request.get(API_ROUTES.alarmList, params);
        const normalizedData = normalize(data.alarms, [alarmSchema]);

        commit(types.FETCH_LIST_COMPLETED, {
          byId: normalizedData.entities.alarm,
          allIds: normalizedData.result,
          meta: {
            first: data.first,
            last: data.last,
            total: data.total,
          },
        });
      } catch (err) {
        console.error(err);
      }
    },
  },
};
