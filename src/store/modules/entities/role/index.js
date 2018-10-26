import i18n from '@/i18n';
import qs from 'qs';

import { roleSchema } from '@/store/schemas';
import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';
import request from '@/services/request';

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
    items: (state, getters, rootState, rootGetters) => rootGetters['entities/getList'](ENTITIES_TYPES.role, state.allIds),
    getItemById: (state, getters, rootState, rootGetters) => id => rootGetters['entities/getItem'](ENTITIES_TYPES.role, id),
    meta: state => state.meta,
    pending: state => state.pending,
  },
  mutations: {
    [types.FETCH_LIST](state, { params }) {
      state.pending = true;
      state.fetchingParams = params;
    },
    [types.FETCH_LIST_COMPLETED](state, { allIds, meta }) {
      state.pending = false;
      state.allIds = allIds;
      state.meta = meta;
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
          route: API_ROUTES.role.list,
          schema: [roleSchema],
          params,
          dataPreparer: d => d.data,
        }, { root: true });

        commit(types.FETCH_LIST_COMPLETED, {
          allIds: normalizedData.result,
          meta: {
            total: data.total,
          },
        });
      } catch (err) {
        commit(types.FETCH_LIST_FAILED);

        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });
      }
    },

    /**
   * Fetch roles list with previous params
   *
   * @param {Function} dispatch
   * @param {Object} state
   * @returns {*}
   */
    fetchListWithPreviousParams({ dispatch, state }) {
      return dispatch('fetchList', {
        params: state.fetchingParams,
      });
    },

    async create({ dispatch }, { data }) {
      try {
        await request.post(API_ROUTES.role.create, qs.stringify({ role: JSON.stringify(data) }), {
          headers: { 'content-type': 'application/x-www-form-urlencoded' },
        });
        await dispatch('popup/add', { type: 'success', text: i18n.t('success.default') }, { root: true });
      } catch (err) {
        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });
      }
    },

    async remove({ dispatch }, { id }) {
      try {
        await request.delete(`${API_ROUTES.role.remove}/${id}`);
        await dispatch('popup/add', { type: 'success', text: i18n.t('success.default') }, { root: true });
      } catch (err) {
        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });
      }
    },
  },
};
