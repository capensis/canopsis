import { normalize } from 'normalizr';

import { API_ROUTES } from '@/config';
import { alarmSchema } from '@/store/schemas';
import request from '@/services/request';
import entitiesTypes from '@/store/modules/entities/types';

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
    pending: false,
    fetchingParams: {},
  },
  getters: {
    byId: state => state.byId,
    allIds: state => state.allIds,
    items: state => state.allIds.map(id => state.byId[id]),
    meta: state => state.meta,
    pending: state => state.pending,
  },
  mutations: {
    [types.FETCH_LIST](state) {
      state.byId = {};
      state.allIds = [];
      state.meta = {};
      state.pending = true;
    },
    [types.FETCH_LIST_COMPLETED](state, { byId, allIds, meta }) {
      state.byId = byId;
      state.allIds = allIds;
      state.meta = meta;
      state.pending = false;
    },
  },
  actions: {
    /* async fetchList({ commit }, params = {}) {
      commit(`entities/${entitiesTypes.ENTITIES_UPDATE}`, normalizedData.entities, { root: true });

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
    }, */
    async fetchList({ commit }, params = {}) {
      commit(types.FETCH_LIST);

      try {
        const [data] = await request.get(API_ROUTES.alarmList, params);
        const normalizedData = normalize(data.alarms, [alarmSchema]);

        commit(`entities/${entitiesTypes.ENTITIES_UPDATE}`, normalizedData.entities, { root: true });
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
