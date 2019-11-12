import qs from 'qs';

import i18n from '@/i18n';
import request from '@/services/request';
import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';
import { userSchema } from '@/store/schemas';

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
    getItemById: (state, getters, rootState, rootGetters) => id =>
      rootGetters['entities/getItem'](ENTITIES_TYPES.user, id),

    items: (state, getters, rootState, rootGetters) =>
      rootGetters['entities/getList'](ENTITIES_TYPES.user, state.allIds),

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
    /**
     * Action for user creating/editing
     *
     * @param {Function} dispatch
     * @param {Object} data
     * @returns {Promise<void>}
     */
    async create({ dispatch }, { data }) {
      try {
        await request.post(API_ROUTES.user.create, qs.stringify({ user: JSON.stringify(data) }), {
          headers: { 'content-type': 'application/x-www-form-urlencoded' },
        });

        await dispatch('popups/success', { text: i18n.t('success.default') }, { root: true });
      } catch (err) {
        await dispatch('popups/error', { text: i18n.t('errors.default') }, { root: true });
      }
    },

    /**
     * Action for user removing
     *
     * @param {Function} dispatch
     * @param {string} id
     * @returns {Promise<void>}
     */
    async remove({ dispatch }, { id }) {
      try {
        await request.delete(`${API_ROUTES.user.remove}/${id}`);
        await dispatch('popups/success', { text: i18n.t('success.default') }, { root: true });
      } catch (err) {
        await dispatch('popups/error', { text: i18n.t('errors.default') }, { root: true });
      }
    },

    /**
     * Fetch users list by params
     *
     * @param {Function} commit
     * @param {Function} dispatch
     * @param {Object} params
     * @returns {Promise<void>}
     */
    async fetchList({ commit, dispatch }, { params } = {}) {
      try {
        commit(types.FETCH_LIST, { params });

        const { normalizedData, data } = await dispatch('entities/fetch', {
          route: API_ROUTES.user.list,
          schema: [userSchema],
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

        await dispatch('popups/error', { text: i18n.t('errors.default') }, { root: true });
      }
    },

    /**
     * Fetch users list with previous params
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
  },
};
