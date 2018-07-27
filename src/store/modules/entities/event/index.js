import { entitySchema } from '@/store/schemas';
import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';
import qs from 'qs';
import i18n from '@/i18n';
import request from '@/services/request';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default {
  namespaced: true,
  state: {
    items: [],
    allIds: [],
    meta: {},
    pending: false,
    fetchingParams: {},
  },
  getters: {
    allIds: state => state.allIds,
    items: state => state.items,
    pending: state => state.pending,
    meta: state => state.meta,
  },
  mutations: {
    [types.FETCH_LIST](state, { params }) {
      state.pending = true;
      state.fetchingParams = params;
    },
    [types.FETCH_LIST_COMPLETED](state, { items, allIds, meta }) {
      state.items = items;
      state.allIds = allIds;
      state.meta = meta;
      state.pending = false;
    },
    [types.FETCH_LIST_FAILED](state) {
      state.pending = false;
    },
  },
  actions: {
    fetch({ dispatch }, { params } = {}) {
      return dispatch('entities/fetch', {
        route: API_ROUTES.events,
        schema: [entitySchema],
        params,
        dataPreparer: d => d.data,
      }, { root: true });
    },

    async fetchList({ commit, dispatch }, { params } = {}) {
      try {
        commit(types.FETCH_LIST, { params });
        const { normalizedData, data } = await dispatch('fetch', { params });
        commit(types.FETCH_LIST_COMPLETED, {
          items: normalizedData.entities.entity,
          allIds: normalizedData.result,
          meta: {
            total: data.total,
          },
        });
      } catch (err) {
        console.error(err);
        commit(types.FETCH_LIST_FAILED);
      }
    },

    async create({ dispatch }, { data }) {
      try {
        await request.post(API_ROUTES.event, qs.stringify({ event: JSON.stringify(data) }), {
          headers: { 'content-type': 'application/x-www-form-urlencoded' },
        });
        await dispatch('popup/add', { type: 'success', text: i18n.t('success.default') }, { root: true });
      } catch (e) {
        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });
        console.warn(e);
      }
    },
  },
};
