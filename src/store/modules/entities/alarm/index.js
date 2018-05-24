import { normalize } from 'normalizr';

import { API_ROUTES } from '@/config';
import { alarmSchema } from '@/store/schemas';
import request from '@/services/request';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_ITEM: 'FETCH_ITEM',
  FETCH_ITEM_COMPLETED: 'FETCH_ITEM_COMPLETED',
};

export default {
  namespaced: true,
  state: {
    byId: {},
    allIds: [],
    meta: {},
    fetchComplete: false,
    itemPending: false,
    item: {},
  },
  getters: {
    byId: state => state.byId,
    allIds: state => state.allIds,
    items: state => state.allIds.map(id => state.byId[id]),
    meta: state => state.meta,
    fetchComplete: state => state.fetchComplete,
    itemPending: state => state.itemPending,
    item: state => state.item,
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
    [types.FETCH_ITEM](state) {
      state.item = {};
      state.itemPending = true;
    },
    [types.FETCH_ITEM_COMPLETED](state, { item }) {
      state.item = item;
      state.itemPending = false;
    },
  },
  actions: {
    async fetchItem({ commit }, params = {}) {
      commit(types.FETCH_ITEM);
      try {
        const [data] = await request.get(API_ROUTES.alarmList, params);
        const normalizedData = normalize(data.alarms, [alarmSchema]);
        const itemId = normalizedData.result[0];
        commit(types.FETCH_ITEM_COMPLETED, {
          item: normalizedData.entities.alarm[itemId],
        });
      } catch (err) {
        console.log(err);
      }
    },
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
