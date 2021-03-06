import request from '@/services/request';
import i18n from '@/i18n';
import schemas from '@/store/schemas';
import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import commentModule from './comment';

const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
  FETCH_BY_ID: 'FETCH_BY_ID',
  FETCH_BY_ID_COMPLETED: 'FETCH_BY_ID_COMPLETED',
  FETCH_BY_ID_FAILED: 'FETCH_BY_ID_FAILED',
};

export default {
  namespaced: true,
  modules: { comment: commentModule },
  state: {
    allIds: [],
    pending: false,
    meta: {},
  },
  getters: {
    allIds: state => state.allIds,
    items: (state, getters, rootState, rootGetters) =>
      rootGetters['entities/getList'](ENTITIES_TYPES.pbehavior, state.allIds),
    getItem: (state, getters, rootState, rootGetters) => id =>
      rootGetters['entities/getItem'](ENTITIES_TYPES.pbehavior, id),
    pending: state => state.pending,
    meta: state => state.meta,
  },
  mutations: {
    [types.FETCH_LIST](state) {
      state.pending = true;
    },
    [types.FETCH_LIST_COMPLETED](state, { allIds, meta }) {
      state.allIds = allIds;
      state.meta = meta;
      state.pending = false;
    },
    [types.FETCH_LIST_FAILED](state) {
      state.pending = false;
    },
    [types.FETCH_BY_ID](state) {
      state.pending = true;
    },
    [types.FETCH_BY_ID_COMPLETED](state, { allIds }) {
      state.allIds = allIds;
      state.pending = false;
    },
    [types.FETCH_BY_ID_FAILED](state) {
      state.pending = false;
    },
  },
  actions: {
    async fetchList({ dispatch, commit }, { params } = {}) {
      try {
        commit(types.FETCH_LIST);

        const { data, normalizedData } = await dispatch('entities/fetch', {
          route: API_ROUTES.planning.pbehaviors,
          schema: [schemas.pbehavior],
          params,
          dataPreparer: d => d.data,
        }, { root: true });

        commit(types.FETCH_LIST_COMPLETED, {
          allIds: normalizedData.result,
          meta: data.meta,
        });
      } catch (err) {
        await dispatch('popups/error', { text: i18n.t('errors.default') }, { root: true });
        commit(types.FETCH_LIST_FAILED);
      }
    },

    async fetchListByEntityId({ commit, dispatch }, { id }) {
      try {
        const { normalizedData } = await dispatch('entities/fetch', {
          route: API_ROUTES.planning.pbehaviorById,
          schema: [schemas.pbehavior],
          params: { id },
        }, { root: true });

        commit(types.FETCH_BY_ID_COMPLETED, { allIds: normalizedData.result });
      } catch (err) {
        commit(types.FETCH_BY_ID_FAILED, err);

        console.warn(err);
      }
    },

    async create({ dispatch }, { data }) {
      try {
        const pbehavior = await request.post(API_ROUTES.planning.pbehaviors, data);

        await dispatch('popups/success', { text: i18n.t('modals.createPbehavior.success.create') }, { root: true });

        return pbehavior;
      } catch (err) {
        console.error(err);
        await dispatch('popups/error', { text: i18n.t('errors.default') }, { root: true });

        throw err;
      }
    },

    async update({ dispatch }, { data, id }) {
      await dispatch('entities/update', {
        route: `${API_ROUTES.planning.pbehaviors}/${id}`,
        schema: schemas.pbehavior,
        body: data,
      }, { root: true });
    },

    async remove({ dispatch }, { id }) {
      try {
        await request.delete(`${API_ROUTES.planning.pbehaviors}/${id}`);
        await dispatch('entities/removeFromStore', {
          id,
          type: ENTITIES_TYPES.pbehavior,
        }, { root: true });
      } catch (err) {
        console.warn(err);
      }
    },
  },
};

