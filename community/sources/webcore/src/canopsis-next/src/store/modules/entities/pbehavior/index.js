import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import request from '@/services/request';

import i18n from '@/i18n';

import schemas from '@/store/schemas';

import commentModule from './comment';
import entitiesModule from './entities';

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
  modules: { comment: commentModule, entities: entitiesModule },
  state: {
    allIds: [],
    pending: false,
    fetchingParams: {},
    meta: {},
  },
  getters: {
    allIds: state => state.allIds,
    items: (state, getters, rootState, rootGetters) => rootGetters['entities/getList'](
      ENTITIES_TYPES.pbehavior,
      state.allIds,
    ),
    getItem: (state, getters, rootState, rootGetters) => id => rootGetters['entities/getItem'](
      ENTITIES_TYPES.pbehavior,
      id,
    ),
    pending: state => state.pending,
    meta: state => state.meta,
  },
  mutations: {
    [types.FETCH_LIST](state, { params } = {}) {
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
        commit(types.FETCH_LIST, { params });

        const { data, normalizedData } = await dispatch('entities/fetch', {
          route: API_ROUTES.pbehavior.pbehaviors,
          schema: [schemas.pbehavior],
          params,
          dataPreparer: d => d.data,
        }, { root: true });

        commit(types.FETCH_LIST_COMPLETED, {
          allIds: normalizedData.result,
          meta: data.meta,
        });
      } catch (err) {
        console.error(err);

        await dispatch('popups/error', { text: i18n.t('errors.default') }, { root: true });

        commit(types.FETCH_LIST_FAILED);
      }
    },

    fetchListWithPreviousParams({ dispatch, state }) {
      dispatch('fetchList', { params: state.fetchingParams });
    },

    fetchListWithoutStore(context, { params }) {
      return request.get(API_ROUTES.pbehavior.pbehaviors, { params });
    },

    async fetchListByEntityId({ commit, dispatch }, { params }) {
      try {
        const { normalizedData } = await dispatch('entities/fetch', {
          route: API_ROUTES.pbehavior.entities,
          schema: [schemas.pbehavior],
          params,
        }, { root: true });

        commit(types.FETCH_BY_ID_COMPLETED, { allIds: normalizedData.result });
      } catch (err) {
        commit(types.FETCH_BY_ID_FAILED, err);

        console.error(err);
      }
    },

    fetchListByEntityIdWithoutStore(context, { id, params = {} }) {
      return request.get(API_ROUTES.pbehavior.entities, { params: { _id: id, ...params } });
    },

    async create({ dispatch }, { data }) {
      try {
        const pbehavior = await request.post(API_ROUTES.pbehavior.pbehaviors, data);

        await dispatch('popups/success', { text: i18n.t('modals.createPbehavior.success.create') }, { root: true });

        return pbehavior;
      } catch (err) {
        console.error(err);

        await dispatch('popups/error', { text: i18n.t('errors.default') }, { root: true });

        throw err;
      }
    },

    bulkCreate(context, { data }) {
      return request.post(API_ROUTES.pbehavior.bulkPbehaviors, data);
    },

    fetchItemWithoutStore(context, { id }) {
      return request.get(`${API_ROUTES.pbehavior.pbehaviors}/${id}`);
    },

    async update({ dispatch }, { data, id }) {
      await dispatch('entities/update', {
        route: `${API_ROUTES.pbehavior.pbehaviors}/${id}`,
        schema: schemas.pbehavior,
        body: data,
      }, { root: true });
    },

    bulkUpdate(context, { data }) {
      return request.put(API_ROUTES.pbehavior.bulkPbehaviors, data);
    },

    async remove({ dispatch }, { id }) {
      try {
        await request.delete(`${API_ROUTES.pbehavior.pbehaviors}/${id}`);
        await dispatch('entities/removeFromStore', {
          id,
          type: ENTITIES_TYPES.pbehavior,
        }, { root: true });
      } catch (err) {
        console.warn(err);
      }
    },

    async removeWithoutStore(context, { id }) {
      return request.delete(`${API_ROUTES.pbehavior.pbehaviors}/${id}`);
    },

    bulkRemove(context, { data }) {
      return request.delete(API_ROUTES.pbehavior.bulkPbehaviors, { data });
    },

    fetchPbehaviorsCalendarWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.pbehavior.calendar, { params });
    },

    fetchEntitiesPbehaviorsCalendarWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.pbehavior.entitiesCalendar, { params });
    },

    bulkCreateEntityPbehaviors(context, { data } = {}) {
      return request.post(API_ROUTES.pbehavior.bulkEntityPbehaviors, data);
    },

    bulkRemoveEntityPbehaviors(context, { data }) {
      return request.delete(API_ROUTES.pbehavior.bulkEntityPbehaviors, { data });
    },
  },
};
