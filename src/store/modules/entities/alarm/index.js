import merge from 'lodash/merge';

import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';
import { alarmSchema } from '@/store/schemas';
import i18n from '@/i18n';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default {
  namespaced: true,
  state: {
    allIds: [],
    meta: {},
    pending: false,
    fetchingParams: {},
  },
  getters: {
    allIds: state => state.allIds,
    items: (state, getters, rootState, rootGetters) =>
      rootGetters['entities/getList'](ENTITIES_TYPES.alarm, state.allIds),
    item: (state, getters, rootState, rootGetters) => id =>
      rootGetters['entities/getItem'](ENTITIES_TYPES.alarm, id),
    meta: state => state.meta,
    pending: state => state.pending,
  },
  mutations: {
    [types.FETCH_LIST](state, { params }) {
      state.pending = true;
      state.fetchingParams = params;
    },
    [types.FETCH_LIST_COMPLETED](state, { allIds, meta }) {
      state.allIds = allIds;
      state.meta = meta;
      state.pending = false;
    },
    [types.FETCH_LIST_FAILED](state) {
      state.pending = false;
    },
  },
  actions: {
    async fetchList({ commit, dispatch }, { params } = {}) {
      try {
        commit(types.FETCH_LIST, { params });

        const { normalizedData, data } = await dispatch('entities/fetch', {
          route: API_ROUTES.alarmList,
          schema: [alarmSchema],
          params,
          dataPreparer: d => d.data[0].alarms,
        }, { root: true });
        commit(types.FETCH_LIST_COMPLETED, {
          allIds: normalizedData.result,
          meta: {
            first: data.data[0].first,
            last: data.data[0].last,
            total: data.data[0].total,
          },
        });
      } catch (err) {
        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });

        commit(types.FETCH_LIST_FAILED);
      }
    },

    fetchListWithPreviousParams({ dispatch, state }) {
      return dispatch('fetchList', { params: state.fetchingParams });
    },

    async fetchItem({ dispatch }, { id, params }) {
      try {
        const paramsWithItemId = merge(params, { filter: { d: id } });

        await dispatch('entities/fetch', {
          route: API_ROUTES.alarmList,
          schema: [alarmSchema],
          params: paramsWithItemId,
          dataPreparer: d => d.data[0].alarms,
        }, { root: true });
      } catch (err) {
        console.error(err);
      }
    },

  },
};
